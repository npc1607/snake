package session

import (
	"testing"
	"time"

	"snake/internal/game"
)

func TestHostClientSnapshotFlow(t *testing.T) {
	host, err := NewHost("127.0.0.1:0")
	if err != nil {
		t.Fatalf("new host: %v", err)
	}
	defer func() {
		if err := host.Close(); err != nil {
			t.Fatalf("close host: %v", err)
		}
	}()

	client, err := NewClient(host.Addr())
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Fatalf("close client: %v", err)
		}
	}()

	state := game.State{
		Width:  12,
		Height: 8,
		Players: []game.PlayerState{
			{
				ID:    1,
				Snake: []game.Point{{X: 3, Y: 4}},
				Alive: true,
			},
		},
	}

	deadline := time.After(2 * time.Second)
	for !host.Ready() {
		select {
		case <-deadline:
			t.Fatal("host never became ready")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	if err := host.SendSnapshot(state); err != nil {
		t.Fatalf("send snapshot: %v", err)
	}

	select {
	case got := <-client.Snapshots():
		if got.Width != state.Width || got.Height != state.Height {
			t.Fatalf("snapshot size = %dx%d, want %dx%d", got.Width, got.Height, state.Width, state.Height)
		}
		if len(got.Players) != 1 || len(got.Players[0].Snake) != 1 {
			t.Fatalf("snapshot players = %+v, want one player with one segment", got.Players)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for snapshot")
	}

	deadline = time.After(3 * time.Second)
	for client.LatencyMS() == 0 || host.LatencyMS() == 0 {
		select {
		case <-deadline:
			t.Fatalf("latency was not reported: client=%d host=%d", client.LatencyMS(), host.LatencyMS())
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
