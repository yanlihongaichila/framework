package gprc

import (
	"github.com/yanlihongaichila/framework/nacos"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	_ "github.com/mbobakov/grpc-consul-resolver"
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

	//return grpc.Dial("10.2.171.14:8077", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.Dial("consul://10.2.171.14:8500/"+cnfs.App.Secret+"?wait=14s", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`))

}
