package svc

import (
	"lab-zero/orders/ordersclient"
	"lab-zero/users/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	OrdersRpc ordersclient.Orders
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		OrdersRpc: ordersclient.NewOrders(zrpc.MustNewClient(c.OrdersServiceRpc)),
	}
}
