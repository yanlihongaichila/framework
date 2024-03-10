package app

import (
	"github.com/yanlihongaichila/framework/mysql"
	"github.com/yanlihongaichila/framework/nacos"
	"strconv"
)

// 初始化服务
func Init(fileName, filePath, nacosName string, apps ...string) error {

	viper, err := nacos.InitViper(fileName, filePath, nacosName)

	if err != nil {
		return err
	}
	atoi, err := strconv.Atoi(viper["port"])
	if err != nil {
		return err
	}
	err = nacos.GetClient(viper["address"], uint64(atoi))

	for _, app := range apps {
		switch app {
		case "mysql":
			err = mysql.InitMysql(viper["group"], viper["name"])
			if err != nil {
				return err
			}

		}
	}

	return err
}
