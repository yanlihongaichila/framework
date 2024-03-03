package grpc

import (
	"fmt"
	"github.com/yanlihongaichila/framework/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

const (
	IPADDR      = "127.0.0.1"
	PORT        = 8848
	NAMESPACEID = "c787d4e6-a673-4b9e-baa5-2437bae2b891"
	GROUP       = "DEFAULT_GROUP"
)

type GatewayStru struct {
	APP struct {
		IPADDR   string `json:"IPADDR"`
		PORT     string `json:"PORT"`
		DATABASE string `json:"DATABASE"`
	} `json:"APP"`
	Mysql struct {
		User     string `json:"user"`
		Pwd      string `json:"pwd"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Datebase string `json:"datebase"`
	} `json:"Mysql"`
}

func Client(toService string) (*grpc.ClientConn, error) {
	cnfStr, err := nacos.InitNacos(&nacos.Nacos{
		IpAddr:      IPADDR,
		Port:        PORT,
		NamespaceId: NAMESPACEID,
		DataId:      toService,
		Group:       GROUP,
	})
	if err != nil {
		return nil, err
	}
	cnf := new(GatewayStru)
	err = yaml.Unmarshal([]byte(cnfStr), &cnf)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("%v:%v", cnf.APP.IPADDR, cnf.APP.PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
