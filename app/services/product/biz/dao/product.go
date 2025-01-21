package dao

import (
	"errors"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/product/biz/dal/postgres"
)

type Product = models.Product

func Create(product *Product) error {
	return postgres.DB.Create(product).Error
}

func GetByID(id uint) (*Product, error) {
	var product Product
	if err := postgres.DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func Update(product *Product) error {
	return postgres.DB.Save(product).Error
}

func Delete(id uint) error {
	result := postgres.DB.Delete(&Product{}, id)
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return result.Error
}

func Search(query string) ([]Product, error) {
	var products []Product
	searchQuery := "%" + query + "%"
	if err := postgres.DB.Where("name LIKE ? OR description LIKE ?", searchQuery, searchQuery).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
