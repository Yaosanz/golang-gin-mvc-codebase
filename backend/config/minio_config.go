package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type MinioConfig struct {
	Endpoint   string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey  string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey  string `mapstructure:"MINIO_SECRET_KEY"`
	BucketName string `mapstructure:"MINIO_BUCKET_NAME"`
	Secure     bool   `mapstructure:"MINIO_SECURE"`
}

// load app config and set default value and marshalling value from file to struct
func (m *MinioConfig) load(v *viper.Viper) {
	v.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	v.SetDefault("MINIO_ACCESS_KEY", "")
	v.SetDefault("MINIO_SECRET_KEY", "")
	v.SetDefault("MINIO_BUCKET_NAME", "")
	v.SetDefault("MINIO_SECURE", false)

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal minioConfig: %v", err.Error())
	}
}

// validate validate config value
func (v *MinioConfig) validate() error {
	// write validation logic here
	if v.Endpoint == "" {
		return &ConfigurationError{
			Component: "MINIO_ENDPOINT",
			Err:       fmt.Errorf("cannot be empty"),
		}
	}

	if v.AccessKey == "" {
		return &ConfigurationError{
			Component: "MINIO_ACCESS_KEY",
			Err:       fmt.Errorf("cannot be empty"),
		}
	}

	if v.SecretKey == "" {
		return &ConfigurationError{
			Component: "MINIO_SECRET_KEY",
			Err:       fmt.Errorf("cannot be empty"),
		}
	}

	if v.BucketName == "" {
		return &ConfigurationError{
			Component: "MINIO_BUCKET_NAME",
			Err:       fmt.Errorf("cannot be empty"),
		}
	}
	return nil
}
