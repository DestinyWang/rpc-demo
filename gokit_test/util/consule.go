package util

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

var consulClient *consulApi.Client

var (
	ServiceName string
	ServicePort int
)

func init() {

}

func RegisterService(id string) (err error) {
	var config = consulApi.DefaultConfig()
	config.Address = "192.168.0.108:8500"
	var reg = &consulApi.AgentServiceRegistration{
		ID:      id,
		Name:    ServiceName,
		Tags:    []string{"primary"},
		Port:    ServicePort,
		Address: "192.168.0.108",
		Check: &consulApi.AgentServiceCheck{
			Interval: "5s",
			HTTP:     fmt.Sprintf("http://192.168.0.108:%d/health", ServicePort),
		},
	}
	logrus.Infof("reg=[%+v]", reg)

	if consulClient, err = consulApi.NewClient(config); err != nil {
		logrus.WithError(err).Errorf("init client fail")
		return err
	}
	if err = consulClient.Agent().ServiceRegister(reg); err != nil {
		logrus.WithError(err).Errorf("register consul fail")
		return err
	}
	return nil
}

func DeregisterService(id string) (err error) {
	logrus.WithField("serviceId", id).Info("deregister service")
	return consulClient.Agent().ServiceDeregister(id)
}
