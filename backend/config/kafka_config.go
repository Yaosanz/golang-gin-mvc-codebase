package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers          []string `mapstructure:"KAFKA_BROKERS"`
	GroupID          string   `mapstructure:"KAFKA_GROUP_ID"`
	Topic            string   `mapstructure:"KAFKA_TOPIC"`
	Username         string   `mapstructure:"KAFKA_USERNAME"`
	Password         string   `mapstructure:"KAFKA_PASSWORD"`
	SasEnabled       bool     `mapstructure:"KAFKA_SAS_ENABLED"`
	SecurityProtocol string   `mapstructure:"KAFKA_SECURITY_PROTOCOL"`
}

// load server config and set default value and marshalling value from file to struct
func (m *KafkaConfig) load(v *viper.Viper) {
	v.SetDefault("KAFKA_BROKERS", []string{"localhost:9092"})
	v.SetDefault("KAFKA_GROUP_ID", "")
	v.SetDefault("KAFKA_TOPIC", false)
	v.SetDefault("KAFKA_SAS_ENABLED", false)
	v.SetDefault("KAFKA_SECURITY_PROTOCOL", "PLAINTEXT")

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal OIDC config: %v", err.Error())
	}
}

// validate validate config value
func (v *KafkaConfig) validate() error {
	if v.Brokers != nil && len(v.Brokers) == 0 {
		return &ConfigurationError{
			Component: "KAFKA_BROKERS",
			Err:       fmt.Errorf("cannot be empty"),
		}
	}

	return nil
}
