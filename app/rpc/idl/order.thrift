namespace go order

struct OrderProductItem{
    1: i32 product_id;
    4: i32 quantity;
}

struct OrderItem{
    1: i32 order_id;
    2: i32 product_id;
    3: string product_name;
    4: i32 quantity;
    5: double price;
}

struct Order{
    1: i32 id;
    2: i32 user_id;
    // 订单号
    3: string number;
    4: double total_price;
    5: i32 status;
    6: list<OrderItem> items;
    7: string created_at;
    8: string pay_time;
}

service OrderService{
    CreateOrderResp createOrder(1: CreateOrderReq req);
    UpdateOrderStatusResp updateOrderStatus(1: UpdateOrderStatusReq req);
    GetOrderResp getOrder(1: GetOrderReq req);
    ListUserOrderResp listUserOrder(1: ListUserOrderReq req);
    ListUserOrderResp adminListOrder(1: ListUserOrderReq req);
}
struct CreateOrderReq{
    1: i32 user_id;
    2: list<OrderProductItem> items;
}
struct CreateOrderResp{
    1: Order order;
}
struct UpdateOrderStatusReq{
    1: i32 order_id;
    2: i32 status;
}
struct UpdateOrderStatusResp{
    1: i32 new_status;
    2: bool success;
}
struct ListUserOrderReq{
    1: i32 user_id;
    2: i32 page_num;
    3: i32 page_size;
}
struct ListUserOrderResp{
    1: list<Order> orders;
}
struct GetOrderReq{
    1: i32 order_id;
}
struct GetOrderResp{
    1: Order order;
}