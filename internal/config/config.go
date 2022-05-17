package config

import "github.com/spf13/viper"

type Config struct {
	Port int
	dbConfig
}

type dbConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int
}

func New(v *viper.Viper) *Config {
	return &Config{
		Port: v.GetInt("port"),
		dbConfig: dbConfig{
			DBHost:     v.GetString("db_host"),
			DBUser:     v.GetString("db_user"),
			DBPassword: v.GetString("db_password"),
			DBName:     v.GetString("db_name"),
			DBPort:     v.GetInt("db_port"),
		},
	}
}
