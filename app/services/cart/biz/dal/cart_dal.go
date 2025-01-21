package dal
//数据访问层

import (
"app/models"
"gorm.io/gorm"
)

// CartDAL 定义购物车数据访问层结构体
type CartDAL struct {
	db *gorm.DB
}

// NewCartDAL 创建新的购物车数据访问层实例
func NewCartDAL(db *gorm.DB) *CartDAL {
	return &CartDAL{
		db: db,
	}
}

// GetCartByUserID 根据用户ID获取购物车信息
func (dal *CartDAL) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := dal.db.Where("user_id =?", userID).Preload("Items.Product").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// AddItemToCart 向购物车中添加商品
func (dal *CartDAL) AddItemToCart(userID, productID uint, quantity int) error {
	var cart models.Cart
	err := dal.db.Where("user_id =?", userID).First(&cart).Error
	if err != nil {
		return err
	}

	// 检查商品是否已在购物车中
	for _, item := range cart.Items {
		if item.ProductID == productID {
			item.Quantity += quantity
			return dal.db.Save(&item).Error
		}
	}

	// 商品不在购物车中，添加新的购物车项
	newItem := models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return dal.db.Create(&newItem).Error
}
