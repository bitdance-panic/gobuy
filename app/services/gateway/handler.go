package main

import (
	"context"
	"net/http"
	"time"

	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"
	"github.com/bitdance-panic/gobuy/app/services/gateway/middleware"
	"github.com/bitdance-panic/gobuy/app/services/gateway/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	cutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
)

func PongHandler(ctx context.Context, c *app.RequestContext) {
	hlog.Info("Pong request received ")
	c.JSON(consts.StatusOK, cutils.H{"message": "pong"})
}

func handleLogin(ctx context.Context, c *app.RequestContext) {
	// 通过 /login?email=1234&pass=1234 测试
	email := c.Query("email")
	password := c.Query("pass")

	hlog.Infof("Login attempt for email=%s", email)

	req := rpc_user.LoginReq{
		Email:    email,
		Password: password,
	}
	resp, err := userservice.Login(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("Login failed for email=%s, error=%v", email, err)
		c.JSON(consts.StatusInternalServerError, cutils.H{"message": err.Error()})
		return
	}
	if resp.Success {
		hlog.Infof("Login successful for email=%s, userID=%d", email, resp.UserId)
		c.JSON(consts.StatusOK, cutils.H{"userID": resp.UserId})
	} else {
		hlog.Warnf("Login failed for email=%s: invalid credentials", email)
		c.JSON(consts.StatusOK, cutils.H{"message": "Invalid credentials"})
	}
}

// 成功响应构造器
func SuccessResponse(data interface{}) cutils.H {
	return cutils.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	}
}

// 错误响应构造器
func ErrorResponse(message string) cutils.H {
	return cutils.H{
		"code":    http.StatusInternalServerError,
		"message": message,
		"data":    nil,
	}
}

// 用户注册请求结构
type RegisterRequest struct {
	Username string `json:"username" vd:"len($)>5 && len($)<20"` // 用户名长度校验
	Password string `json:"password" vd:"len($)>6"`              // 密码复杂度校验
	Email    string `json:"email"`
}

// 用户注册处理函数
func RegisterHandler(ctx context.Context, c *app.RequestContext) {
	var req RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Errorf("Registration failed for username: %s, validation error: %v", req.Username, err)
		c.JSON(http.StatusBadRequest, ErrorResponse("Parameter validation failed"))
		return
	}

	// 检查用户是否已存在
	_, err := dao.GetUserByEmail(req.Email)
	if err == nil {
		hlog.Warnf("Registration failed for email: %s, email already exists", req.Email)
		c.JSON(http.StatusConflict, ErrorResponse("email already exists"))
		return
	}

	// 创建用户
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		hlog.Errorf("Password hashing failed for username: %s, error: %v", req.Username, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Internal Server Error"))
		return
	}

	user := &dao.User{
		Username:       req.Username,
		PasswordHashed: hashedPassword,
		Email:          req.Email,
	}
	if err := dao.CreateUser(user); err != nil {
		hlog.Errorf("User creation failed for username: %s, error: %v", req.Username, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("User registration failed"))
		return
	}

	hlog.Infof("User created successfully, user ID: %d", user.ID)
	c.JSON(http.StatusCreated, SuccessResponse(map[string]interface{}{
		"user_id": user.ID,
	}))
}

// 用户登录请求结构
type LoginRequest struct {
	Email    string `json:"email" vd:"required"`
	Password string `json:"password" vd:"required"`
}

// 用户登录处理函数
func LoginHandler(ctx context.Context, c *app.RequestContext) {
	var req LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("Login failed for email: %s, validation error: %v", req.Email, err)
		c.JSON(http.StatusBadRequest, ErrorResponse("Parameter validation failed"))
		return
	}

	// 获取用户信息
	user, err := dao.GetUserByEmail(req.Email)
	if err != nil {
		hlog.Warnf("Login failed for email: %s, user not found", req.Email)
		c.JSON(http.StatusUnauthorized, ErrorResponse("Invalid email or password"))
		return
	}

	// 验证密码
	if !utils.VerifyPassword(req.Password, user.PasswordHashed) {
		hlog.Warnf("Login failed for email: %s, incorrect password", req.Email)
		c.JSON(http.StatusUnauthorized, ErrorResponse("Invalid email or password"))
		return
	}

	// 生成双令牌
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, conf.GetConf().Hertz.JWTSecret)
	if err != nil {
		hlog.Errorf("Token generation failed for email: %s, error: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Internal Server Error"))
		return
	}

	// 保存刷新令牌
	if err := dao.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		hlog.Errorf("Failed to save refresh token for email: %s, error: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Internal Server Error"))
		return
	}

	hlog.Infof("Login successful for email: %s, generated tokens", req.Email)
	c.JSON(http.StatusOK, SuccessResponse(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}

// 用户信息响应结构
type ProfileResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	JoinTime string `json:"join_time"`
}

// 获取用户信息处理函数
func GetProfileHandler(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		hlog.Warn("Unauthorized access attempt, user ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse("Unauthorized access"))
		return
	}

	// 查询用户信息
	user, err := dao.GetUserByID(userID.(int))
	if err != nil {
		hlog.Errorf("Failed to query user profile for user ID: %d, error: %v", userID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Internal Server Error"))
		return
	}

	hlog.Infof("Successfully retrieved profile for user ID: %d", user.ID)
	c.JSON(http.StatusOK, SuccessResponse(ProfileResponse{
		UserID:   user.ID,
		Username: user.Username,
		JoinTime: user.CreatedAt.Format(time.RFC3339),
	}))
}

// 商品查询参数
type ListProductsRequest struct {
	Page     int `query:"page" vd:"$>=0"`
	PageSize int `query:"page_size" vd:"$>0 && $<=100"`
}

// 商品列表响应
type ProductItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// 获取商品列表处理函数
func ListProductsHandler(ctx context.Context, c *app.RequestContext) {
	var req ListProductsRequest
	if err := c.BindAndValidate(&req); err != nil {
		hlog.Warnf("Failed to list products, validation error: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse("Parameter validation failed"))
		return
	}

	products, total, err := dao.ListProducts(req.Page, req.PageSize)
	if err != nil {
		hlog.Errorf("Failed to query products, error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("Internal Server Error"))
		return
	}

	// 构造响应数据
	var items []ProductItem
	for _, p := range products {
		items = append(items, ProductItem{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	hlog.Infof("Successfully listed products, total: %d", total)
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"items": items,
		"total": total,
	}))
}

// 路由注册函数
func RegisterRoutes(h *server.Hertz) {

	// 公共路由
	public := h.Group("/")
	{
		public.POST("/register", RegisterHandler)
		public.POST("/login", LoginHandler)
		public.GET("/ping", PongHandler)
	}

	// 需要认证的路由
	auth := h.Group("/")
	auth.Use(middleware.AuthMiddleware.MiddlewareFunc())
	{
		auth.GET("/profile", GetProfileHandler)
		auth.GET("/products", ListProductsHandler)
	}
}
