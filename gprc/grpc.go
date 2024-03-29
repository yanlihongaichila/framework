package gprc

import (
	"fmt"
	"github.com/yanlihongaichila/framework/consul"
	"github.com/yanlihongaichila/framework/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

func getConfig(serviceName string) (*Config, error) {
	configInfo, err := nacos.GetConfig("DEFAULT_GROUP", serviceName)
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

func ConcentGrpc(serviceName string, fu func(s *grpc.Server)) error {
	cof, err := getConfig(serviceName)
	if err != nil {
		return err
	}
	ip, err := consul.GetIp()
	if err != nil {
		return err
	}
	//本地的话修改端口(端口随机)
	//docker的话改ip(ip随机)
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", ip, cof.App.Port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err.Error())
	}
	_, err = consul.InitRegisterServer("DEFAULT_GROUP", serviceName)
	if err != nil {

		return err
	}

	//健康检测
	s := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	//反射
	reflection.Register(s)

	fu(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}

func ConcentGrpcCert(port int, fu func(s *grpc.Server), cert, key string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)

	}
	s := grpc.NewServer(grpc.Creds(creds))
	//反射
	reflection.Register(s)
	fu(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return err
}
