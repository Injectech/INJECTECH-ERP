package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	AppName          string        `mapstructure:"APP_NAME"`
	Env              string        `mapstructure:"APP_ENV"`
	HTTPPort         string        `mapstructure:"HTTP_PORT"`
	ReadTimeout      time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
	WriteTimeout     time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
	ShutdownTimeout  time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
	DatabaseURL      string        `mapstructure:"DATABASE_URL"`
	JWTAccessSecret  string        `mapstructure:"JWT_ACCESS_SECRET"`
	JWTRefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"`
	AccessTTL        time.Duration `mapstructure:"JWT_ACCESS_TTL"`
	RefreshTTL       time.Duration `mapstructure:"JWT_REFRESH_TTL"`
}

// Load initializes Viper and populates Config with defaults and env overrides.
func Load() (Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("APP_NAME", "erp-backend")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HTTP_READ_TIMEOUT", "5s")
	v.SetDefault("HTTP_WRITE_TIMEOUT", "10s")
	v.SetDefault("SHUTDOWN_TIMEOUT", "10s")
	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h") // 7d
	v.SetDefault("JWT_ACCESS_SECRET", "change-this-access")
	v.SetDefault("JWT_REFRESH_SECRET", "change-this-refresh")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// MustPort ensures port string is prefixed with ':' for http.ListenAndServe.
func (c Config) MustPort() string {
	if c.HTTPPort == "" {
		return ":8080"
	}
	if c.HTTPPort[0] == ':' {
		return c.HTTPPort
	}
	return ":" + c.HTTPPort
}
