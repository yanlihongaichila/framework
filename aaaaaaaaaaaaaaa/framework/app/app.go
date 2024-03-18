package app

import (
	"github.com/JobNing/framework/config"
	"github.com/JobNing/framework/mysql"
)

func Init(
	serviceName string,
	nacosIP, naocsPort string,
	apps ...string,
) error {
	if err := config.GetClient(nacosIP, naocsPort); err != nil {
		return err
	}

	for _, val := range apps {
		switch val {
		case "mysql":
			err := mysql.InitMysql(serviceName)
			if err != nil {
				panic(err)
			}
		}

	}
	return nil
}
