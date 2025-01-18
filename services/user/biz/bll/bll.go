package bll

import (
	"context"
	"user/biz/dal/mysql"
	"user/biz/model"

	rpc_user "rpc/kitex_gen/user"
	// "golang.org/x/crypto/bcrypt"
)

type UserBLL struct{}

func NewUserBLL() *UserBLL {
	return &UserBLL{}
}

func (s *UserBLL) Login(ctx context.Context, req *rpc_user.LoginReq) (*rpc_user.LoginResp, error) {
	userPO, err := model.GetByEmailAndPass(mysql.DB, ctx, req.Email, req.Password)
	resp := rpc_user.LoginResp{}
	// 没查到
	if err != nil {
		resp.Success = false
	} else {
		resp.Success = true
		resp.UserId = int32(userPO.ID)
	}
	return &resp, nil
}
