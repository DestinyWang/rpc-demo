package main

import (
	"context"
	"github.com/DestinyWang/gokit-test/grpc_test/services"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 证书
	transportCredentials, err := credentials.NewClientTLSFromFile("grpc_test/keys/server.crt", "destiny.com")
	if err != nil {
		logrus.WithError(err).Error("new client tls from file fail")
		return
	}
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		logrus.WithError(err).Error("dial fail")
		return
	}
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	
	prodClient := services.NewProdServiceClient(conn)
	prodResponse, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 10})
	if err != nil {
		logrus.WithError(err).WithField("resp", prodResponse).Error("rpc fail")
		return
	}
	logrus.Infof("prodResp=[%+v]", prodResponse)
}