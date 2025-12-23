package config

import (
	"github.com/spf13/viper"
	"log"
)

type DatabaseConfig struct {
	Host                  string `mapstructure:"DB_HOST"`
	Port                  string `mapstructure:"DB_PORT"`
	Username              string `mapstructure:"DB_USERNAME"`
	Password              string `mapstructure:"DB_PASSWORD"`
	Name                  string `mapstructure:"DB_NAME"`
	Debug                 bool   `mapstructure:"DB_DEBUG"`
	Timezone              string `mapstructure:"DB_TIMEZONE"`
	SSLMode               string `mapstructure:"DB_SSL_MODE"`
	MaxOpenConnection     string `mapstructure:"DB_MAX_OPEN_CONNECTION"`
	MaxConnectionLifetime string `mapstructure:"DB_MAX_CONNECTION_LIFETIME"`
	MaxIdleLifetime       string `mapstructure:"DB_MAX_IDLE_LIFETIME"`
}

// load database config and set default value and marshalling value from file to struct
func (m *DatabaseConfig) load(v *viper.Viper) {
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "3306")
	v.SetDefault("DB_USERNAME", "root")
	v.SetDefault("DB_PASSWORD", "root")
	v.SetDefault("DB_NAME", "postgres")
	v.SetDefault("DB_DEBUG", false)
	v.SetDefault("DB_TIMEZONE", "UTC")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNECTION", "10")
	v.SetDefault("DB_MAX_CONNECTION_LIFETIME", "5m")
	v.SetDefault("DB_MAX_IDLE_LIFETIME", "5m")

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal DatabaseConfig: %v", err.Error())
	}
}
