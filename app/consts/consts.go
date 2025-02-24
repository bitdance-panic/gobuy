package consts

// 定义订单状态常量
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 1 // 待支付
	OrderStatusPaid      OrderStatus = 2 // 已支付
	OrderStatusCancelled OrderStatus = 3 // 已取消
)

const (
	CONTEXT_UID_KEY       string = "uid"
	ITEM_INITIAL_QUANTITY int    = 1
)
