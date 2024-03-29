package consul

import (
	"errors"
	"fmt"
	"github.com/yanlihongaichila/framework/nacos"
	"gopkg.in/yaml.v2"
	"log"
	"net"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var (
	ConsulClient *api.Client
)

type ConsulConfigs struct {
	Consul struct {
		Ip      string `yaml:"ip"`
		Port    int    `yaml:"port"`
		Secret  string `yaml:"secret"`
		Version string `yaml:"version"`
	} `yaml:"consul"`
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

// 获取地址
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

func GetIp() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, isVailIpNet := addr.(*net.IPNet)
			if isVailIpNet && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					// 添加一些额外的检测逻辑，例如判断IP地址是否在本地网络范围内
					if ipNet.IP.IsGlobalUnicast() {
						// 添加详细的日志输出
						log.Printf("获取到的IP地址：%s，对应网络接口：%s\n", ipNet.IP.String(), i.Name)
						return ipNet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", errors.New("Unable to find a valid global unicast IP address")
}

/*
***************************写到连接rpc中
 */
//shh := grpc.NewServer() // 创建gRPC服务器
//healthcheck := health.NewServer()
//healthpb.RegisterHealthServer(shh, healthcheck)
func InitRegisterServer(group, service string) (string, error) {
	ip, err := GetIp()
	consulConfig, err := getConsulConfig(group, service)

	if err != nil {
		return "", err
	}
	consulCon := ConsulConfigs{}
	err = yaml.Unmarshal([]byte(consulConfig), &consulCon)
	if err != nil {
		return "", err
	}
	//consul配置
	cfig := consulCon.Consul

	//rpc配置
	rfig := consulCon.Rpc

	//使用默认配置
	config := api.Config{Address: fmt.Sprintf("%v:%v", cfig.Ip, 8500)}

	//示例化客户端
	ConsulClient, err = api.NewClient(&config)

	if err != nil {
		fmt.Println(err)
		zap.S().Panic(err.Error())
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%v:%v", ip, rfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "20s",
	}

	//健康检查,检查我们注册的微服务
	err = ConsulClient.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s", uuid.NewV4()),
		Name:    rfig.Key,
		Tags:    []string{cfig.Version},
		Port:    rfig.Port,
		Address: ip,
		Check:   check,
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return ip, nil
}
