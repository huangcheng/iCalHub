package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	UserAgent string
	RedisHost string
	RedisPort string
	RedisDB   int
	CacheTTL  uint
}

func ReadConfig() Config {
	var config Config

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("USER_AGENT", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:130.0) Gecko/20100101 Firefox/130.0")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_DB", "0")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("CACHE_TTL", "3600")

	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")

	viper.ReadInConfig()

	config.Port = viper.GetString("PORT")
	config.UserAgent = viper.GetString("USER_AGENT")
	config.RedisHost = viper.GetString("REDIS_HOST")
	config.RedisPort = viper.GetString("REDIS_PORT")
	config.RedisDB = viper.GetInt("REDIS_DB")
	config.CacheTTL = viper.GetUint("CACHE_TTL")

	return config
}
