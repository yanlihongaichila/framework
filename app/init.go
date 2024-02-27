package app

import "github.com/yanlihongaichila/framework/mysql"

// 初始化服务
func Init(apps ...string) error {
	var err error

	for _, app := range apps {
		switch app {
		case "mysql":
			mysql.InitMysql()
		}
	}

	return err
}
