package config

import (
	"log"

	"github.com/spf13/viper"
)

type SwaggerConfig struct {
	Title       string `mapstructure:"SWAGGER_TITLE"`
	Description string `mapstructure:"SWAGGER_DESCRIPTION"`
	Version     string `mapstructure:"SWAGGER_VERSION"`
	Host        string `mapstructure:"SWAGGER_HOST"`
	Schema      string `mapstructure:"SWAGGER_SCHEMA"`
	BasePath    string `mapstructure:"SWAGGER_BASE_PATH"`
}

// load smtp config and set default value and marshalling value from file to struct
func (m *SwaggerConfig) load(v *viper.Viper) {
	v.SetDefault("SWAGGER_TITLE", "Go Starter")
	v.SetDefault("SWAGGER_DESCRIPTION", "Go Starter API Documentation")
	v.SetDefault("SWAGGER_VERSION", "1.0")
	v.SetDefault("SWAGGER_HOST", "")
	v.SetDefault("SWAGGER_SCHEMA", "http,https")
	v.SetDefault("SWAGGER_BASE_PATH", "")

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal Swagger Config: %v", err.Error())
	}
}
