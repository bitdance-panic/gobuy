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
		resp.UserId = int32(userPO.ID)
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
		UserId:   strconv.Itoa(user.ID), //将 user.ID 转换为 string
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

// 封禁用户
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
