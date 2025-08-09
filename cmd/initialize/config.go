package initialize

import (
	"fmt"
	"os"

	"github.com/xinliangnote/go-gin-api/cmd/global"

	"github.com/spf13/viper"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig(env string) {
	//从配置文件中读取对应的配置
	pwd, _ := os.Getwd()
	configFileName := fmt.Sprintf("%s/etc/config-%s.yaml", pwd, env)

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}

func InitConfigTest(filepath string) {
	//从配置文件中读取对应的配置
	v := viper.New()
	v.SetConfigFile(filepath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
