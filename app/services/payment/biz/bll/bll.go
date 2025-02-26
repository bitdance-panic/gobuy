package bll

// import (
// 	"context"

// 	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/payment"

// 	"github.com/bitdance-panic/gobuy/app/services/payment/biz/dao"

// 	"github.com/bitdance-panic/gobuy/app/models"
// )

// type PaymentBLL struct{}

// // NewPaymentBLL 创建 PaymentBLL 实例
// func NewPaymentBLL() *PaymentBLL {
// 	return &PaymentBLL{}
// }

// // CreatePayment 创建支付记录
// func (b *PaymentBLL) CreatePayment(ctx context.Context, req *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
// 	p := models.Payment{
// 		UserID:  int(req.UserId),
// 		OrderID: int(req.OrderId),
// 		Amount:  req.Amount,
// 		Status:  1,
// 	}

// 	if err := dao.CreatePayment(ctx, &p); err != nil {
// 		return nil, err
// 	}

// 	return &payment.CreatePaymentResponse{
// 		Payment: convertToProtoPayment(&p),
// 	}, nil
// }

// // UpdatePayment 更新支付记录
// func (b *PaymentBLL) UpdatePayment(ctx context.Context, req *payment.UpdatePaymentRequest) (*payment.UpdatePaymentResponse, error) {
// 	// 查询支付记录
// 	p, err := dao.GetPaymentByID(ctx, int(req.PaymentId))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 根据请求更新字段
// 	if req.Status != payment.PaymentStatus_PENDING {
// 		p.Status = int(req.Status)
// 	}

// 	if err := dao.UpdatePayment(ctx, p); err != nil {
// 		return nil, err
// 	}

// 	return &payment.UpdatePaymentResponse{
// 		Payment: convertToProtoPayment(p),
// 	}, nil
// }

// // DeletePayment 删除支付记录
// func (b *PaymentBLL) DeletePayment(ctx context.Context, req *payment.DeletePaymentRequest) (*payment.DeletePaymentResponse, error) {
// 	// 删除支付记录
// 	err := dao.DeletePayment(ctx, int(req.PaymentId))
// 	if err != nil {
// 		return &payment.DeletePaymentResponse{Success: false}, err
// 	}

// 	return &payment.DeletePaymentResponse{Success: true}, nil
// }

// // GetPayment 查询单个支付记录
// func (b *PaymentBLL) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
// 	// 根据 PaymentId 获取支付记录
// 	p, err := dao.GetPaymentByID(ctx, int(req.PaymentId))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &payment.GetPaymentResponse{
// 		Payment: convertToProtoPayment(p),
// 	}, nil
// }

// // convertToProtoPayment 将 models.Payment 转换为 Thrift 的 payment.Payment
// func convertToProtoPayment(p *models.Payment) *payment.Payment {
// 	return &payment.Payment{
// 		Id:        int32(p.ID),
// 		UserId:    int32(p.UserID),
// 		OrderId:   int32(p.OrderID),
// 		Amount:    p.Amount,
// 		Status:    payment.PaymentStatus(p.Status),
// 		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UpdatedAt: p.UpdatedAt.Format("2006-01-02 15:04:05"),
// 	}
// }
