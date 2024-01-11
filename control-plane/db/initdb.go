package db

import (
	"log"
	"time"

	"control-plane/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Orm *gorm.DB

func initDB() {
}

func gormDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DbConfig.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	mysqlDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	mysqlDb.SetMaxOpenConns(50)
	mysqlDb.SetConnMaxIdleTime(10)
	mysqlDb.SetConnMaxLifetime(time.Hour)
	return db
}
func InitDb() {
	Orm = gormDB()
}
