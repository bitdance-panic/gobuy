package bll

import (
	"app/models"
	"app/services/cart/biz/dal"
)

// CartBLL 定义购物车业务逻辑层结构体
type CartBLL struct {
	cartDAL *dal.CartDAL
}

// NewCartBLL 创建新的购物车业务逻辑层实例
func NewCartBLL(cartDAL *dal.CartDAL) *CartBLL {
	return &CartBLL{
		cartDAL: cartDAL,
	}
}

// GetCartByUserID 根据用户ID获取购物车信息
func (bll *CartBLL) GetCartByUserID(userID uint) (*models.Cart, error) {
	return bll.cartDAL.GetCartByUserID(userID)
}

// AddItemToCart 向购物车中添加商品
func (bll *CartBLL) AddItemToCart(userID, productID uint, quantity int) error {
	return bll.cartDAL.AddItemToCart(userID, productID, quantity)
}
