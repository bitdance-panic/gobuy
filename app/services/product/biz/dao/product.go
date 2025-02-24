package dao

import (
	"errors"

	"gorm.io/gorm"

	"github.com/bitdance-panic/gobuy/app/models"
)

type Product = models.Product

func Create(db *gorm.DB, product *Product) error {
	return db.Create(product).Error
}

func List(db *gorm.DB, pageNum int, pageSize int) (*[]Product, int64, error) {
	var products []Product
	if err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Where("is_deleted = false").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&Product{}).Where("is_deleted = false").Count(&count)
	return &products, count, nil
}

func AdminList(db *gorm.DB, pageNum int, pageSize int) (*[]Product, int64, error) {
	var products []Product
	if err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&Product{}).Count(&count)
	return &products, count, nil
}

func GetByID(db *gorm.DB, id int) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 处理记录未找到的情况
			return nil, nil
		}
		// 处理其他错误
		return nil, err
	}
	return &product, nil
}

func Update(db *gorm.DB, product *Product) error {
	if product == nil {
		return errors.New("product is nil")
	}
	return db.Save(product).Error
}

func Remove(db *gorm.DB, id int) error {
	result := db.Model(&Product{}).Where("id = ? AND is_deleted = false", id).Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}
	return nil
}

func Search(db *gorm.DB, query string, pageNum int, pageSize int) ([]Product, int64, error) {
	var products []Product
	searchQuery := "%" + query + "%"
	if err := db.Limit(pageSize).Offset((pageNum-1)*pageSize).Where("is_deleted = ? AND (name LIKE ? OR description LIKE ?)", false, searchQuery, searchQuery).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	db.Model(&Product{}).Where("is_deleted = ? AND (name LIKE ? OR description LIKE ?)", false, searchQuery, searchQuery).Count(&count)
	return products, count, nil
}
