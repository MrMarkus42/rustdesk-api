package admin

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

// Lifetime of the login token the web client hands to hbbs (--must-login) in its punch request.
const webClientPunchTokenTTL = 10 * time.Minute

type WebClient struct {
}

func webClientEnabled(c *gin.Context) bool {
	if global.Config.App.WebClient != 1 {
		response.Fail(c, 101, "web client disabled")
		return false
	}
	return true
}

// Config 浏览器客户端配置
// @Tags WEBCLIENT
// @Summary Web client configuration
// @Description Server key and a short-lived login token for hbbs
// @Produce  json
// @Success 200 {object} response.Response
// @Router /admin/webclient/config [get]
// @Security token
func (wc *WebClient) Config(c *gin.Context) {
	if !webClientEnabled(c) {
		return
	}
	u := service.AllService.UserService.CurUser(c)
	// hbbs only checks the signature and expiry of this token. Without a JWT key it checks
	// nothing but that one is present, so the session token is passed on unchanged.
	token := ""
	if len(global.Jwt.Key) > 0 {
		token = global.Jwt.GenerateTokenWithExpire(u.Id, webClientPunchTokenTTL)
	} else {
		token = c.GetString("token")
	}
	name := u.Nickname
	if name == "" {
		name = u.Username
	}
	response.Success(c, gin.H{
		"key":         strings.TrimSpace(global.Config.Rustdesk.Key),
		"punch_token": token,
		"username":    u.Username,
		"name":        name,
	})
}

// Ticket 浏览器客户端一次性票据
// @Tags WEBCLIENT
// @Summary One-time ticket for one WebSocket connection
// @Produce  json
// @Success 200 {object} response.Response
// @Router /admin/webclient/ticket [post]
// @Security token
func (wc *WebClient) Ticket(c *gin.Context) {
	if !webClientEnabled(c) {
		return
	}
	u := service.AllService.UserService.CurUser(c)
	t := service.WebClient.NewTicket(u.Id)
	if t == "" {
		response.Fail(c, 101, "too many requests")
		return
	}
	response.Success(c, gin.H{"ticket": t})
}

// Peer 浏览器客户端设备信息
// @Tags WEBCLIENT
// @Summary Address book entry of the current user for a device (alias, saved password hash)
// @Produce  json
// @Param id query string true "device id"
// @Success 200 {object} response.Response
// @Router /admin/webclient/peer [get]
// @Security token
func (wc *WebClient) Peer(c *gin.Context) {
	if !webClientEnabled(c) {
		return
	}
	id := c.Query("id")
	if id == "" || len(id) > 100 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	u := service.AllService.UserService.CurUser(c)
	ab := service.AllService.AddressBookService.InfoByUserIdAndId(u.Id, id)
	if ab.RowId == 0 {
		response.Success(c, gin.H{"id": id})
		return
	}
	response.Success(c, gin.H{
		"id":       id,
		"alias":    ab.Alias,
		"hostname": ab.Hostname,
		"username": ab.Username,
		"platform": ab.Platform,
		"hash":     ab.Hash,
	})
}
