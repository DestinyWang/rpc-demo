package main

import (
	"context"
	services "github.com/DestinyWang/gokit_test/client/services"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httpTransport "github.com/go-kit/kit/transport/http"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main2() {
	var err error
	// step1: 创建 client
	var conf = consulApi.DefaultConfig()
	conf.Address = "192.168.0.108:8500" // 注册中心地址
	var apiClient *consulApi.Client
	if apiClient, err = consulApi.NewClient(conf); err != nil {
		logrus.WithError(err).Error("init consul apiClient fail")
		panic(err)
	}
	var client = consul.NewClient(apiClient)
	// step2: 创建 instance
	var logger = log.NewLogfmtLogger(os.Stdout)
	var instance = consul.NewInstancer(client, logger, "userservice1", []string{"primary"}, true)
	// step3: 创建 EndPointer
	var factoryFuc = func(serviceUrl string) (endpoint.Endpoint, io.Closer, error) {
		var target, _ = url.Parse("http://" + serviceUrl)
		return httpTransport.NewClient(http.MethodGet, target, services.EncodeRequestFunc, services.DecodeRequestFunc).Endpoint(), nil, nil
	}
	var endPointer = sd.NewEndpointer(instance, factoryFuc, logger)
	var myLb = lb.NewRoundRobin(endPointer)
	for {
		// step4: 调用 userService 服务
		var userResp interface{}
		var getUserInfo, _ = myLb.Endpoint()
		if userResp, err = getUserInfo(context.Background(), services.UserReq{
			Uid: 101,
		}); err != nil {
			logrus.WithError(err).Errorf("services fail")
			panic(err)
		}
		var userInfo = userResp.(*services.UserResp)
		logrus.Infof("userInfo=[%+v]", userInfo)
		time.Sleep(time.Second)
	}
}

func routeEndPoints(endPoints []endpoint.Endpoint) endpoint.Endpoint {
	random := rand.Intn(len(endPoints))
	endPoint := endPoints[random]
	logrus.WithFields(logrus.Fields{
		"length": len(endPoints),
		"random": random,
	}).Infof("route result")
	return endPoint
}
