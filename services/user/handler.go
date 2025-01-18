package main

import (
	"context"
	"fmt"
	user "rpc/kitex_gen/user"
	"user/biz/bll"
)

var ub *bll.UserBLL

func init() {
	ub = bll.NewUserBLL()
}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	fmt.Printf("Login called, req: %+v\n", *req)
	// 做validate...
	resp, err = ub.Login(ctx, req)
	return
}
