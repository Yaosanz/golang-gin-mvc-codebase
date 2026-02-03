package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type OidcConfig struct {
	Issuer            string `mapstructure:"OIDC_ISSUER"`
	ClientID          string `mapstructure:"OIDC_CLIENT_ID"`
	SkipClientIDCheck bool   `mapstructure:"OIDC_SKIP_CLIENT_ID_CHECK"`
	Enabled           bool   `mapstructure:"OIDC_ENABLED"`
}

// load server config and set default value and marshalling value from file to struct
func (m *OidcConfig) load(v *viper.Viper) {
	v.SetDefault("OIDC_ISSUER", "")
	v.SetDefault("OIDC_CLIENT_ID", "")
	v.SetDefault("OIDC_SKIP_CLIENT_ID_CHECK", false)
	v.SetDefault("OIDC_ENABLED", true)

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal OIDC config: %v", err.Error())
	}
}

// validate validate config value
func (v *OidcConfig) validate() error {
	// write validation logic here
	if v.Issuer == "" {
		return &ConfigurationError{
			Component: "OIDC_ISSUER",
			Err:       fmt.Errorf("cannot be empty"),
		}
	}

	if v.ClientID == "" {
		return &ConfigurationError{
			Component: "OIDC_CLIENT_ID",
			Err:       fmt.Errorf("cannot be empty"),
		}
	}

	return nil
}
