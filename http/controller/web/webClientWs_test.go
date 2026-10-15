package web

import "testing"

func TestTicketFromSubprotocols(t *testing.T) {
	cases := []struct {
		header string
		ticket string
		ok     bool
	}{
		{"rustdesk-web, t.abc", "abc", true},
		{"t.abc, rustdesk-web", "abc", true},
		{"t.abc", "", false},
		{"rustdesk-web", "", false},
		{"", "", false},
		{"other, t.abc", "", false},
	}
	for _, c := range cases {
		tk, ok := ticketFromSubprotocols(c.header)
		if ok != c.ok || (ok && tk != c.ticket) {
			t.Errorf("%q: got (%q, %v), want (%q, %v)", c.header, tk, ok, c.ticket, c.ok)
		}
	}
}

func TestWsTarget(t *testing.T) {
	if u, ok := wsTarget("ws://127.0.0.1:21118"); !ok || u.Scheme != "http" || u.Host != "127.0.0.1:21118" {
		t.Errorf("ws target: %v %v", u, ok)
	}
	if u, ok := wsTarget("wss://hbbs.example.org/ws/id"); !ok || u.Scheme != "https" {
		t.Errorf("wss target: %v %v", u, ok)
	}
	for _, bad := range []string{"", "http://x", "ws://", "file:///etc/passwd"} {
		if _, ok := wsTarget(bad); ok {
			t.Errorf("%q accepted", bad)
		}
	}
}
