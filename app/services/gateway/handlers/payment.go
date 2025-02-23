package handlers

// // HandleCreatePayment 创建支付记录
// // @Summary 创建支付记录
// // @Description 创建一个新的支付记录并返回相关信息
// // @Accept application/json
// // @Produce application/json
// // @Router /payment [post]
// func HandleCreatePayment(ctx context.Context, c *app.RequestContext) {
// 	req := rpc_payment.CreatePaymentRequest{}

// 	// 从请求体中绑定参数并验证
// 	if err := c.BindAndValidate(&req); err != nil {
// 		hlog.Warnf("CreatePayment failed , validation error: %v", err)
// 		utils.Fail(c, err.Error())
// 		return
// 	}

// 	resp, err := services.PaymentService.CreatePayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
// 	if err != nil {
// 		utils.FailFull(c, consts.StatusInternalServerError, fmt.Sprintf("Create payment failed: %v", err.Error()), nil)
// 		return
// 	} else {
// 		utils.Success(c, utils.H{"payment": resp.Payment})
// 	}
// }

// // HandleGetPayment 获取支付记录
// // @Summary 获取支付记录
// // @Description 根据支付记录ID获取支付记录详情
// // @Accept application/json
// // @Produce application/json
// // @Router /payment/{id} [get]
// func HandleGetPayment(ctx context.Context, c *app.RequestContext) {
// 	req := rpc_payment.GetPaymentRequest{}

// 	// 从请求体中绑定参数并验证
// 	if err := c.BindAndValidate(&req); err != nil {
// 		hlog.Warnf("GetPayment failed , validation error: %v", err)
// 		utils.Fail(c, err.Error())
// 		return
// 	}

// 	resp, err := services.PaymentService.GetPayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
// 	if err != nil {
// 		utils.FailFull(c, consts.StatusNotFound, fmt.Sprintf("Payment not found: %v", err.Error()), nil)
// 		return
// 	} else {
// 		utils.Success(c, utils.H{"payment": resp.Payment})
// 	}
// }

// // HandleUpdatePayment 更新支付记录
// // @Summary 更新支付记录
// // @Description 更新指定支付记录的状态或其他信息
// // @Accept application/json
// // @Produce application/json
// // @Router /payment/{id} [put]
// func HandleUpdatePayment(ctx context.Context, c *app.RequestContext) {
// 	req := rpc_payment.UpdatePaymentRequest{}

// 	// 从请求体中绑定参数并验证
// 	if err := c.BindAndValidate(&req); err != nil {
// 		hlog.Warnf("UpdatePayment failed , validation error: %v", err)
// 		utils.Fail(c, err.Error())
// 		return
// 	}

// 	resp, err := services.PaymentService.UpdatePayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
// 	if err != nil {
// 		utils.FailFull(c, consts.StatusInternalServerError, fmt.Sprintf("UpdatePayment failed: %v", err.Error()), nil)
// 		return
// 	} else {
// 		utils.Success(c, utils.H{"payment": resp.Payment})
// 	}
// }

// // HandleDeletePayment 删除支付记录
// // @Summary 删除支付记录
// // @Description 删除指定ID的支付记录
// // @Accept application/json
// // @Produce application/json
// // @Router /payment/{id} [delete]
// func HandleDeletePayment(ctx context.Context, c *app.RequestContext) {
// 	req := rpc_payment.DeletePaymentRequest{}

// 	// 从请求体中绑定参数并验证
// 	if err := c.BindAndValidate(&req); err != nil {
// 		hlog.Warnf("DeletePayment failed , validation error: %v", err)
// 		utils.Fail(c, err.Error())
// 		return
// 	}

// 	_, err := services.PaymentService.DeletePayment(context.Background(), &req, callopt.WithRPCTimeout(3*time.Second))
// 	if err != nil {
// 		utils.FailFull(c, consts.StatusInternalServerError, fmt.Sprintf("DeletePayment failed: %v", err.Error()), nil)
// 		return
// 	} else {
// 		utils.Success(c, utils.H{"payment_id": req.PaymentId})
// 	}
// }
