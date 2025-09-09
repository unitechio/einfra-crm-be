
package config

import (
	"github.com/spf13/viper"
)

// Config holds all the configuration for the application.
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	DB      DBConfig      `mapstructure:"db"`
	Auth    AuthConfig    `mapstructure:"auth"`
	SMTP    SmtpConfig    `mapstructure:"smtp"`
	Storage StorageConfig `mapstructure:"storage"`
}

// ServerConfig holds the server configuration.
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DBConfig holds the database configuration.
type DBConfig struct {
	PostgresURL string `mapstructure:"postgres_url"`
}

// AuthConfig holds the authentication configuration.
type AuthConfig struct {
	Google OAuthConfig `mapstructure:"google"`
	Azure  OAuthConfig `mapstructure:"azure"`
}

// OAuthConfig holds the OAuth configuration.
type OAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	TenantID     string `mapstructure:"tenant_id,omitempty"` // For Azure AD
}

// SmtpConfig holds the SMTP configuration.
type SmtpConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// StorageConfig holds the file storage configuration.
type StorageConfig struct {
	ImagePath string `mapstructure:"image_path"`
}

// LoadConfig loads the configuration from a file.
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config") // Path to look for the config file in
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
