package main

import (
	"github.com/DestinyWang/gokit-test/grpc_test/services"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
)

func main() {
	// 添加证书验证
	transportCredentials, err := credentials.NewServerTLSFromFile("grpc_test/keys/server.crt", "grpc_test/keys/server_no_passwd.key")
	if err != nil {
		logrus.WithError(err).Error("load server tls fail")
		return
	}
	rpcServer := grpc.NewServer(grpc.Creds(transportCredentials))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))
	//listen, err := net.Listen("tcp", ":8081")
	//if err != nil {
	//	logrus.Error("listen fail")
	//	return
	//}
	//logrus.Info("listen")
	//if err = rpcServer.Serve(listen); err != nil {
	//	logrus.Error("grpc serve fail")
	//}
	logrus.Info("serve")
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": request.Method,
			"header": request.Header,
			"proto": request.Proto,
		}).Info("server")
	})
	httpServer := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	if err = httpServer.ListenAndServeTLS("grpc_test/keys/server.crt", "grpc_test/keys/server_no_passwd.key"); err != nil {
		logrus.WithError(err).Error("listen and serve tls fail")
	}
}
