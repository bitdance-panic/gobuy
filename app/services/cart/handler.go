package main

import (
	"context"
	"github.com/bitdance-panic/gobuy/app/services/cart/biz/bll"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var cartBLL *bll.CartBLL

// 获取购物车信息的Handler
func GetCartHandler(ctx context.Context, c *app.RequestContext) {
	userID := c.QueryInt("user_id")
	if userID == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "需要提供user_id(用户标识)"})
		return
	}

	cart, err := cartBLL.GetCartByUserID(uint(userID))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, cart)
}

//向购物车中添加商品的Handler
func AddItemToCartHandler(ctx context.Context, c *app.RequestContext) {
	userID := c.QueryInt("user_id")
	if userID == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "需要提供user_id(用户标识)"})
		return
	}

	productID := c.QueryInt("product_id")
	if productID == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "需要提供product_id(产品编号)"})
		return
	}

	quantity := c.QueryInt("quantity")
	if quantity == 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "需要数量"})
		return
	}

	err := cartBLL.AddItemToCart(uint(userID), uint(productID), quantity)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"message": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, utils.H{"message": "商品已成功添加到购物车"})
}
