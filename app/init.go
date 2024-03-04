package app

import (
	"github.com/yanlihongaichila/framework/mysql"
	"github.com/yanlihongaichila/framework/nacos"
)

// 初始化服务
func Init(serviceName string, apps ...string) error {
	var err error
	err = nacos.GetClient()

	for _, app := range apps {
		switch app {
		case "mysql":
			mysql.InitMysql(serviceName)
		}
	}

	return err
}
