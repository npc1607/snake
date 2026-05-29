package session

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"snake/internal/game"
	"snake/internal/helper"
)

// DefaultAddress is the default LAN listen address.
const DefaultAddress = ":7777"

// PlayerInput is one remote player direction update.
type PlayerInput struct {
	PlayerID  int
	Direction game.Direction
}

type wireMessage struct {
	Type         string         `json:"type"`
	PlayerID     int            `json:"player_id,omitempty"`
	Direction    game.Direction `json:"direction,omitempty"`
	State        *game.State    `json:"state,omitempty"`
	PingID       int64          `json:"ping_id,omitempty"`
	SentUnixNano int64          `json:"sent_unix_nano,omitempty"`
	LatencyMS    int64          `json:"latency_ms,omitempty"`
}

// Host manages one LAN host session.
type Host struct {
	listener net.Listener
	inputs   chan PlayerInput

	mu      sync.RWMutex
	conn    net.Conn
	latency int64
	writeMu sync.Mutex
}

// Addr returns the listening address.
func (h *Host) Addr() string {
	return h.listener.Addr().String()
}

// NewHost starts a host session.
func NewHost(address string) (*Host, error) {
	if address == "" {
		address = DefaultAddress
	}

	listener, err := net.Listen(helper.NetworkForAddress(address), address)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", address, err)
	}

	host := &Host{
		listener: listener,
		inputs:   make(chan PlayerInput, 32),
	}
	go host.acceptLoop()

	return host, nil
}

// Inputs returns the remote input queue.
func (h *Host) Inputs() <-chan PlayerInput {
	return h.inputs
}

// Ready reports whether a remote player is connected.
func (h *Host) Ready() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.conn != nil
}

// RemoteAddr returns the currently connected player's remote address.
func (h *Host) RemoteAddr() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.conn == nil {
		return ""
	}

	return h.conn.RemoteAddr().String()
}

// LatencyMS returns the latest reported remote player round-trip latency in milliseconds.
func (h *Host) LatencyMS() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return int(h.latency)
}

// SendSnapshot pushes a full game snapshot to the client.
func (h *Host) SendSnapshot(state game.State) error {
	conn := h.currentConn()
	if conn == nil {
		return nil
	}

	return h.writeMessage(conn, wireMessage{
		Type:  "snapshot",
		State: &state,
	})
}

