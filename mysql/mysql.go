package mysql

import (
	"fmt"
	"github.com/yanlihongaichila/framework/nacos"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var Db *gorm.DB

// 连接mysql
type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}
type Val struct {
	Mysql MysqlConfig `yaml:"mysql"`
}

var MysqlConfigVal Val

func GetMysqlConfig(serviceGroup, serviceName string) error {
	/*
		竣文
	*/
	//initNacos, err := nacos.InitNacos(&nacos.Nacos{
	//	IpAddr:      "127.0.0.1",
	//	Port:        8848,
	//	NamespaceId: "c787d4e6-a673-4b9e-baa5-2437bae2b891",
	//	DataId:      "user_mysql",
	//	Group:       "DEFAULT_GROUP",
	//})
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//err = json.Unmarshal([]byte(initNacos), &MysqlInfo)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//return nil

	content, err := nacos.GetConfig(serviceGroup, serviceName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(content), &MysqlConfigVal)
	if err != nil {
		fmt.Println("**********errr")
		return err
	}
	fmt.Println("22222222222222")
	fmt.Println(content)
	fmt.Println(MysqlConfigVal)

	return nil
}

func InitMysql(serviceGroup, serviceName string) error {
	err := GetMysqlConfig(serviceGroup, serviceName)
	if err != nil {
		log.Println("failed to get mysql config")
		return err
	}
	mConfig := MysqlConfigVal.Mysql
	fmt.Println(mConfig)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mConfig.Username, mConfig.Password, mConfig.Host, mConfig.Port, mConfig.Database)
	//dsn := "root:root@tcp(127.0.0.1:3306)/zg6?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	return nil
}

// 回滚
func WithTX(txFc func(tx *gorm.DB) error) {
	var err error
	tx := Db.Begin()
	err = txFc(tx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}
