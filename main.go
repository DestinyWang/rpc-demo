package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/DestinyWang/gokit-test/services"
	"github.com/DestinyWang/gokit-test/util"
	"github.com/afex/hystrix-go/hystrix"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main1() {
	// flag
	var name = flag.String("name", "", "service name")
	var port = flag.Int("p", 0, "service port")
	flag.Parse()
	//
	logrus.WithFields(logrus.Fields{
		"name": *name,
		"port": *port,
	}).Info("server start")
	util.ServicePort = *port
	util.ServiceName = *name
	var id = fmt.Sprintf("%s:%s", util.ServiceName, uuid.New().String())
	// 初始化服务
	var userService = &services.UserService{}
	var limit = rate.NewLimiter(1, 3)
	var endPoint = services.RateLimit(limit)(services.GenUserEndpoint(userService))
	// 创建 handler
	var options = []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
			var contentType, body = "text/plain; charset=utf-8", []byte(err.Error())
			w.Header().Set("Content-type", contentType)
			if commonErr, ok := err.(*util.CommonErr); ok {
				w.WriteHeader(commonErr.Code)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(body)
		}),
	}
	
	var serverHandler = httpTransport.NewServer(endPoint, services.DecodeUserReq, services.EncodeUserResp, options...)
	// 路由
	var router = mux.NewRouter()
	router.Methods(http.MethodGet, http.MethodDelete, http.MethodPost).Path("/user/{uid:\\d+}").Handler(serverHandler)
	router.Methods(http.MethodGet).Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-type", "application/json")
		writer.Write([]byte(`{"status": "ok"}`))
	})
	var errCh = make(chan error)
	go func() {
		// 注册服务
		if err := util.RegisterService(id); err != nil {
			errCh <- err
		}
		if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
			logrus.WithError(err).Errorf("http server start fail")
			errCh <- err
		}
	}()
	go func() {
		// 信号监听
		var signCh = make(chan os.Signal)
		signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM) // 分别拦截 ctrl+c, kill 信号
		errCh <- fmt.Errorf("%s", <-signCh)
	}()
	var getErr = <-errCh
	logrus.WithError(getErr).Errorf("检测到服务异常, 开始注销服务")
	if err := util.DeregisterService(id); err != nil {
		logrus.WithError(err).Errorf("deregister service fail")
		panic(err)
	}
	time.Sleep(time.Second)
}

type Product struct {
	Id    int
	Title string
	Price int
}

func GetProduct() *Product {
	var r = rand.Intn(10)
	if r < 6 {
		time.Sleep(3 * time.Second) // 随机延迟三秒
	}
	return &Product{
		Id:    101,
		Title: "Product Title",
		Price: 12,
	}
}

func main() {
	var err error
	var configGetProduct = hystrix.CommandConfig{
		Timeout:                2000, // 超时时间
		MaxConcurrentRequests:  10,   // 最大并发数
		RequestVolumeThreshold: 20,   // 请求阈值, 默认 20, 有 20 个请求才进行错误百分比计算
		ErrorPercentThreshold:  50,    // 错误百分比
	}
	hystrix.ConfigureCommand("GetProduct", configGetProduct) // 关联
	for {
		time.Sleep(time.Second)
		if err = hystrix.Go("GetProduct", func() error {
			var p = GetProduct()
			logrus.Info(time.Now().Format("2006-01-02 15:04:05"), p)
			return nil
		}, func(err error) error {
			return nil
		}); err != nil {
			logrus.WithError(err).Error(time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}
