package mysql

import (
	"fmt"
	"github.com/JobNing/framework/config"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type mysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func InitMysql(serviceName string) error {
	type Val struct {
		Mysql mysqlConfig `yaml:"mysql"`
	}
	mysqlConfigVal := Val{}
	content, err := config.GetConfig("DEFAULT_GROUP", serviceName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(content), &mysqlConfigVal)
	if err != nil {
		return err
	}
	configM := mysqlConfigVal.Mysql
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		configM.Username,
		configM.Password,
		configM.Host,
		configM.Port,
		configM.Database,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func WithTX(txFc func(tx *gorm.DB) error) {
	var err error
	tx := DB.Begin()
	err = txFc(tx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}
