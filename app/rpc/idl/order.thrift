namespace go order

struct OrderItem{
    1: i32 order_id;
    2: i32 product_id;
    3: string product_name;
    4: i32 quantity;
    5: double price;
    6: string product_image;
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

struct OrderAddress{
    1: i32 id;
    2: i32 user_id;
    3: string phone;
    // 订单号
    4: i32 order_id;
    5: string order_address;
}

struct CreateOrderAddressReq{
    1: i32 user_id;
    2: string phone;
    3: i32 order_id;
    4: string order_address;
}

struct CreateOrderAddressResp{
    1: i32 order_id;
    2: bool success;
}

struct DeleteOrderAddressReq{
    1: i32 order_id;
}

struct DeleteOrderAddressResp{
    1: i32 order_id;
    2: bool success;
}

struct UpdateOrderAddressReq{
    1: i32 order_id;
    2: string order_address;
}

struct UpdateOrderAddressResp{
    1: string order_address;
    2: bool success;
}

struct GetOrderAddressReq{
    1: i32 order_id;
}

struct GetOrderAddressResp{
    1: OrderAddress order_address;
}

service OrderService{
    CreateOrderResp createOrder(1: CreateOrderReq req);
    UpdateOrderStatusResp updateOrderStatus(1: UpdateOrderStatusReq req);
    GetOrderResp getOrder(1: GetOrderReq req);
    ListOrderResp listUserOrder(1: ListOrderReq req);
    ListOrderResp adminListOrder(1: ListOrderReq req);
    CreateOrderAddressResp createOrderAddress(1: CreateOrderAddressReq req);
    DeleteOrderAddressResp deleteOrderAddress(1: DeleteOrderAddressReq req);
    UpdateOrderAddressResp updateOrderAddress(1: UpdateOrderAddressReq req);
    GetOrderAddressResp getOrderAddress(1: GetOrderAddressReq req);
}


struct CreateOrderReq{
    1: i32 user_id;
    2: list<i32> cartItemIDs;
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
struct ListOrderReq{
    1: i32 user_id;
    2: i32 page_num;
    3: i32 page_size;
}
struct ListOrderResp{
    1: list<Order> orders;
    2: i64 total_count
}
struct GetOrderReq{
    1: i32 order_id;
}
struct GetOrderResp{
    1: Order order;
}