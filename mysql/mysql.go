package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

// 连接mysql
func InitMysql() {
	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/zg6?charset=utf8mb4&parseTime=True&loc=Local"
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
