package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/bitdance-panic/gobuy/app/models"
	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"

	//"github.com/bitdance-panic/gobuy/app/services/gateway/conf"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	gutils "github.com/bitdance-panic/gobuy/app/services/gateway/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/hertz-contrib/jwt"
)

type User = models.User

var identityKey = "uid"

var secret = "panic-bitdance"

var UserClient userservice.Client

const (
	AccessTokenExpire  = 12 * time.Hour
	RefreshTokenExpire = 7 * 24 * time.Hour
)

var AuthMiddleware *jwt.HertzJWTMiddleware

// BlackList      = gutils.NewSyncSet()

func InitAuth() {
	// JWT中间件配置
	initJWTMiddleware()
}

func initJWTMiddleware() {
	var err error
	AuthMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:            "gobuy auth",
		SigningAlgorithm: "HS256",
		Key:              []byte(secret),
		Timeout:          AccessTokenExpire,
		MaxRefresh:       RefreshTokenExpire,
		// 从请求头的Authorization提取token
		TokenLookup:                 "header: Authorization",
		TokenHeadName:               "Bearer",
		SendAuthorization:           true,
		WithoutDefaultTokenHeadName: false,
		TimeFunc:                    time.Now,
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			return e.Error()
		},
		DisabledAbort: false,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			userID := int(data.(int32))

			return jwt.MapClaims{
				identityKey: userID,
			}
		},
		Authenticator: authenticate,
		LoginResponse: loginResponse,
		// Authorizator: authorize,
		Unauthorized: unauthorizedHandler,
		// 登出的响应
		// 登出会自动删除token
		LogoutResponse: func(ctx context.Context, c *app.RequestContext, code int) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": http.StatusOK,
			})
		},
		RefreshResponse: refreshResponse,
		IdentityKey:     identityKey,
		IdentityHandler: identityHandler,
	})

	_ = AuthMiddleware
	if err != nil {
		hlog.Fatalf("auth middleware Error:" + err.Error())
	}
}

// 认证处理，用于用户登录时提取并验证登录凭据
func authenticate(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	loginReq := rpc_user.LoginReq{}

	if err := c.Bind(&loginReq); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	// 检查黑名单...

	loginResp, err := UserClient.Login(context.Background(), &loginReq, callopt.WithRPCTimeout(10*time.Second))
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	c.Set("uid", int(loginResp.UserId))
	return loginResp.UserId, nil
}

// 登录响应
func loginResponse(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
	// 生成RefreshToken
	refreshToken, err := gutils.GenerateRefreshToken(c.GetInt("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "刷新Token失败",
		})
		return
	}

	// 保存refreshToken
	if err := dao.UpdateRefreshToken(tidb.DB, c.GetInt("uid"), refreshToken); err != nil {
		hlog.Errorf("Failed to save refresh token, error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "保存refreshToken失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":          http.StatusOK,
		"access_token":  token,
		"expire":        expire.Unix(),
		"refresh_token": refreshToken,
	})
}

// 授权处理
// func authorize(data interface{}, ctx context.Context, c *app.RequestContext) bool {
// 	uid, ok := data.(int)
// 	if !ok {
// 		return false
// 	}

// refreshResponse 刷新Token响应
func refreshResponse(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
	// 生成新的RefreshToken
	refreshToken, err := gutils.GenerateRefreshToken(c.GetInt("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "刷新refreshToken失败",
		})
		return
	}

	// 更新数据库中的RefreshToken
	if err := dao.UpdateRefreshToken(tidb.DB, c.GetInt("uid"), refreshToken); err != nil {
		hlog.Errorf("Failed to save refresh token, error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "保存Token失败",
		})
		return
	}

	// 返回新的双Token
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":          http.StatusOK,
		"access_token":  token,
		"expire":        expire.Unix(),
		"refresh_token": refreshToken,
	})
}

// 身份处理
func identityHandler(ctx context.Context, c *app.RequestContext) interface{} {
	claims := jwt.ExtractClaims(ctx, c)
	return claims["roles"]
}

// 统一错误处理
func unauthorizedHandler(ctx context.Context, c *app.RequestContext, code int, message string) {
	c.JSON(code, map[string]interface{}{
		"code":    code,
		"message": message,
	})
	c.Abort()
}
