package app

import (
    "github.com/spf13/viper"
)

type Config struct {}

func (*Config) GetString(name string) string {
    return viper.GetString(name)
}
