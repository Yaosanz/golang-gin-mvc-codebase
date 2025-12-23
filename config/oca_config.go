package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type OcaConfig struct {
	AppID             string        `mapstructure:"OCA_APP_ID"`
	MaxRetries        int           `mapstructure:"OCA_MAX_RETRIES"`
	RetryBackoff      time.Duration `mapstructure:"OCA_RETRY_BACKOFF"`
	EmailBaseURL      string        `mapstructure:"OCA_EMAIL_BASE_URL"`
	EmailAuthToken    string        `mapstructure:"OCA_EMAIL_AUTH_TOKEN"`
	EmailSenderEmail  string        `mapstructure:"OCA_EMAIL_SENDER_EMAIL"`
	EmailSenderName   string        `mapstructure:"OCA_EMAIL_SENDER_NAME"`
	WhatsAppBaseURL   string        `mapstructure:"OCA_WHATSAPP_BASE_URL"`
	WhatsAppAuthToken string        `mapstructure:"OCA_WHATSAPP_AUTH_TOKEN"`
}

// load app config and set default value and marshalling value from file to struct
func (m *OcaConfig) load(v *viper.Viper) {
	defaultRetryBackOff := 2 * time.Second
	v.SetDefault("OCA_MAX_RETRIES", 3)
	v.SetDefault("OCA_RETRY_BACKOFF", defaultRetryBackOff)

	// set default value for OCA email
	v.SetDefault("OCA_EMAIL_BASE_URL", "https://webapigw.ocatelkom.co.id")
	v.SetDefault("OCA_EMAIL_SENDER_EMAIL", "noreply@daemon.com")
	v.SetDefault("OCA_EMAIL_SENDER_NAME", "Postmaster")
	v.SetDefault("OCA_EMAIL_AUTH_TOKEN", "")

	// set default value for OCA whatsapp
	v.SetDefault("OCA_WHATSAPP_BASE_URL", "https://wa01.ocatelkom.co.id")
	v.SetDefault("OCA_WHATSAPP_AUTH_TOKEN", "")

	err := v.Unmarshal(&m)
	if err != nil {
		log.Fatalf("unable to unmarshal OcaConfig: %v", err.Error())
	}
}
