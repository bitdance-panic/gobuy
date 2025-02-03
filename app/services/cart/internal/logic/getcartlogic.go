package logic

import (
"context"
"app/models"
"app/services/cart/internal/svc"
"app/services/cart/proto/cart"

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

func (l *EmptyCartLogic) EmptyCart(in *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log

	// 查询用户的购物车
	var cart model.Cart
	err := db.Preload("Products").Where("user_id = ?", in.UserId).Take(&cart).Error
	if err != nil {
		log.Errorw("获取购物车失败", "error", err, "user_id", in.UserId)
		return nil, fmt.Errorf("获取购物车失败: %v", err)
	}

	// 直接删除与该用户相关的所有购物车项
	err = db.Where("cart_id = ?", cart.ID).Delete(&model.CartProducts{}).Error
	if err != nil {
		log.Errorw("清空购物车失败", "error", err, "cart_id", cart.ID)
		return nil, fmt.Errorf("清空购物车失败: %v", err)
	}

	log.Infow("购物车已成功清空", "user_id", in.UserId, "cart_id", cart.ID)

	return &cart.EmptyCartResp{}, nil
}
