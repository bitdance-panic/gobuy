package main

import (
	"context"
	payment "gobuy/app/rpc/kitex_gen/payment"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// CreatePayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, req *payment.CreatePaymentRequest) (resp *payment.CreatePaymentResponse, err error) {
	// TODO: Your code here...
	return
}

// GetPayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (resp *payment.GetPaymentResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdatePayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) UpdatePayment(ctx context.Context, req *payment.UpdatePaymentRequest) (resp *payment.UpdatePaymentResponse, err error) {
	// TODO: Your code here...
	return
}

// DeletePayment implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) DeletePayment(ctx context.Context, req *payment.DeletePaymentRequest) (resp *payment.DeletePaymentResponse, err error) {
	// TODO: Your code here...
	return
}
