package main

import (
	"context"
	user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// TODO: Your code here...
	return
}

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserReq) (resp *user.GetUserResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserReq) (resp *user.UpdateUserResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserReq) (resp *user.DeleteUserResp, err error) {
	// TODO: Your code here...
	return
}

// BlockUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) BlockUser(ctx context.Context, req *user.BlockUserReq) (resp *user.BlockUserResp, err error) {
	// TODO: Your code here...
	return
}

// UnblockUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnblockUser(ctx context.Context, req *user.UnblockUserReq) (resp *user.UnblockUserResp, err error) {
	// TODO: Your code here...
	return
}
