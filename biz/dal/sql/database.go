package config

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlConn *gorm.DB
var RedisConn *redis.Pool

// Mysql 连接mysql
func Mysql() error {

	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	ip := viper.GetString("mysql.ip")
	name := viper.GetString("mysql.name")

	//parseTime=True&loc=Local MySQL 默认时间是格林尼治时间，与我们差八小时，需要定位到我们当地时间
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, ip, name)
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	MysqlConn = conn
	return err
}

// Redis 连接redis
func Redis() error {
	ip := viper.GetString("redis.ip.json")
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
	get := RedisConn.Get()
	db, _ := MysqlConn.DB()
	_ = get.Close()
	_ = db.Close()
}
