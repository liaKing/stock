package util

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func GetDataViper() (err error) {
	viper.SetConfigFile("./conf/data.json") // 指定配置文件路径
	err = viper.ReadInConfig()              // 读取配置信息
	if err != nil {                         // 读取配置信息失败
		panic(fmt.Errorf("GetDataViper	Fatal error config file: %s \n", err))
		return
	}
	//监视配置文件是否改变
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("DataViper配置文件修改了...")
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("GetDataViper	viper.Unmarshal failed, err:%v\n", err)
			return
		}
	})
	return
}
