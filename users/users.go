package main

import (
	"flag"
	"fmt"

	"lab-zero/users/internal/config"
	"lab-zero/users/internal/server"
	"lab-zero/users/internal/svc"
	"lab-zero/users/users"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/users.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	gw := gateway.MustNewServer(c.Gateway)
	group := service.NewServiceGroup()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		users.RegisterUsersServer(grpcServer, server.NewUsersServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	group.Add(s)

	group.Add(gw)
	defer group.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	group.Start()
}
