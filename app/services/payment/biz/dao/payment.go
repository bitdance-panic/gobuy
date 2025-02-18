package dao

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/payment/biz/dal/tidb"
)

// CreatePayment 插入新的 Payment 记录
func CreatePayment(ctx context.Context, p *models.Payment) error {
	return tidb.DB.WithContext(ctx).Create(p).Error
}

// GetPaymentByID 根据主键ID查询支付记录
func GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	var payment models.Payment
	if err := tidb.DB.WithContext(ctx).First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// UpdatePayment 更新支付记录
func UpdatePayment(ctx context.Context, p *models.Payment) error {
	return tidb.DB.WithContext(ctx).Save(p).Error
}

// DeletePayment 根据ID删除支付记录(物理删除)
func DeletePayment(ctx context.Context, id int) error {
	return tidb.DB.WithContext(ctx).Delete(&models.Payment{}, id).Error
}

// ListPaymentsByUserID 查询某个用户的所有支付记录
func ListPaymentsByUserID(ctx context.Context, userID int) ([]models.Payment, error) {
	var payments []models.Payment
	if err := tidb.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
