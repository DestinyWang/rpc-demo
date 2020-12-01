package main

import (
	"context"
	"github.com/DestinyWang/gokit-test/grpc_test/services"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Error("dial fail")
	}
	defer conn.Close()
	
	prodServiceClient := services.NewProdServiceClient(conn)
	prodResponse, err := prodServiceClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 10})
	if err != nil {
		logrus.WithError(err).Error("rpc fail")
	}
	logrus.Infof("prodResp=[%+v]", prodResponse)
}