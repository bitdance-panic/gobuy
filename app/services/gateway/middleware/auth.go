package middleware

import (
	"context"
	"time"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
)

var AuthMiddleware *jwt.HertzJWTMiddleware

func InitJWT() {

	var err error
	AuthMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "gobuy",
		Key:           []byte("panic-bitdance"),
		Timeout:       12 * time.Hour,
		MaxRefresh:    7 * 24 * time.Hour,
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		Authenticator: authenticator,
		Unauthorized:  unauthorizedHandler,
	})

	if err != nil {
		hlog.Fatalf("JWT中间件初始化失败: %v", err)
	}
}

func authenticator(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&loginReq); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	user, err := dao.GetUserByUsername(loginReq.Username)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	// 密码验证逻辑需要自行实现
	// if !checkPassword(loginReq.Password, user.PasswordHashed) {
	//     return nil, jwt.ErrFailedAuthentication
	// }

	return user.ID, nil
}

func unauthorizedHandler(ctx context.Context, c *app.RequestContext, code int, message string) {
	c.JSON(code, map[string]interface{}{
		"code":    code,
		"message": message,
	})
}
