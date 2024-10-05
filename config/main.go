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
	viper.SetDefault("USER_AGENT", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_DB", "0")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("CACHE_TTL", "86400")

	viper.AutomaticEnv()

	config.Port = viper.GetString("PORT")
	config.UserAgent = viper.GetString("USER_AGENT")
	config.RedisHost = viper.GetString("REDIS_HOST")
	config.RedisPort = viper.GetString("REDIS_PORT")
	config.RedisDB = viper.GetInt("REDIS_DB")
	config.CacheTTL = viper.GetUint("CACHE_TTL")

	return config
}
