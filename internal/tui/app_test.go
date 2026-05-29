package tui

import (
	"testing"
	"time"
)

func TestJoinAppReceivesSnapshotFromHostApp(t *testing.T) {
	hostApp := NewHostApp("127.0.0.1:0")
	if hostApp.err != nil {
		t.Fatalf("new host app: %v", hostApp.err)
	}
	defer func() {
		if hostApp.hostSession != nil {
			_ = hostApp.hostSession.Close()
		}
	}()

	joinApp := NewJoinApp(hostApp.hostSession.Addr())
	connectMsg := connectClient(joinApp.joinAddr)()
	model, cmd := joinApp.Update(connectMsg)
	joinApp = model.(App)
	if joinApp.err != nil {
		t.Fatalf("join app connect: %v", joinApp.err)
	}
	if joinApp.clientSession == nil {
		t.Fatal("join app did not store client session")
	}

	deadline := time.After(2 * time.Second)
	for !hostApp.hostSession.Ready() {
		select {
		case <-deadline:
			t.Fatal("host app never became ready")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	model, _ = hostApp.Update(tickMsg(time.Now()))
	hostApp = model.(App)

	stateMessage := cmd()
	model, _ = joinApp.Update(stateMessage)
	joinApp = model.(App)

	if joinApp.remoteState.Width == 0 || joinApp.remoteState.Height == 0 {
		t.Fatalf("join app remote state = %+v, want non-zero board size", joinApp.remoteState)
	}
	if len(joinApp.remoteState.Players) != 2 {
		t.Fatalf("join app player count = %d, want 2", len(joinApp.remoteState.Players))
	}
}
