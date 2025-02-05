package logic

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/cart/internal/svc"
	"github.com/bitdance-panic/gobuy/app/services/cart/proto/cart"
)

type EmptyCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmptyCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmptyCartLogic {
	return &EmptyCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmptyCartLogic) EmptyCart(in *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	db := l.svcCtx.DB
	c := models.Cart{UserID: uint(in.UserId)}

	if err := db.Preload("Products").Where("user_id = ?", c.UserID).Take(&c).Error; err != nil {
		return nil, err
	}

	ProductID := make([]uint, len(c.ID))
	for i, v := range c.Products {
		ProductID[i] = v.ID
	}
	//log.Debug("empty cart_id:" + strconv.Itoa(int(c.ID)))

	err := db.Where("cart_id = ?", c.ID).Where("product_id in ?", ProductID).Delete(&models.CartItem{}).Error
	if err != nil {
		return nil, err
	}
	return &cart.EmptyCartResp{}, nil

}
