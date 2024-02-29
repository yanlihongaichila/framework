package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/yanlihongaichila/framework/nacos"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var Db *gorm.DB

// 连接mysql
type MysqlConfig struct {
	Mysql struct {
		User     string `json:"user"`
		Pwd      string `json:"pwd"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Datebase string `json:"datebase"`
	} `json:"Mysql"`
}

// 获取mysql配置
var MysqlInfo MysqlConfig

func GetMysqlConfig() error {
	initNacos, err := nacos.InitNacos(&nacos.Nacos{
		IpAddr:      "127.0.0.1",
		Port:        8848,
		NamespaceId: "c787d4e6-a673-4b9e-baa5-2437bae2b891",
		DataId:      "user_mysql",
		Group:       "DEFAULT_GROUP",
	})
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal([]byte(initNacos), &MysqlInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func InitMysql() {
	err := GetMysqlConfig()
	if err != nil {
		log.Println("failed to get mysql config")
		return
	}
	mConfig := MysqlInfo.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mConfig.User, mConfig.Pwd, mConfig.Host, mConfig.Port, mConfig.Datebase)
	//dsn := "root:root@tcp(127.0.0.1:3306)/zg6?charset=utf8mb4&parseTime=True&loc=Local"

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to open Mysql database")
	}
}

// 实时监控Mysql配置

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
