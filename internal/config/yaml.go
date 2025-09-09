package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
// The values are read by Viper from a config file and/or environment variables.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Security SecurityConfig `mapstructure:"security"`
}

// ServerConfig holds server-specific configuration.
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"readTimeout"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
	IdleTimeout  int    `mapstructure:"idleTimeout"`
}

// SecurityConfig holds security-related configuration.
type SecurityConfig struct {
	JWTSecret      string   `mapstructure:"jwtSecret"`
	CORSEnabled    bool     `mapstructure:"corsEnabled"`
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
	RateLimitRPS   float64  `mapstructure:"rateLimitRps"` // Requests per second
	RateLimitBurst int      `mapstructure:"rateLimitBurst"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfigFile() (cfg Config, err error) {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../") // For tests
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	// Use a replacer to map env vars like SERVER_PORT to server.port
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err = viper.ReadInConfig(); err != nil {
		// If config file is not found, we can proceed if env vars are set.
		// We will rely on default values and env vars.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	// Set default values
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.readTimeout", 15)
	v.SetDefault("server.writeTimeout", 15)
	v.SetDefault("server.idleTimeout", 60)
	v.SetDefault("security.corsEnabled", true)
	v.SetDefault("security.allowedOrigins", []string{"http://localhost:3000"})
	v.SetDefault("security.rateLimitRps", 10)
	v.SetDefault("security.rateLimitBurst", 20)

	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}

	return
}
