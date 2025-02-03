package logic


import (
	"context"
	"errors"
	"gorm.io/gorm"
	"app/models"
	"strconv"

	"app/services/cart/internal/svc"
	"app/services/cart/proto/cart"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddItemLogic {
	return &AddItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddItemLogic) AddItem(in *cart.AddItemReq) (*cart.AddItemResp, error) {
	db := l.svcCtx.DB
	log := l.svcCtx.Log

	// 查询用户及其购物车信息
	var user model.User
	err := db.Preload("Cart").Where("id = ?", in.UserId).Take(&user).Error
	if err != nil {
		log.Errorw("take user with cart", "error", err)
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 确保用户有购物车
	if user.Cart == nil {
		return nil, errors.New("用户没有购物车")
	}

	// 查询产品信息
	var product model.Product
	err = db.Where("id = ?", in.Item.ProductId).Take(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warnw("AddItem not found product", "product_id", in.Item.ProductId)
		return nil, fmt.Errorf("产品不存在: %d", in.Item.ProductId)
	} else if err != nil {
		log.Errorw("AddItem", "error", err)
		return nil, fmt.Errorf("查询产品信息失败: %v", err)
	}

	// 检查是否已存在该商品
	var cartProduct model.CartProducts
	err = db.Where("cart_id = ? AND product_id = ?", user.Cart.ID, in.Item.ProductId).First(&cartProduct).Error
	if err == nil {
		// 更新现有商品数量
		cartProduct.Quantity += uint(in.Item.Quantity)
		err = db.Save(&cartProduct).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新的购物车项
		cartProduct = model.CartProducts{
			CartID:    user.Cart.ID,
			ProductID: uint(in.Item.ProductId),
			Quantity:  uint(in.Item.Quantity),
		}
		err = db.Create(&cartProduct).Error
	} else {
		log.Errorw("查询购物车项失败", "error", err)
		return nil, fmt.Errorf("查询购物车项失败: %v", err)
	}

	if err != nil {
		log.Errorw("save cart", "error", err)
		return nil, fmt.Errorf("保存购物车项失败: %v", err)
	}

	return &cart.AddItemResp{}, nil
}
