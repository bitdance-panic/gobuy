package logic

import (
	"context"
	"errors"
	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/cart/internal/svc"
	"github.com/bitdance-panic/gobuy/app/services/cart/proto/cart"
	"gorm.io/gorm"
)

type AddItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddItemLogic {
	return &AddItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddItemLogic) AddItem(in *cart.AddItemReq) (*cart.AddItemResp, error) {
	db := l.svcCtx.DB
	p := models.Product{}
	u := models.User{}

	err := db.Preload("Cart").Where("id = ?", in.UserId).Take(&u).Error
	if err != nil {
		return nil, err
	}

	err = db.Where("id = ?", in.Item.ProductId).Take(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	c := models.CartItem{CartID: uint(u.Cart.ID), ProductID: uint(in.Item.ProductId), Quantity: uint(in.Item.Quantity)}
	err = db.Save(&c).Error
	if err != nil {
		return nil, err
	}

	return &cart.AddItemResp{}, nil

}
