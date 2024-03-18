package grpc

import (
	"fmt"
	"github.com/JobNing/framework/config"
	"github.com/JobNing/framework/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"log"
	"net"
)

type Config struct {
	App struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"app"`
}

func getConfig(nacosGroup, serviceName string) (*Config, error) {
	configInfo, err := config.GetConfig(nacosGroup, serviceName)
	if err != nil {
		return nil, err
	}
	cnf := new(Config)
	err = yaml.Unmarshal([]byte(configInfo), cnf)
	if err != nil {
		return nil, err
	}
	return cnf, nil
}

func RegisterGRPC(nacosGroup, serviceName string, register func(s *grpc.Server)) error {
	cof, err := getConfig(nacosGroup, serviceName)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "0.0.0.0", cof.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	err = consul.ServiceRegister(nacosGroup, serviceName, cof.App.Ip, cof.App.Port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	//反射接口支持查询
	reflection.Register(s)
	//支持健康检查
	healthpb.RegisterHealthServer(s, health.NewServer())

	register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return err
}
