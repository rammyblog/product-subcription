package config

import (
	"github.com/rammyblog/go-paystack"
	"gorm.io/gorm"
)

// AppConfig represents the global configuration object
type AppConfig struct {
	DB             *gorm.DB
	PaystackClient *paystack.Client
}

// GlobalConfig is the global instance of AppConfig
var GlobalConfig *AppConfig
