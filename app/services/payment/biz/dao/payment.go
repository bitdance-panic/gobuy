package dao

import (
	"gobuy/app/models" // 调整为你项目实际 import 路径

	"gobuy/app/services/user/biz/dal/tidb"
)

// CreatePayment 插入新的 Payment 记录
func CreatePayment(p *models.Payment) error {
	return tidb.DB.Create(p).Error
}

// GetPaymentByID 根据主键ID查询支付记录
func GetPaymentByID(id int) (*models.Payment, error) {
	var payment models.Payment
	if err := tidb.DB.First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// UpdatePayment 更新支付记录
func UpdatePayment(p *models.Payment) error {
	return tidb.DB.Save(p).Error
}

// DeletePayment 根据ID删除支付记录(物理删除)
func DeletePayment(id int) error {
	return tidb.DB.Delete(&models.Payment{}, id).Error
}

// ListPaymentsByUserID 查询某个用户的所有支付记录
func ListPaymentsByUserID(userID int) ([]models.Payment, error) {
	var payments []models.Payment
	if err := tidb.DB.Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
