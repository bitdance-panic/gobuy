package logic

import (
	"context"
	"fmt"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/cart/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EmptyCartLogic) EmptyCart(in *models.EmptyCartReq) (*models.EmptyCartResp, error) {
	db := l.svcCtx.DB
	// 查询用户的购物车
	var cart models.Cart
	err := db.Preload("Products").Where("user_id = ?", in.UserId).Take(&cart).Error
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败: %v", err)
	}

	// 直接删除与该用户相关的所有购物车项
	err = db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
	if err != nil {
		return nil, fmt.Errorf("清空购物车失败: %v", err)
	}

	return &cart.EmptyCartResp{}, nil
}
