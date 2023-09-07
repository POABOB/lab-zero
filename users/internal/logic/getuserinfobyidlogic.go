package logic

import (
	"context"

	"lab-zero/orders/orders"
	"lab-zero/users/internal/svc"
	"lab-zero/users/users"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByIdLogic {
	return &GetUserInfoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoByIdLogic) GetUserInfoById(in *users.GetUserInfoRequest) (*users.GetUserInfoResponse, error) {
	// validate
	if err := in.ValidateAll(); err != nil {
		return &users.GetUserInfoResponse{Message: err.Error()}, nil
	}

	// call Orders via gRPC
	// OrdersRpc := ordersclient.NewOrders(zrpc.MustNewClient(l.svcCtx.Config.OrdersServiceRpc))
	resp, err := l.svcCtx.OrdersRpc.GetOrdersByUserId(l.ctx, &orders.GetOrdersByUserIdRequest{
		UserId: in.UserId,
	})
	if err != nil {
		return nil, err
	}

	// copy struct
	userOrders := []*users.UserOrders{}
	copier.Copy(&userOrders, &resp.Data)

	return &users.GetUserInfoResponse{Status: true, Message: "", Data: userOrders}, nil
}
