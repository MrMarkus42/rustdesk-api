package web

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

// The browser client reaches hbbs/hbbr through this proxy, so it needs no extra ports and
// nobody without a login can use it. Browsers cannot send headers with a WebSocket, so the
// one-time ticket travels as a second subprotocol: "rustdesk-web, t.<ticket>".
const (
	webClientSubprotocol  = "rustdesk-web"
	webClientTicketPrefix = "t."
	webClientMaxConns     = 256
)

var webClientConns atomic.Int32

type WebClientWs struct {
}

func ticketFromSubprotocols(header string) (ticket string, ok bool) {
	hasProto := false
	for _, p := range strings.Split(header, ",") {
		p = strings.TrimSpace(p)
		if p == webClientSubprotocol {
			hasProto = true
		} else if strings.HasPrefix(p, webClientTicketPrefix) {
			ticket = strings.TrimPrefix(p, webClientTicketPrefix)
		}
	}
	return ticket, hasProto && ticket != ""
}

func wsTarget(raw string) (*url.URL, bool) {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return nil, false
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	default:
		return nil, false
	}
	return u, true
}

// Proxy forwards /webclient/ws/id to hbbs and /webclient/ws/relay to hbbr.
func (w *WebClientWs) Proxy(c *gin.Context) {
	if global.Config.App.WebClient != 1 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var raw string
	switch c.Param("kind") {
	case "id":
		raw = global.Config.Rustdesk.WsIdTarget
	case "relay":
		raw = global.Config.Rustdesk.WsRelayTarget
	default:
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ticket, ok := ticketFromSubprotocols(c.GetHeader("Sec-WebSocket-Protocol"))
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := service.WebClient.UseTicket(ticket)
	if userId == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user := service.AllService.UserService.InfoById(userId)
	if !service.AllService.UserService.CheckUserEnable(user) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	target, ok := wsTarget(raw)
	if !ok {
		global.Logger.Error("webclient: invalid ws target ", raw)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	if webClientConns.Add(1) > webClientMaxConns {
		webClientConns.Add(-1)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	defer webClientConns.Add(-1)

	kind := c.Param("kind")
	clientIp := c.ClientIP()
	global.Logger.Infof("webclient: user %s (%s) opened %s connection", user.Username, clientIp, kind)

	proxy := &httputil.ReverseProxy{
		// Rewrite starts from an outgoing request without X-Forwarded-* headers.
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.URL.Path = target.Path
			r.Out.URL.RawPath = ""
			r.Out.URL.RawQuery = ""
			r.Out.Host = target.Host
			for _, h := range []string{"Cookie", "Authorization", "Api-Token", "Origin",
				"Sec-WebSocket-Protocol", "X-Real-Ip", "X-Forwarded-For"} {
				r.Out.Header.Del(h)
			}
			// hbbr blocks and limits by client IP; hbbs keys pending punch requests by
			// address, so it keeps the unique local address of this connection instead.
			if kind == "relay" {
				r.Out.Header.Set("X-Real-IP", clientIp)
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			if resp.StatusCode == http.StatusSwitchingProtocols {
				// The browser drops a connection whose server picked no offered subprotocol.
				resp.Header.Set("Sec-WebSocket-Protocol", webClientSubprotocol)
			}
			return nil
		},
		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			global.Logger.Warn("webclient: ", kind, " proxy error: ", err)
			rw.WriteHeader(http.StatusBadGateway)
		},
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
