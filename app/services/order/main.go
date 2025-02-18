package main

import (
	"github.com/bitdance-panic/gobuy/app/models"
	rpc_order "github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal"
	"github.com/bitdance-panic/gobuy/app/services/order/conf"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/kitex/server"
)

type OrderServiceAdapter struct {
	orderService *OrderServiceImpl
}

func kitexInit() (opts []server.Option) {
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		localIp := "0.0.0.0"
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))
	return
}
func main() {
	dal.Init()
	opts := kitexInit()

	svr := orderservice.NewServer(new(OrderServiceImpl), opts...)
	orderservice := &OrderServiceImpl{}
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
	//创建适配器
	adapter := NewOrderServiceAdapter(orderservice)
	//初始化Gin路由
	router := gin.Default()
	//注册订单相关路由
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/", adapter.CreateOrderAdapter)
		//orderGroup.GET("/:order_id", GetOrderAdapter)
		orderGroup.GET("/:user_id", adapter.GetOrdersByUserIDAdapter)
		orderGroup.PUT("/:order_id/status", adapter.UpdateOrderAdapter)
		orderGroup.DELETE("/:order_id", adapter.SoftDeleteOrderAdapter)
	}
	router.Run(":8882")

	// 创建一个新的cron实例
	c := cron.New()
	// 添加定时任务，每分钟执行一次
	_, err1 := c.AddFunc("@every 1m", checkAndUpdateExpiredOrders)
	if err1 != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	// 启动cron
	c.Start()
	// 保持主程序运行
	select {}
}

func NewOrderServiceAdapter(orderservice *OrderServiceImpl) *OrderServiceAdapter {
	return &OrderServiceAdapter{
		orderService: orderservice,
	}
}

// 处理创建订单的HTTP请求
func (a *OrderServiceAdapter) CreateOrderAdapter(context *gin.Context) {
	var req rpc_order.CreateOrderRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	resp, err := a.orderService.CreateOrder(context, &req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, resp)
}

// 通过订单id获取订单
// func (a *OrderServiceAdapter) GetOrderAdapter(context *gin.Context) {
//
// }
func (a *OrderServiceAdapter) GetOrdersByUserIDAdapter(context *gin.Context) {
	userIdStr := context.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	req := &rpc_order.GetUserOrdersRequest{
		UserId: int32(userId),
	}
	resp, err := a.orderService.GetUserOrders(context, req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, resp)
}
func (a *OrderServiceAdapter) UpdateOrderAdapter(context *gin.Context) {
	orderIDStr := context.Param("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}
	var req rpc_order.UpdateOrderRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req.OrderId = int32(orderID)
	resp, err := a.orderService.UpdateOrder(context, &req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, resp)
}
func (a *OrderServiceAdapter) SoftDeleteOrderAdapter(context *gin.Context) {
	orderIDStr := context.Param("order_id")
	_, err := strconv.Atoi(orderIDStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}
	status := a.orderService.SoftDeleteOrder(context)
	if status != http.StatusOK {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete order"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Order soft deleted successfully"})
}
func checkAndUpdateExpiredOrders() {
	//ctx := context.Background()
	//获取所有状态为“待支付”的订单
	orders, err := dal.GetOrdersByStatus(models.OrderStatusPending)
	if err != nil {
		log.Println("Failed to get pending orders:%v", err)
		return
	}
	//遍历订单，检查是否过期
	for _, order := range orders {
		if isOrderExpired(order) {
			//更新订单状态为"已取消"
			order.Status = int(models.OrderStatusCancelled)
			err, _ := dal.UpdateOrderStatus(int32(order.ID), int(order.Status))
			if err != nil {
				log.Println("Failed to update order:%v", err)
				continue
			}
			log.Println("Order %d has been cancelled due to expiration", order.ID)
		}
	}
}

// 检查订单是否过期
func isOrderExpired(order models.Order) bool {
	//假设订单创建后30分钟未支付则过期
	expirationTime := order.CreatedAt.Add(30 * time.Minute)
	return time.Now().After(expirationTime)
}
