package app

import (
    "github.com/spf13/viper"
)

type Config struct {}

func (*Config) GetString(name string) string {
    return viper.GetString(name)
}

func (*Config) GetInt(name string) int {
    return viper.GetInt(name)
}

func NewConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    err := viper.ReadInConfig()
    if err != nil {
        return nil, err
    }

    return &Config{}, nil
}
