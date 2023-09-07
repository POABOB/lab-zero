// Code generated by goctl. DO NOT EDIT.
// Source: users.proto

package server

import (
	"context"

	"lab-zero/users/internal/logic"
	"lab-zero/users/internal/svc"
	"lab-zero/users/users"
)

type UsersServer struct {
	svcCtx *svc.ServiceContext
	users.UnimplementedUsersServer
}

func NewUsersServer(svcCtx *svc.ServiceContext) *UsersServer {
	return &UsersServer{
		svcCtx: svcCtx,
	}
}

func (s *UsersServer) GetUserInfoById(ctx context.Context, in *users.GetUserInfoRequest) (*users.GetUserInfoResponse, error) {
	l := logic.NewGetUserInfoByIdLogic(ctx, s.svcCtx)
	return l.GetUserInfoById(in)
}