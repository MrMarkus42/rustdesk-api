package service

import (
	crand "crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Tickets let the browser open the WebSocket proxy (/webclient/ws/...), which cannot carry the
// api-token header. A ticket is random, bound to one user, valid for a short time and only once.
const (
	webClientTicketTTL     = 30 * time.Second
	webClientTicketMax     = 2000
	webClientTicketPerUser = 20
)

type webClientTicket struct {
	userId  uint
	expires time.Time
}

type WebClientService struct {
	mu      sync.Mutex
	tickets map[string]webClientTicket
}

var WebClient = &WebClientService{tickets: map[string]webClientTicket{}}

// pruneLocked drops expired tickets; the caller holds mu.
func (ws *WebClientService) pruneLocked(now time.Time) {
	for k, t := range ws.tickets {
		if now.After(t.expires) {
			delete(ws.tickets, k)
		}
	}
}

// NewTicket returns a one-time ticket for userId, or "" when too many are outstanding.
func (ws *WebClientService) NewTicket(userId uint) string {
	b := make([]byte, 32)
	if _, err := crand.Read(b); err != nil {
		return ""
	}
	ticket := hex.EncodeToString(b)
	now := time.Now()
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.pruneLocked(now)
	if len(ws.tickets) >= webClientTicketMax {
		return ""
	}
	n := 0
	for _, t := range ws.tickets {
		if t.userId == userId {
			n++
		}
	}
	if n >= webClientTicketPerUser {
		return ""
	}
	ws.tickets[ticket] = webClientTicket{userId: userId, expires: now.Add(webClientTicketTTL)}
	return ticket
}

// UseTicket consumes a ticket and returns its user id, or 0 if it is unknown or expired.
func (ws *WebClientService) UseTicket(ticket string) uint {
	if len(ticket) != 64 {
		return 0
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	t, ok := ws.tickets[ticket]
	if !ok {
		return 0
	}
	delete(ws.tickets, ticket)
	if time.Now().After(t.expires) {
		return 0
	}
	return t.userId
}
