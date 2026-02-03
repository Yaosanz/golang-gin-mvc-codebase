package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// ConfigurationError Custom error types for better error handling
type ConfigurationError struct {
	Component string
	Err       error
}

func (e *ConfigurationError) Error() string {
	return fmt.Sprintf("configuration error in %s: %v", e.Component, e.Err)
}

type Config struct {
	app            AppConfig
	server         ServerConfig
	db             DatabaseConfig
	smtp           SmtpConfig
	minio          MinioConfig
	oidc           OidcConfig
	swagger        SwaggerConfig
	kafka          KafkaConfig
	redis          RedisConfig
	oca            OcaConfig
	jwt            JwtConfig
	MigrationPath  string
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")  // set the type of the config file (JSON, TOML, YAML, HCL, INI, ENV)
	v.AddConfigPath(".")    // tell Viper the location of the config file
	v.SetConfigFile(".env") // the actual config file

	// discover and read configuration file (this will override default value)
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		var pathError *os.PathError
		switch {
		case errors.As(err, &configFileNotFoundError):
			log.Printf("No config file found, using environment variables")
		case errors.As(err, &pathError):
			log.Printf("config file permission denied: %s", err)
		default:
			log.Printf("error reading config file: %s", err)
		}
		log.Printf("Falling back to server environment variables")
	}

	// Fallback and also read from server environment variables (this will override both value from default & file)
	v.AutomaticEnv()

	// load and unmarshal to map value to config struct
	var cfg Config
	cfg.app.load(v)
	cfg.server.load(v)
	cfg.swagger.load(v)
	cfg.db.load(v)
	cfg.smtp.load(v)
	cfg.minio.load(v)
	cfg.oidc.load(v)
	cfg.kafka.load(v)
	cfg.redis.load(v)
	cfg.oca.load(v)
	cfg.jwt.load(v)
	
	// Set migration path
	if migrationPath := os.Getenv("MIGRATION_PATH"); migrationPath != "" {
		cfg.MigrationPath = migrationPath
	} else {
		cfg.MigrationPath = "database/migrations"
	}

	// Validate all configurations
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate this is config validation runner
func (c *Config) validate() error {
	// register required validation method here

	// validate app config
	if err := c.app.validate(); err != nil {
		return err
	}
	return nil
}

// Jwt Getter method
func (c *Config) Jwt() JwtConfig {
	return c.jwt
}

// App Getter method
func (c *Config) App() AppConfig {
	return c.app
}

// Server Getter method
func (c *Config) Server() ServerConfig {
	return c.server
}

// Database Getter method
func (c *Config) Database() DatabaseConfig {
	return c.db
}

// Smtp Getter method
func (c *Config) Smtp() SmtpConfig {
	return c.smtp
}

// Minio Getter method
func (c *Config) Minio() MinioConfig {
	return c.minio
}

func (c *Config) Oidc() OidcConfig {
	return c.oidc
}

func (c *Config) Swagger() SwaggerConfig {
	return c.swagger
}

func (c *Config) Kafka() KafkaConfig {
	return c.kafka
}

func (c *Config) Redis() RedisConfig {
	return c.redis
}

func (c *Config) Oca() OcaConfig {
	return c.oca
}
