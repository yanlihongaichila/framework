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
	Mysql struct {
		User     string `json:"user"     yaml:"User"`
		Pwd      string `json:"pwd"      yaml:"Pwd"`
		Host     string `json:"host"     yaml:"Host"`
		Port     string `json:"port"     yaml:"Port"`
		Datebase string `json:"datebase" yaml:"Datebase"`
	} `json:"Mysql"`
}

// 获取mysql配置
var MysqlInfo MysqlConfig

func GetMysqlConfig(serviceName string) error {
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
	type Val struct {
		Mysql MysqlConfig `yaml:"mysql"`
	}
	mysqlConfigVal := Val{}
	content, err := nacos.GetConfig("DEFAULT_GROUP", serviceName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(content), &mysqlConfigVal)
	if err != nil {
		fmt.Println("**********errr")
		return err
	}
	fmt.Println(content)
	fmt.Println(mysqlConfigVal)

	return nil
}

func InitMysql(serviceName string) {
	err := GetMysqlConfig(serviceName)
	if err != nil {
		log.Println("failed to get mysql config")
		return
	}
	mConfig := MysqlInfo.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mConfig.User, mConfig.Pwd, mConfig.Host, mConfig.Port, mConfig.Datebase)
	//dsn := "root:root@tcp(127.0.0.1:3306)/zg6?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to open Mysql database")
	}
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
