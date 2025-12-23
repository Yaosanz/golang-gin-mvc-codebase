package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name           string   `mapstructure:"APP_NAME"`
	Env            string   `mapstructure:"APP_ENV"`
	Timezone       string   `mapstructure:"APP_TIMEZONE"`
	Debug          bool     `mapstructure:"APP_DEBUG"`
	TrustedProxies []string `mapstructure:"APP_TRUSTED_PROXIES"`
	AllowedOrigins []string `mapstructure:"APP_ALLOWED_ORIGINS"`
	EnableCron     bool     `mapstructure:"APP_ENABLE_CRON"`
}

// load app config and set default value and marshalling value from file to struct
func (m *AppConfig) load(v *viper.Viper) {
	v.SetDefault("APP_NAME", "golang-app")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_TIMEZONE", "UTC")
	v.SetDefault("APP_DEBUG", true)
	v.SetDefault("APP_TRUSTED_PROXIES", []string{})
	v.SetDefault("APP_ALLOWED_ORIGINS", []string{})
	v.SetDefault("APP_ENABLE_CRON", false)

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal AppConfig: %v", err.Error())
	}
}

// validate validate config value
func (v *AppConfig) validate() error {
	// write validation logic here
	if !isValidEnvironment(v.Env) {
		return &ConfigurationError{
			Component: "APP_ENV",
			Err:       fmt.Errorf("invalid value %s", v.Env),
		}
	}
	return nil
}

// Helper function for environment validation
func isValidEnvironment(env string) bool {
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}
	return validEnvs[strings.ToLower(env)]
}
