package settings

import (
	"errors"

	"github.com/spf13/viper"
)

//新建一个全局对象
var Cnf = new(AppConf)

type AppConf struct {
	App App
	Mysql DBConfig
	Redis RedisConfig

}


type App struct {
	Name string `mapstructure:"name"`
	version string `mapstructure:"version"`
	Port string	`mapstructure:"port"`
	Log string `mapstructure:"logPath"`
	Level string `mapstructure:"logLevel"`
	StartTime string `mapstructure:"start_time"`
	MechineId int64 `mapstructure:"machine_id"`
}


type DBConfig struct  {
	Host string `mapstructure:"host"`
	DBName string `mapstructure:"dbname"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Charset string `mapstructure:"charset"`
}


type RedisConfig struct {
	Host string `mapstructure:"host"`
	Passwrod string `mapstructure:"password"`
	Port string `mapstructure:"port"`
}
//定义viper初始化
func Init() error {
	//定义viper读取路径\
	viper.AddConfigPath("settings/")
	viper.AddConfigPath("./")

	viper.SetConfigName("config")

	viper.SetConfigType("yaml")

	//读取文件到viper中
	err := viper.ReadInConfig()
	if err != nil {
		return  errors.New("viper读取配置失败")
	}

	//反序列化到Cnf结构体中
	if err := viper.Unmarshal(Cnf); err != nil {
		return errors.New("viper解析配置失败")
	}
	return  nil
}