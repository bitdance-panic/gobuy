package consts

// 定义订单状态常量
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 0 // 待支付
	OrderStatusPaid      OrderStatus = 1 // 已支付
	OrderStatusCancelled OrderStatus = 2 // 已取消
)

const (
	CONTEXT_UID_KEY       string = "uid"
	ITEM_INITIAL_QUANTITY int    = 1
)
