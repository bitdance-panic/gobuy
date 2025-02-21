package bll

import (
	// "common/model"
	"context"
	"strconv"

	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dao"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	// "golang.org/x/crypto/bcrypt"
)

type UserBLL struct{}

func NewUserBLL() *UserBLL {
	return &UserBLL{}
}

// Register 业务逻辑：注册新用户
func Register(ctx context.Context, req *rpc_user.RegisterReq) (*rpc_user.RegisterResp, error) {
	// 调用 DAO 层创建用户
	user, err := dao.RegisterUser(tidb.DB, ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return &rpc_user.RegisterResp{UserId: "", Success: false}, err
	}

	// 返回成功响应，包含新用户的 ID
	return &rpc_user.RegisterResp{
		UserId:  strconv.Itoa(user.ID), // 返回 user ID
		Success: true,
	}, nil
}

// GetUsers 获取所有用户信息
func GetUsers(ctx context.Context, req *rpc_user.GetUsersReq) (*rpc_user.GetUsersResp, error) {
	// 查询数据库中的所有用户信息，假设分页处理
	users, err := dao.GetUsers(tidb.DB, ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return &rpc_user.GetUsersResp{
			Users:   nil,
			Message: "Failed to retrieve users",
		}, err
	}

	// 将从数据库中查询到的用户数据转换为 *rpc_user.User 格式
	var userList []*rpc_user.User
	for _, u := range users {
		userList = append(userList, &rpc_user.User{
			UserId:       int32(u.ID), // 转换为 int32 类型
			Username:     u.Username,
			Email:        u.Email,
			RefreshToken: u.RefreshToken,
		})
	}

	// 返回响应
	return &rpc_user.GetUsersResp{
		Users:   userList, // 返回转换后的用户列表
		Message: "Users retrieved successfully",
	}, nil
}

func (s *UserBLL) Login(ctx context.Context, req *rpc_user.LoginReq) (*rpc_user.LoginResp, error) {
	hlog.Infof("Login attempt for email=%s", req.Email)

	userPO, err := dao.GetUserByEmailAndPass(tidb.DB, ctx, req.Email, req.Password)
	resp := rpc_user.LoginResp{}
	// 没查到
	if err != nil {
		hlog.Errorf("Login failed for email=%s, error=%v. Invalid email or password", req.Email, err)
		resp.Success = false
	} else {
		resp.Success = true
		resp.UserId = strconv.Itoa(userPO.ID)
	}
	return &resp, err
}

func GetUser(ctx context.Context, userID int) (*rpc_user.GetUserResp, error) {
	if userID <= 0 {
		return &rpc_user.GetUserResp{Success: false}, nil
	}

	user, err := dao.GetUserByID(tidb.DB, ctx, userID)
	if err != nil {
		return &rpc_user.GetUserResp{Success: false}, nil
	}

	return &rpc_user.GetUserResp{
		Success:  true,
		UserId:   strconv.Itoa(user.ID), // 将 user.ID 转换为 string
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

// 更新用户信息
func UpdateUser(ctx context.Context, req *rpc_user.UpdateUserReq) (*rpc_user.UpdateUserResp, error) {
	userID, err := strconv.Atoi(req.UserId)
	if err != nil || userID <= 0 {
		return &rpc_user.UpdateUserResp{Success: false}, nil
	}

	// 处理 nil 指针
	username := ""
	email := ""

	if req.Username != nil {
		username = *req.Username
	}
	if req.Email != nil {
		email = *req.Email
	}

	err = dao.UpdateUserByID(tidb.DB, ctx, userID, username, email)
	if err != nil {
		return &rpc_user.UpdateUserResp{Success: false}, nil
	}

	return &rpc_user.UpdateUserResp{Success: true}, nil
}

// 删除用户
func DeleteUser(ctx context.Context, req *rpc_user.DeleteUserReq) (*rpc_user.DeleteUserResp, error) {
	userID, err := strconv.Atoi(req.UserId)
	if err != nil || userID <= 0 {
		return &rpc_user.DeleteUserResp{Success: false}, nil
	}

	err = dao.DeleteUserByID(tidb.DB, ctx, userID)
	if err != nil {
		return &rpc_user.DeleteUserResp{Success: false}, nil
	}

	return &rpc_user.DeleteUserResp{Success: true}, nil
}

// 封禁用户 ：将用户加入黑名单
func BlockUser(ctx context.Context, req *rpc_user.BlockUserReq) (*rpc_user.BlockUserResp, error) {
	user, err := dao.BlockUser(tidb.DB, ctx, req.Identifier, req.Reason, req.ExpiresAt)
	if err != nil {
		return &rpc_user.BlockUserResp{BlockId: "", Success: false}, err
	}

	return &rpc_user.BlockUserResp{
		BlockId: strconv.Itoa(user.ID), // 返回 user ID
		Success: true,
	}, nil
}

// 解禁用户
func UnblockUser(ctx context.Context, req *rpc_user.UnblockUserReq) (*rpc_user.UnblockUserResp, error) {
	err := dao.UnblockUser(tidb.DB, ctx, req.Identifier)
	if err != nil {
		return &rpc_user.UnblockUserResp{Success: false}, err
	}

	return &rpc_user.UnblockUserResp{
		Success: true,
	}, nil
}
