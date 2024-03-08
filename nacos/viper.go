package nacos

import (
	"fmt"
	"github.com/spf13/viper"
)

// 用viper拿取nacos的配置
func InitViper(fileName, filePath, serviceName string) (map[string]string, error) {
	viper.SetConfigFile(fileName)
	viper.AddConfigPath(filePath)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			return nil, fmt.Errorf("no such config file")

		} else {

			return nil, fmt.Errorf("read config error")

		}
	}

	mapString := viper.GetStringMapString(serviceName)
	return mapString, nil
}
