package pkg

import (
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
