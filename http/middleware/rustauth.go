package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

func RustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Println(c.Request.URL, c.Request.Header)
		//获取HTTP_AUTHORIZATION
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		if len(token) <= 7 {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		//提取token，格式是Bearer {token}
		//这里只是简单的提取
		token = token[7:]

		//验证token

		//检查是否设置了jwt key
		// Only the signature is checked here: the expiry inside the JWT is fixed at login,
		// while the user_tokens row below is extended on use. Checking the JWT expiry
		// logged every client out a fixed time after login, however actively it was used.
		var jwtUid uint
		if len(global.Jwt.Key) > 0 {
			jwtUid, _ = service.AllService.UserService.VerifyJWTSignature(token)
			if jwtUid == 0 {
				c.JSON(401, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}
		}

		user, ut := service.AllService.UserService.InfoByAccessToken(token)
		if user.Id == 0 || (jwtUid != 0 && jwtUid != user.Id) {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		if !service.AllService.UserService.CheckUserEnable(user) {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("curUser", user)
		c.Set("token", token)

		service.AllService.UserService.AutoRefreshAccessToken(ut)

		c.Next()
	}
}
