package config

import (
	"strings"
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
	CorsOrigins      string        `mapstructure:"CORS_ORIGINS"`
	CookieSecure     bool          `mapstructure:"COOKIE_SECURE"`
}

// Load initializes Viper and populates Config with defaults and env overrides.
func Load() (Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		// Allow running from repo root where .env lives in backend/.
		v.SetConfigFile("backend/.env")
		_ = v.ReadInConfig()
	}
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
	v.SetDefault("CORS_ORIGINS", "http://localhost:5173")
	v.SetDefault("COOKIE_SECURE", false)

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

// CorsOriginList splits comma-separated CORS origins.
func (c Config) CorsOriginList() []string {
	origins := strings.Split(c.CorsOrigins, ",")
	out := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
