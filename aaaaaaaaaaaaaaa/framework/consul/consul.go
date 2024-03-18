package consul

import (
	"context"
	"fmt"
	"github.com/JobNing/framework/config"
	"github.com/JobNing/framework/redis"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
	"net"
	"strconv"
	"time"
)

const CONSUL_KEY = "consul:node:index"

type ConsulConfig struct {
	Consul struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"consul"`
}

func getConfig(nacosGroup, serviceName string) (*ConsulConfig, error) {
	cnf, err := config.GetConfig(nacosGroup, serviceName)
	if err != nil {
		return nil, err
	}

	consulCnf := new(ConsulConfig)
	err = yaml.Unmarshal([]byte(cnf), consulCnf)
	if err != nil {
		return nil, err
	}

	return consulCnf, err
}

func getIndex(ctx context.Context, serviceName string, indexLen int) (int, error) {
	exist, err := redis.ExistKey(ctx, serviceName, CONSUL_KEY)
	if err != nil {
		return 0, err
	}

	if exist {
		indexStr, err := redis.GetByKey(ctx, serviceName, CONSUL_KEY)
		if err != nil {
			return 0, err
		}
		index, err := strconv.Atoi(indexStr)
		newIndex := index + 1

		if newIndex >= indexLen {
			newIndex = 0
		}
		err = redis.SetKey(ctx, serviceName, CONSUL_KEY, newIndex, time.Duration(0))
		if err != nil {
			return 0, err
		}

		return index, nil
	}

	err = redis.SetKey(ctx, serviceName, "consul:node:index", 0, time.Duration(0))
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func AgentHealthService(ctx context.Context, nacosGroup, serviceName string) (string, error) {
	cof, err := getConfig(nacosGroup, serviceName)
	if err != nil {
		return "", err
	}

	client, err := capi.NewClient(&capi.Config{
		Address: fmt.Sprintf("%v:%v", cof.Consul.Ip, cof.Consul.Port),
	})
	if err != nil {
		return "", err
	}
	sr, infos, err := client.Agent().AgentHealthServiceByName(serviceName)
	if err != nil {
		return "", err
	}
	if sr != "passing" {
		return "", fmt.Errorf("is not have health service")
	}
	return fmt.Sprintf("%v:%v", infos[0].Service.Address, infos[0].Service.Port), nil
}

func getIps() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interfaces ipAddress: %v\n", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		// 检查ip地址判断是否回环地址
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func ServiceRegister(nacosGroup, serviceName string, address, port string) error {
	cof, err := getConfig(nacosGroup, serviceName)
	if err != nil {
		return err
	}
	client, err := capi.NewClient(&capi.Config{
		Address: fmt.Sprintf("%v:%v", cof.Consul.Ip, cof.Consul.Port),
	})
	if err != nil {
		return err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	err = client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    serviceName,
		Tags:    []string{"GRPC"},
		Port:    portInt,
		Address: address,
		Check: &capi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%v:%v", address, port),
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	return err
}
