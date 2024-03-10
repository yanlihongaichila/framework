package nacos

import (
	"fmt"
	"github.com/spf13/viper"
)

// 用viper拿取nacos的配置
func InitViper(fileName, filePath, serviceName string) (map[string]string, error) {

	viper.SetConfigName(fileName)
	viper.AddConfigPath(filePath)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	mapString := viper.GetStringMapString(serviceName)
	fmt.Println(mapString)
	return mapString, nil
}
