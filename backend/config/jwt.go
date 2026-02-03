package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type JwtConfig struct {
	Secret     string
	ExpiredIn time.Duration
	Issuer     string
}

func (c *JwtConfig) load(v *viper.Viper) {
	expired, err := time.ParseDuration(v.GetString("JWT_EXPIRED_IN"))
	if err != nil {
		expired = 24 * time.Hour // default fallback
	}

	c.Secret = v.GetString("JWT_SECRET")
	c.Issuer = v.GetString("JWT_ISSUER")
	c.ExpiredIn = expired
}

func (c *JwtConfig) validate() error {
	if c.Secret == "" {
		return &ConfigurationError{
			Component: "JWT",
			Err:       errors.New("JWT_SECRET is required"),
		}
	}

	if c.Issuer == "" {
		return &ConfigurationError{
			Component: "JWT",
			Err:       errors.New("JWT_ISSUER is required"),
		}
	}

	return nil
}
