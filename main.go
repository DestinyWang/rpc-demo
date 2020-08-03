package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/DestinyWang/gokit-test/services"
	"github.com/DestinyWang/gokit-test/util"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
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
