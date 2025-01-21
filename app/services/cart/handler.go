package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"app/services/cart/biz/bll"
)

var cartBLL *bll.CartBLL

// GetCartHandler 获取购物车信息的Handler
func GetCartHandler(ctx context.Context, c *app.RequestContext) {
	userID := c.QueryInt("user_id")
	if userID == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "user_id is required"})
		return
	}

	cart, err := cartBLL.GetCartByUserID(uint(userID))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, cart)
}

// AddItemToCartHandler 向购物车中添加商品的Handler
func AddItemToCartHandler(ctx context.Context, c *app.RequestContext) {
	userID := c.QueryInt("user_id")
	if userID == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "user_id is required"})
		return
	}

	productID := c.QueryInt("product_id")
	if productID == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "product_id is required"})
		return
	}

	quantity := c.QueryInt("quantity")
	if quantity == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "quantity is required"})
		return
	}

	err := cartBLL.AddItemToCart(uint(userID), uint(productID), quantity)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, utils.H{"message": "Item added to cart successfully"})
}