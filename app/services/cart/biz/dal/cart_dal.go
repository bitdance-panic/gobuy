package dal

// 数据访问层

import (
	"github.com/bitdance-panic/gobuy/app/models"
	"gorm.io/gorm"
)

type CartDAL struct {
	db *gorm.DB
}

func NewCartDAL(db *gorm.DB) *CartDAL {
	return &CartDAL{
		db: db,
	}
}

func (dal *CartDAL) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	// 先确保用户的购物车存在且未被删除
	err := dal.db.Where("user_id =? AND is_deleted =?", userID, false).Preload("Items.Product").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// 向购物车中添加商品
func (dal *CartDAL) AddItemToCart(userID, productID uint, quantity int) error {
	var cart models.Cart
	// 先确保用户的购物车存在且未被删除
	err := dal.db.Where("user_id =? AND is_deleted =?", userID, false).First(&cart).Error
	if err != nil {
		return err
	}

	// 检查商品是否已在购物车中且未被删除
	for _, item := range cart.Items {
		if item.ProductID == productID && !item.IsDeleted {
			item.Quantity += uint(quantity)
			return dal.db.Save(&item).Error
		}
	}

	// 商品不在购物车中，添加新的购物车项
	newItem := models.CartItem{
		UserID:    userID,
		CartID:    uint(cart.ID),
		ProductID: productID,
		Quantity:  uint(quantity),
		IsDeleted: false,
	}
	return dal.db.Create(&newItem).Error
}

// 从购物车中移除商品
func (dal *CartDAL) RemoveItemFromCart(userID, productID uint) error {
	var cart models.Cart
	// 先确保用户的购物车存在且未被删除
	err := dal.db.Where("user_id =? AND is_deleted =?", userID, false).First(&cart).Error
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		if item.ProductID == productID && !item.IsDeleted {
			item.IsDeleted = true
			return dal.db.Save(&item).Error
		}
	}

	return nil
}

// 清空购物车
func (dal *CartDAL) ClearCart(userID uint) error {
	var cart models.Cart
	// 先确保用户的购物车存在且未被删除
	err := dal.db.Where("user_id =? AND is_deleted =?", userID, false).First(&cart).Error
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		if !item.IsDeleted {
			item.IsDeleted = true
			if err := dal.db.Save(&item).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
