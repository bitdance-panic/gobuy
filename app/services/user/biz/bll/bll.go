package bll

import (
	// "common/model"
	"context"

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
