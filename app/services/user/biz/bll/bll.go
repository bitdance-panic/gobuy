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

// Register 业务逻辑：注册新用户
func Register(ctx context.Context, req *rpc_user.RegisterReq) (*rpc_user.RegisterResp, error) {
	user, err := dao.RegisterUser(tidb.DB, ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return nil, err
	}
	return &rpc_user.RegisterResp{
		UserId:  int32(user.ID),
		Success: true,
	}, nil
}

func Login(ctx context.Context, req *rpc_user.LoginReq) (*rpc_user.LoginResp, error) {
	hlog.Infof("Login attempt for email=%s", req.Email)
	userPO, err := dao.GetUserByEmailAndPass(tidb.DB, ctx, req.Email, req.Password)
	resp := rpc_user.LoginResp{}
	// 没查到
	if err != nil {
		hlog.Errorf("Login failed for email=%s, error=%v. Invalid email or password", req.Email, err)
		resp.Success = false
	} else {
		resp.Success = true
		resp.UserId = int32(userPO.ID)
	}
	return &resp, err
}

// GetUsers 获取所有用户信息
func AdminListUser(ctx context.Context, req *rpc_user.AdminListUserReq) (*rpc_user.AdminListUserResp, error) {
	// 查询数据库中的所有用户信息，假设分页处理
	users, err := dao.AdminListUser(tidb.DB, ctx, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return &rpc_user.AdminListUserResp{
			Users:   nil,
			Message: "Failed to retrieve users",
		}, err
	}
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
	return &rpc_user.AdminListUserResp{
		Users:   userList,
		Message: "Users retrieved successfully",
	}, nil
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
		UserId:   int32(user.ID),
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

// 更新用户信息
func UpdateUser(ctx context.Context, req *rpc_user.UpdateUserReq) (*rpc_user.UpdateUserResp, error) {
	// 处理 nil 指针
	username := ""
	email := ""
	if req.Username != nil {
		username = *req.Username
	}
	if req.Email != nil {
		email = *req.Email
	}
	err := dao.UpdateUserByID(tidb.DB, ctx, int(req.UserId), username, email)
	if err != nil {
		return &rpc_user.UpdateUserResp{Success: false}, nil
	}
	return &rpc_user.UpdateUserResp{Success: true}, nil
}

// 删除用户
func RemoveUser(ctx context.Context, req *rpc_user.RemoveUserReq) (*rpc_user.RemoveUserResp, error) {
	err := dao.DeleteUserByID(tidb.DB, ctx, int(req.UserId))
	if err != nil {
		return &rpc_user.RemoveUserResp{Success: false}, nil
	}
	return &rpc_user.RemoveUserResp{Success: true}, nil
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
