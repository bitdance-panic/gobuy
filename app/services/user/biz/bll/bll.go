package bll

import (
	// "common/model"
	"context"

	rpc_user "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dal/postgres"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dao"
	// "golang.org/x/crypto/bcrypt"
)

type UserBLL struct{}

func NewUserBLL() *UserBLL {
	return &UserBLL{}
}

func (s *UserBLL) Login(ctx context.Context, req *rpc_user.LoginReq) (*rpc_user.LoginResp, error) {
	userPO, err := dao.GetUserByEmailAndPass(postgres.DB, ctx, req.Email, req.Password)
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
