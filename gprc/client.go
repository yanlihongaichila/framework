package gprc

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

type AppConfig struct {
	Ip     string `yaml:"ip"`
	Port   string `yaml:"port"`
	Secret string `yaml:"secret"`
}
type Val struct {
	App AppConfig `yaml:"app"`
}

func Client(toService string) (*grpc.ClientConn, error) {
	cnfStr, err := nacos.GetConfig("DEFAULT_GROUP", toService)
	if err != nil {
		return nil, err
	}
	cnfs := new(Val)
	err = yaml.Unmarshal([]byte(cnfStr), &cnfs)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("%v:%v", cnfs.App.Ip, cnfs.App.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
