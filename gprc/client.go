package gprc

import (
	"fmt"
	"github.com/yanlihongaichila/framework/consul"
	"github.com/yanlihongaichila/framework/nacos"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	_ "github.com/mbobakov/grpc-consul-resolver"
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
	//拿取consul配置
	con := consul.ConsulConfigs{}
	err = yaml.Unmarshal([]byte(cnfStr), &con)
	//return grpc.Dial("10.2.171.14:8077", grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.Dial(fmt.Sprintf("consul://%v:%v/", con.Consul.Ip, con.Consul.Port)+cnfs.App.Secret+"?wait=14s", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`))

}
