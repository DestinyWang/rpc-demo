package util

import (
	consulApi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

var consulClient *consulApi.Client

func RegisterService() (err error) {
	var config = consulApi.DefaultConfig()
	config.Address = "192.168.0.108:8500"
	var reg = &consulApi.AgentServiceRegistration{
		ID:      "userservice",
		Name:    "userservice",
		Tags:    []string{"primary"},
		Port:    8000,
		Address: "192.168.0.108",
		Check: &consulApi.AgentServiceCheck{
			Interval: "5s",
			HTTP:     "http://192.168.0.108:8000/health",
		},
	}

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

func DeregisterService() (err error) {
	return consulClient.Agent().ServiceDeregister("userservice")
}
