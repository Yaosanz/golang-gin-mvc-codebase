package config

import (
	"github.com/spf13/viper"
	"log"
)

type SmtpConfig struct {
	Host     string `mapstructure:"SMTP_HOST"`
	Port     string `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
	From     string `mapstructure:"SMTP_FROM"`
}

// load smtp config and set default value and marshalling value from file to struct
func (m *SmtpConfig) load(v *viper.Viper) {
	v.SetDefault("SMTP_HOST", "")
	v.SetDefault("SMTP_PORT", "25")
	v.SetDefault("SMTP_USERNAME", "")
	v.SetDefault("SMTP_PASSWORD", "")
	v.SetDefault("SMTP_FROM", "noreply@example.com")

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal SMTPConfig: %v", err.Error())
	}
}
