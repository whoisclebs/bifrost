package config

import (
	"github.com/spf13/viper"
	"time"
)

type ProxyConfig struct {
	App    *AppConfig    `mapstructure:"app"`
	Redis  *RedisConfig  `mapstructure:"redis"`
	Eureka *EurekaConfig `mapstructure:"eureka"`
	Sentry *SentryConfig `mapstructure:"sentry"`
}

type AppConfig struct {
	Name         string        `mapstructure:"name"`
	Port         int           `mapstructure:"port"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	LoggingLevel string        `mapstructure:"logging_level"`
	Language     string        `mapstructure:"lang"`
	Interceptors []string      `mapstructure:"interceptors"`
}

type RedisConfig struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
	LocalCache     struct {
		Size           int           `mapstructure:"size"`
		ExpirationTime time.Duration `mapstructure:"expiration_time"`
	} `mapstructure:"local_cache"`
}

type EurekaConfig struct {
	Addresses []string `mapstructure:"addresses"`
}

type SentryConfig struct {
	SensitiveHeaders string                       `mapstructure:"sensitive-headers"`
	IgnoredHeaders   string                       `mapstructure:"ignored-headers"`
	Host             SentryHostConfig             `mapstructure:"host"`
	Routes           map[string]SentryRouteConfig `mapstructure:"routes"`
}

type SentryHostConfig struct {
	SocketTimeoutMillis  int `mapstructure:"socket-timeout-millis"`
	ConnectTimeoutMillis int `mapstructure:"connect-timeout-millis"`
	ReadTimeoutMillis    int `mapstructure:"read-timeout-millis"`
}

type SentryRouteConfig struct {
	SensitiveHeaders string `mapstructure:"sensitive-headers"`
	Path             string `mapstructure:"path"`
	ServiceId        string `mapstructure:"serviceId"`
	Provider         string `mapstructure:"provider"`
}

var (
	c   ProxyConfig
	err error
)

func Init() {
	viper.SetDefault("hostname", "hostname-not-defined")
	viper.AddConfigPath("resources")
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		panic(err)
	}
}

func Get() ProxyConfig {
	return c
}
