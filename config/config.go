package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Payment  PaymentConfig
	Server   ServerConfig
	CORS     CORSConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string
	Environment string // dev, staging, production
	Version     string
	LogLevel    string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	TLS             TLSConfig
}

// TLSConfig holds TLS/HTTPS configuration
type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Database       string
	SSLMode        string // disable, require, verify-ca, verify-full
	MaxConnections int
	MaxIdleConns   int
	MaxLifetime    time.Duration
	MaxIdleTime    time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Database int
	PoolSize int
	TTL      time.Duration
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret           string
	JWTExpiration       time.Duration
	RefreshTokenExpiry  time.Duration
	BcryptCost          int
	SessionCookieName   string
	SessionCookieSecure bool
	SessionCookieDomain string
}

// PaymentConfig holds payment gateway configuration
type PaymentConfig struct {
	Provider   string // stripe, paypal, etc.
	PublicKey  string
	SecretKey  string
	WebhookKey string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")
		v.AddConfigPath(".")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		// Config file is optional; continue if not found
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Environment variables override config file
	v.SetEnvPrefix("ECOMMERCE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "ecommerce")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.loglevel", "info")

	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.readtimeout", "15s")
	v.SetDefault("server.writetimeout", "15s")
	v.SetDefault("server.shutdowntimeout", "30s")
	v.SetDefault("server.tls.enabled", false)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.database", "ecommerce")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.maxconnections", 25)
	v.SetDefault("database.maxidleconns", 5)
	v.SetDefault("database.maxlifetime", "5m")
	v.SetDefault("database.maxidletime", "10m")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.database", 0)
	v.SetDefault("redis.poolsize", 10)
	v.SetDefault("redis.ttl", "1h")

	// Auth defaults
	v.SetDefault("auth.jwtsecret", "change-me-in-production")
	v.SetDefault("auth.jwtexpiration", "15m")
	v.SetDefault("auth.refreshtokenexpiry", "7d")
	v.SetDefault("auth.bcryptcost", 12)
	v.SetDefault("auth.sessioncookiename", "session")
	v.SetDefault("auth.sessioncookiesecure", false)
	v.SetDefault("auth.sessioncookiedomain", "")

	// Payment defaults
	v.SetDefault("payment.provider", "stripe")
	v.SetDefault("payment.publickey", "")
	v.SetDefault("payment.secretkey", "")
	v.SetDefault("payment.webhookkey", "")

	// CORS defaults
	v.SetDefault("cors.allowedorigins", []string{"*"})
	v.SetDefault("cors.allowedmethods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowedheaders", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"})
	v.SetDefault("cors.exposedheaders", []string{})
	v.SetDefault("cors.allowcredentials", true)
	v.SetDefault("cors.maxage", 300)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate environment
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.App.Environment] {
		return fmt.Errorf("invalid environment: %s (must be development, staging, or production)", c.App.Environment)
	}

	// Validate database
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate auth in production
	if c.App.Environment == "production" {
		if c.Auth.JWTSecret == "change-me-in-production" {
			return fmt.Errorf("JWT secret must be changed in production")
		}
		if !c.Server.TLS.Enabled {
			return fmt.Errorf("TLS must be enabled in production")
		}
	}

	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// DatabaseDSN returns the PostgreSQL connection string
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}

// RedisAddr returns the Redis address
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// ServerAddr returns the HTTP server address
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
