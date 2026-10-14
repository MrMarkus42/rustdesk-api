package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/middleware"
	"github.com/lejianwen/rustdesk-api/v2/http/router"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func ApiInit() {
	gin.SetMode(global.Config.Gin.Mode)
	g := gin.New()

	//[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
	//Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
	// Without trust-proxy gin trusted every proxy, so any client could set X-Forwarded-For
	// and get a fresh "IP" per request, bypassing the login captcha/ban and faking login-log IPs.
	// Behind a reverse proxy, set trust-proxy to the proxy's address.
	var pro []string
	if global.Config.Gin.TrustProxy != "" {
		pro = strings.Split(global.Config.Gin.TrustProxy, ",")
	}
	if err := g.SetTrustedProxies(pro); err != nil {
		panic(err)
	}

	if global.Config.Gin.Mode == gin.ReleaseMode {
		//修改gin Recovery日志 输出为logger的输出点
		if global.Logger != nil {
			gin.DefaultErrorWriter = global.Logger.WriterLevel(logrus.ErrorLevel)
		}
	}
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 not found")
	})
	g.Use(middleware.Logger(), middleware.Limiter(), gin.Recovery())
	router.WebInit(g)
	router.Init(g)
	router.ApiInit(g)
	Run(g, global.Config.Gin.ApiAddr)
}
