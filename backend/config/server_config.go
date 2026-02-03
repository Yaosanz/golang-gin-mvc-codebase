package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host            string        `mapstructure:"SERVER_HOST"`
	Port            string        `mapstructure:"SERVER_PORT"`
	ReadTimeout     time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `mapstructure:"SERVER_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
	TLS             bool          `mapstructure:"SERVER_TLS"`
	TLSVersion      string        `mapstructure:"SERVER_TLS_VERSION"`
	CRTFile         string        `mapstructure:"SERVER_CRT_FILE"`
	KeyFile         string        `mapstructure:"SERVER_KEY_FILE"`
}

// load server config and set default value and marshalling value from file to struct
func (m *ServerConfig) load(v *viper.Viper) {
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("SERVER_READ_TIMEOUT", "15s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "15s")
	v.SetDefault("SERVER_IDLE_TIMEOUT", "60s")
	v.SetDefault("SERVER_SHUTDOWN_TIMEOUT", "5s")
	v.SetDefault("SERVER_TLS", false)
	v.SetDefault("SERVER_TLS_VERSION", "TLSv1.2")
	v.SetDefault("SERVER_CRT_FILE", "server.crt")
	v.SetDefault("SERVER_KEY_FILE", "server.key")

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal ServerConfig: %v", err.Error())
	}
}
