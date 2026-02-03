package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Enabled  bool   `mapstructure:"REDIS_ENABLED"`
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

// load database config and set default value and marshalling value from file to struct
func (m *RedisConfig) load(v *viper.Viper) {
	v.SetDefault("REDIS_ENABLED", true)
	v.SetDefault("REDIS_HOST", "127.0.0.1")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	if err := v.Unmarshal(m); err != nil {
		log.Fatalf("unable to unmarshal RedisConfig: %v", err)
	}
}

// Addr returns the address of the Redis server in "host:port" format.
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
