package pkgs

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _db *gorm.DB
var _redis *redis.Client

func InitDB() {
	dsn := fmt.Sprintf("root:123456@tcp(%s:3306)/restaurant?charset=utf8mb4&parseTime=True&loc=Local", "localhost")
	// var dsn = fmt.Sprintf("restaurant:Qwert_54321@tcp(%s:3306)/restaurant?charset=utf8mb4&parseTime=True&loc=Local", "localhost")
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	sqlDB, _ := _db.DB()
	sqlDB.SetMaxOpenConns(50) // 设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(10)

	_redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetDB() *gorm.DB {
	return _db
}

func GetRedis() *redis.Client {
	return _redis
}
