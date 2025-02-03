namespace go order
include "product.thrift"
//定义订单状态枚举
enum OrderStatus{
    PENDING = 1  //待处理
    COMPLETED = 2 //已完成
    CANCELED = 3 //已取消
    SHIPPED =4 //已发货
}
//复制过来的商品信息
//struct Product {
//    1: i32 id;  // 来自 Base 的 ID
//    2: string name;
//    3: double price;
//    4: i32 stock;
//    5: string img;  // 对应 Go 结构体中的 Image
//    6: string description;  // 对应 Go 结构体中的 Description
//    7: string created_at;  // 来自 Base 的 CreatedAt，格式为字符串
//    8: string updated_at;  // 来自 Base 的 UpdatedAt，格式为字符串
//    9: bool is_deleted;  // 来自 Base 的 IsDeleted
//}
//订单
struct Order{
    1: i32 id; //来自Base的ID  订单id
    2: i32 user_id;//用户ID
    3: string order_number;//订单号
    4: double total_amount;//订单总金额
    5: OrderStatus status;//订单状态
    6: list<OrderItem> items;//订单项列表
    7: string create_at;//
    8: string update_at;//
    9: bool is_deleted;
}
//订单项
struct OrderItem{
    1: i32 order_id;//订单ID
    2: i32 product_id; //商品ID
    3: i32 quantity;//商品数量
    4: double price; //商品单价
    5: string product_name;//商品名称
    6: product.Product product;//关联商品
}
service OrderService{
    CreateOrderResponse createOrder(1: CreateOrderRequest req);
    UpdateOrderResponse updateOrder(1: UpdateOrderRequest req);
    GetUserOrdersResponse getUserOrders(1: GetUserOrdersRequest req);

}
//创建订单请求
struct CreateOrderRequest{
    1: i32 user_id;//用户id
    2: list<OrderItem> items;//订单项列表
}
struct CreateOrderResponse{
    1: Order order;
}
//更新订单请求
struct UpdateOrderRequest{
    1: OrderStatus status;//更新的订单状态
}
struct UpdateOrderResponse{
    1: Order order;
}
//获取用户订单请求
struct GetUserOrdersRequest{
    1: i32 user_id;//用户ID
}
//获取用户订单响应
struct GetUserOrdersResponse{
    1: list<Order> orders;//用户订单列表
}













