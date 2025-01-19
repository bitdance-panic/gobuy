package consts

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusShipped   OrderStatus = "shipped"   // 已发货
	OrderStatusCompleted OrderStatus = "completed" // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
)
