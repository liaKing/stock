package config

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisConn *redis.Pool

// Pgsql 连接 PostgreSQL（读配置 pgsql.*，端口 5432）
func Pgsql() error {
	host := viper.GetString("pgsql.host")
	port := viper.GetInt("pgsql.port")
	user := viper.GetString("pgsql.user")
	password := viper.GetString("pgsql.password")
	dbname := viper.GetString("pgsql.dbname")
	sslmode := viper.GetString("pgsql.sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("pgsql connect err:", err)
		return err
	}
	DB = conn
	return nil
}

// Redis 连接redis
func Redis() error {
	ip := viper.GetString("redis.ip")
	maxIdle := viper.GetInt("redis.maxIdle")
	maxActive := viper.GetInt("redis.maxActive")

	redisCoon := &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ip)
		},
	}
	RedisConn = redisCoon
	return nil
}

func Close() {
	if RedisConn != nil {
		get := RedisConn.Get()
		_ = get.Close()
	}
	if DB != nil {
		sqlDB, _ := DB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}
}
