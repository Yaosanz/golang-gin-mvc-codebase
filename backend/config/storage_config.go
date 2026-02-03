package config

import (
	"log"

	"github.com/spf13/viper"
)

// StorageConfig holds configuration for storage service
// todo: this is just example, you can remove this if you don't use storage config
type StorageConfig struct {
	Driver string `mapstructure:"STORAGE_DRIVER"`
	Bucket string `mapstructure:"STORAGE_BUCKET"`
}

// load smtp config and set default value and marshalling value from file to struct
func (m *StorageConfig) load(v *viper.Viper) {
	v.SetDefault("STORAGE_DRIVER", "")
	v.SetDefault("STORAGE_BUCKET", "25")

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal StorageConfig: %v", err.Error())
	}
}
