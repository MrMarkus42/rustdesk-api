package service

import (
	"testing"
	"time"
)

func TestWebClientTicketOnce(t *testing.T) {
	ws := &WebClientService{tickets: map[string]webClientTicket{}}
	tk := ws.NewTicket(7)
	if len(tk) != 64 {
		t.Fatalf("ticket length %d", len(tk))
	}
	if uid := ws.UseTicket(tk); uid != 7 {
		t.Fatalf("first use: got user %d", uid)
	}
	if uid := ws.UseTicket(tk); uid != 0 {
		t.Fatalf("second use must fail, got user %d", uid)
	}
	if uid := ws.UseTicket("short"); uid != 0 {
		t.Fatalf("malformed ticket accepted")
	}
}

func TestWebClientTicketExpiry(t *testing.T) {
	ws := &WebClientService{tickets: map[string]webClientTicket{}}
	tk := ws.NewTicket(1)
	ws.tickets[tk] = webClientTicket{userId: 1, expires: time.Now().Add(-time.Second)}
	if uid := ws.UseTicket(tk); uid != 0 {
		t.Fatalf("expired ticket accepted")
	}
}

func TestWebClientTicketLimitPerUser(t *testing.T) {
	ws := &WebClientService{tickets: map[string]webClientTicket{}}
	for i := 0; i < webClientTicketPerUser; i++ {
		if ws.NewTicket(3) == "" {
			t.Fatalf("ticket %d refused", i)
		}
	}
	if ws.NewTicket(3) != "" {
		t.Fatalf("limit per user not enforced")
	}
	if ws.NewTicket(4) == "" {
		t.Fatalf("other user blocked by the limit")
	}
}
