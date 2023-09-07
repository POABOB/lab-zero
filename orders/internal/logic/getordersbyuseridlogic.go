package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"lab-zero/orders/internal/svc"
	"lab-zero/orders/orders"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrdersByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrdersByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrdersByUserIdLogic {
	// 建立假資料給user撈
	svcCtx.RedisCli.Setnx("orders.userid:1", `[{"buyer_id":1,"seller_id":4,"seller_account":123767,"good_id":73,"good_name":123,"good_price":6000,"status":1,"remark":"100","created_at":1693675667,"updated_at":1693675667},{"buyer_id":1,"seller_id":4,"seller_account":123767,"good_id":73,"good_name":123,"good_price":6000,"status":1,"remark":"100","created_at":1693675667,"updated_at":1693675667},{"buyer_id":1,"seller_id":4,"seller_account":123767,"good_id":73,"good_name":123,"good_price":6000,"status":1,"remark":"100","created_at":1693675667,"updated_at":1693675667},{"buyer_id":1,"seller_id":4,"seller_account":123767,"good_id":73,"good_name":123,"good_price":6000,"status":1,"remark":"100","created_at":1693675667,"updated_at":1693675667},{"buyer_id":1,"seller_id":4,"seller_account":123767,"good_id":73,"good_name":123,"good_price":6000,"status":1,"remark":"100","created_at":1693675667,"updated_at":1693675667},{"buyer_id":1,"seller_id":4,"seller_account":123767,"good_id":73,"good_name":123,"good_price":6000,"status":1,"remark":"100","created_at":1693675667,"updated_at":1693675667}]`)

	return &GetOrdersByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrdersByUserIdLogic) GetOrdersByUserId(in *orders.GetOrdersByUserIdRequest) (*orders.GetOrdersByUserIdResponse, error) {
	// validate
	if err := in.ValidateAll(); err != nil {
		return &orders.GetOrdersByUserIdResponse{Message: err.Error()}, nil
	}

	// redis
	val, err := l.svcCtx.RedisCli.Get(fmt.Sprintf("orders.userid:%v", in.UserId))
	// v, _ := l.svcCtx.RedisCli.Incr("no-single")
	// fmt.Println(v)

	// redis singleflight
	// val, err := l.svcCtx.SingleFlight.Do("get_orders_by_userid", func() (interface{}, error) {
	// 	val, err := l.svcCtx.RedisCli.Get(fmt.Sprintf("orders.userid:%v", in.UserId))
	// 	v, _ := l.svcCtx.RedisCli.Incr("single")
	// 	fmt.Println(v)
	// 	return val, err
	// })
	if err != nil {
		return nil, err
	}

	var data []*orders.UserOrders
	// json.Unmarshal([]byte(val.(string)), &data)
	json.Unmarshal([]byte(val), &data)

	return &orders.GetOrdersByUserIdResponse{Status: true, Data: data}, nil
}
