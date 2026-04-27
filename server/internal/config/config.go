package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	DB       DBConfig       `mapstructure:"db"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	SMTP     SMTPConfig     `mapstructure:"smtp"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Log      LogConfig      `mapstructure:"log"`
	Attach   AttachConfig   `mapstructure:"attach"`
}

type AppConfig struct {
	Env         string `mapstructure:"env"`
	Port        string `mapstructure:"port"`
	FrontendURL string `mapstructure:"frontend_url"`
}

type DBConfig struct {
	Path         string        `mapstructure:"path"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type AttachConfig struct {
	Path string `mapstructure:"path"`
}

var C Config

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}

	return nil
}