// Close shuts down the host session.
func (h *Host) Close() error {
	var firstErr error
	if err := h.listener.Close(); err != nil {
		firstErr = err
	}

	conn := h.currentConn()
	if conn != nil {
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func (h *Host) acceptLoop() {
	for {
		conn, err := h.listener.Accept()
		if err != nil {
			return
		}

		if err := h.handleConnection(conn); err != nil {
			_ = conn.Close()
			continue
		}
	}
}

func (h *Host) handleConnection(conn net.Conn) error {
	decoder := json.NewDecoder(conn)

	var joinMessage wireMessage
	if err := decoder.Decode(&joinMessage); err != nil {
		return fmt.Errorf("read join: %w", err)
	}
	if joinMessage.Type != "join" {
		return fmt.Errorf("unexpected message type %q", joinMessage.Type)
	}

	h.setConn(conn)
	if err := h.writeMessage(conn, wireMessage{
		Type:     "welcome",
		PlayerID: 2,
	}); err != nil {
		h.clearConn(conn)
		return fmt.Errorf("write welcome: %w", err)
	}

	h.readLoop(conn, decoder)
	h.clearConn(conn)

	return nil
}

func (h *Host) readLoop(conn net.Conn, decoder *json.Decoder) {
	for {
		var message wireMessage
		if err := decoder.Decode(&message); err != nil {
			return
		}

		switch message.Type {
		case "input":
			playerID := message.PlayerID
			if playerID <= 0 {
				playerID = 2
			}

			input := PlayerInput{
				PlayerID:  playerID,
				Direction: message.Direction,
			}

			select {
			case h.inputs <- input:
			default:
				<-h.inputs
				h.inputs <- input
			}
		case "ping":
			_ = h.writeMessage(conn, wireMessage{
				Type:         "pong",
				PingID:       message.PingID,
				SentUnixNano: message.SentUnixNano,
			})
		case "latency":
			h.setLatency(message.LatencyMS)
		default:
			continue
		}
	}
}

func (h *Host) writeMessage(conn net.Conn, message wireMessage) error {
	h.writeMu.Lock()
	defer h.writeMu.Unlock()

	if err := json.NewEncoder(conn).Encode(message); err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	return nil
}

func (h *Host) currentConn() net.Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.conn
}

func (h *Host) setConn(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conn != nil {
		_ = h.conn.Close()
	}
	h.conn = conn
	h.latency = 0
}

func (h *Host) clearConn(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conn == conn {
		_ = h.conn.Close()
		h.conn = nil
		h.latency = 0
	}
}

func (h *Host) setLatency(latency int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if latency < 0 {
		latency = 0
	}
	h.latency = latency
}

// Client manages one LAN client session.
type Client struct {
	conn      net.Conn
	playerID  int
	snapshots chan game.State
	latencyMu sync.RWMutex
	latency   int64
	writeMu   sync.Mutex
}

// PlayerID returns the assigned player number.
func (c *Client) PlayerID() int {
	return c.playerID
}

// LocalAddr returns the client's local connection address.
func (c *Client) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

// LatencyMS returns the latest host round-trip latency in milliseconds.
func (c *Client) LatencyMS() int {
	c.latencyMu.RLock()
	defer c.latencyMu.RUnlock()

	return int(c.latency)
}

// NewClient connects to a host session.
func NewClient(address string) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("join address is required")
	}

	conn, err := net.DialTimeout(helper.NetworkForAddress(address), address, 3*time.Second)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", address, err)
	}

	client := &Client{
		conn:      conn,
		snapshots: make(chan game.State, 32),
	}

	if err := client.writeMessage(wireMessage{Type: "join"}); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("write join: %w", err)
	}

	decoder := json.NewDecoder(conn)
	var welcome wireMessage
	if err := decoder.Decode(&welcome); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("read welcome: %w", err)
	}
	if welcome.Type != "welcome" {
		_ = conn.Close()
		return nil, fmt.Errorf("unexpected welcome type %q", welcome.Type)
	}

	client.playerID = welcome.PlayerID
	go client.readLoop(decoder)
	go client.pingLoop()

	return client, nil
}

// Snapshots returns the host snapshot queue.
func (c *Client) Snapshots() <-chan game.State {
	return c.snapshots
}

// SendInput sends one local direction to the host.
func (c *Client) SendInput(direction game.Direction) error {
	return c.writeMessage(wireMessage{
		Type:      "input",
		PlayerID:  c.playerID,
		Direction: direction,
	})
}

// Close shuts down the client session.
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) readLoop(decoder *json.Decoder) {
	for {
		var message wireMessage
		if err := decoder.Decode(&message); err != nil {
			return
		}

		switch message.Type {
		case "snapshot":
			if message.State == nil {
				continue
			}

			select {
			case c.snapshots <- *message.State:
			default:
				<-c.snapshots
				c.snapshots <- *message.State
			}
		case "pong":
			c.recordLatency(message.SentUnixNano)
		default:
			continue
		}
	}
}

func (c *Client) pingLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().UnixNano()
		if err := c.writeMessage(wireMessage{
			Type:         "ping",
			PingID:       now,
			SentUnixNano: now,
		}); err != nil {
			return
		}
	}
}

func (c *Client) recordLatency(sentUnixNano int64) {
	if sentUnixNano <= 0 {
		return
	}

	latency := time.Since(time.Unix(0, sentUnixNano)).Milliseconds()
	if latency < 1 {
		latency = 1
	}

	c.latencyMu.Lock()
	c.latency = latency
	c.latencyMu.Unlock()

	_ = c.writeMessage(wireMessage{
		Type:      "latency",
		LatencyMS: latency,
	})
}

func (c *Client) writeMessage(message wireMessage) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if err := json.NewEncoder(c.conn).Encode(message); err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	return nil
}
