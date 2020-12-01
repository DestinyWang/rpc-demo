package main

import (
	"github.com/DestinyWang/gokit-test/grpc_test/services"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// 添加证书验证
	//transportCredentials, err := credentials.NewServerTLSFromFile("grpc_test/keys/server.crt", "grpc_test/keys/server_no_passwd.key")
	//if err != nil {
	//	logrus.WithError(err).Error("load server tls fail")
	//	return
	//}
	rpcServer := grpc.NewServer()
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		logrus.Error("listen fail")
		return
	}
	logrus.Info("listen")
	if err = rpcServer.Serve(listen); err != nil {
		logrus.Error("grpc serve fail")
	}
	logrus.Info("serve")
}
