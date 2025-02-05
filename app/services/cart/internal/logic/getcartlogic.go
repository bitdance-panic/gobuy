package logic

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/cart/internal/svc"
	"github.com/bitdance-panic/gobuy/app/services/cart/proto/cart"
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

func (l *GetCartLogic) GetCart(in *cart.GetCartReq) (*cart.GetCartResp, error) {
	db := l.svcCtx.DB
	u := models.User{}

	err := db.Preload("Cart").First(&u, in.UserId).Error
	if err != nil {
		return nil, err
	}

	c := make([]models.Product, 0)
	err = db.Model(&models.CartItem{}).Where("cart_id = ?", u.Cart.ID).Find(&c).Error
	if err != nil {
		return nil, err
	}

	res := cart.GetCartResp{Items: make([]*cart.CartItem, len(c))}
	for i, v := range c {
		res.Items[i] = new(cart.CartItem)
		res.Items[i].ProductId = uint32(v.ProductID)
		res.Items[i].Quantity = uint32(v.Quantity)
	}
	return &res, nil

}
