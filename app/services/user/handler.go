package main

import (
	"context"
	"strconv"

	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/bll"
)

var ub *bll.UserBLL

func init() {
	ub = bll.NewUserBLL()
}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *rpc_user.RegisterReq) (resp *rpc_user.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *rpc_user.LoginReq) (resp *rpc_user.LoginResp, err error) {
	// 做validate...
	resp, err = ub.Login(ctx, req)
	return
}

// 获取用户信息
func (s *UserServiceImpl) GetUser(ctx context.Context, req *rpc_user.GetUserReq) (*rpc_user.GetUserResp, error) {
	userID, err := strconv.Atoi(req.UserId)
	if err != nil {
		return &rpc_user.GetUserResp{Success: false}, nil
	}
	return bll.GetUser(ctx, userID)
}

// 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *rpc_user.UpdateUserReq) (*rpc_user.UpdateUserResp, error) {
	return bll.UpdateUser(ctx, req)
}

// 封禁用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *rpc_user.DeleteUserReq) (*rpc_user.DeleteUserResp, error) {
	return bll.DeleteUser(ctx, req)
}
