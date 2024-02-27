package nacos

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var ConfigAll Config

type Config struct {
	Mysql struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	} `json:"Mysql"`
}

func Nacos() {

	//create clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}
	// Create config client for dynamic configuration
	config, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	getConfig, err := config.GetConfig(vo.ConfigParam{
		DataId: "user",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		return
	}
	json.Unmarshal([]byte(getConfig), &ConfigAll)

}
