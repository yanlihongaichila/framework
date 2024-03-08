package consul

import (
	"fmt"
	"github.com/yanlihongaichila/framework/nacos"
	"gopkg.in/yaml.v2"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var (
	ConsulClient *api.Client
	SrvId        string
	err          error
)

type ConsulConfigs struct {
	Consul struct {
		Ip      string `yaml:"ip"`
		Port    int    `yaml:"port"`
		Version string `yaml:"version"`
	} `yaml:"consul"`
}

type RpcConfigs struct {
	Rpc struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
		Key     string `yaml:"key"`
	} `yaml:"rpc"`
}

func getConsulConfig(group, service string) (string, error) {
	config, err := nacos.GetConfig(group, service)
	if err != nil {
		return "", err
	}
	return config, nil
}

/*
***************************写到连接rpc中
 */
//shh := grpc.NewServer() // 创建gRPC服务器
//healthcheck := health.NewServer()
//healthpb.RegisterHealthServer(shh, healthcheck)
func InitRegisterServer(group, service string) error {
	consulConfig, err := getConsulConfig(group, service)
	if err != nil {
		return err
	}
	consulCon := ConsulConfigs{}
	gprcCon := RpcConfigs{}
	err = yaml.Unmarshal([]byte(consulConfig), &consulCon)
	err = yaml.Unmarshal([]byte(consulConfig), &gprcCon)
	if err != nil {
		return err
	}

	cfig := consulCon.Consul
	rfig := gprcCon.Rpc
	//使用默认配置
	config := api.DefaultConfig()

	//配置consul的连接地址
	config.Address = fmt.Sprintf("%s:%d", cfig.Ip, cfig.Port)

	//示例化客户端
	ConsulClient, err = api.NewClient(config)

	if err != nil {
		fmt.Println(err)
		zap.S().Panic(err.Error())
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", rfig.Address, rfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "20s",
	}

	//健康检查,检查我们注册的微服务
	Registration := api.AgentServiceRegistration{}
	Registration.Address = rfig.Address
	Registration.Port = rfig.Port
	Registration.Name = rfig.Key
	Registration.Tags = []string{cfig.Version}
	Registration.ID = fmt.Sprintf("%s", uuid.NewV4())
	SrvId = Registration.ID
	Registration.Check = check

	err = ConsulClient.Agent().ServiceRegister(&Registration)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	return nil
}
