package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"common/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/jwt"
)

type User = model.User

var identityKey = "username"

func RegsterAuth(h *server.Hertz) {
	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:            "test zone",
		SigningAlgorithm: "HS256",
		Key:              []byte("panic-bitdance"),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		// 登录验证器
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			// 处理登录逻辑
			// 返回token
			return map[string]interface{}{
				"username": "1234",
			}, nil
			// 出错
			// return nil, jwt.ErrFailedAuthentication
		},
		// 是否可以访问Realm
		// 鉴权验证器
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			if v, ok := data.(*User); ok && v.Username == "admin" {
				return true
			}
			return false
		},
		// 生成token时放入的参数
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		// 无权限的响应
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},
		// 登录鉴权的响应
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		// 登出的响应
		// 登出可能会自动删除token
		LogoutResponse: func(ctx context.Context, c *app.RequestContext, code int) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": http.StatusOK,
			})
		},
		// refresh token
		RefreshResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		// 不知道干嘛的
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &User{
				Username: claims[identityKey].(string),
			}
		},
		IdentityKey: identityKey,
		// 从请求头的Authorization提取token
		TokenLookup:                 "header: Authorization",
		TokenHeadName:               "Bearer",
		WithoutDefaultTokenHeadName: false,
		TimeFunc:                    time.Now,
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			return e.Error()
		},
		SendAuthorization: true,
		DisabledAbort:     false,
	})
	if err != nil {
		log.Fatal("auth middleware Error:" + err.Error())
	}
	// 之后在这里写handler
	// h.POST("/login", authMiddleware.LoginHandler)
	// h.POST("/logout", authMiddleware.LogoutHandler)
	// h.NoRoute(authMiddleware.MiddlewareFunc(), func(ctx context.Context, c *app.RequestContext) {
	// 	claims := jwt.ExtractClaims(ctx, c)
	// 	hlog.Infof("NoRoute claims: %#v\n", claims)
	// 	c.JSON(404, map[string]string{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	// })
	// 该组接口都需要走middleware
	auth := h.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	// auth.Use(authMiddleware.MiddlewareFunc())
	// {
	// 	auth.GET("/ping", PingHandler)
	// }
}
