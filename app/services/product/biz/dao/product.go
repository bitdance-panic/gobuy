package dao

import (
	"errors"

	"github.com/bitdance-panic/gobuy/app/services/product/biz/dal/tidb"

	"github.com/bitdance-panic/gobuy/app/models"
)

type Product = models.Product

func Create(product *Product) error {
	return tidb.DB.Create(product).Error
}

func GetByID(id uint) (*Product, error) {
	var product Product
	if err := tidb.DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func Update(product *Product) error {
	return tidb.DB.Save(product).Error
}

func Delete(id uint) error {
	result := tidb.DB.Delete(&Product{}, id)
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return result.Error
}

func Search(query string) ([]Product, error) {
	var products []Product
	searchQuery := "%" + query + "%"
	if err := tidb.DB.Where("name LIKE ? OR description LIKE ?", searchQuery, searchQuery).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
