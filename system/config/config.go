package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port  string
	Env   string
}

type DatabaseConfig struct {
	PostgresURL string
}

type RedisConfig struct {
	Address  string
	Password string
}

type ElasticsearchConfig struct {
	URL string
}

type Neo4jConfig struct {
	URL      string
	User     string
	Password string
}

type JWTConfig struct {
	Secret string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

type MonitoringConfig struct {
	LokiHost string
}

type Config struct {
	App           AppConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Elasticsearch ElasticsearchConfig
	Neo4j         Neo4jConfig
	JWT           JWTConfig
	OAuth         OAuthConfig
	Monitoring    MonitoringConfig
	CookieSecure  bool
	FrontendURL   string
}

func (c *Config) IsDev() bool {
	return c.App.Env == "development"
}

func (c *Config) IsProd() bool {
	return c.App.Env == "production"
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Port:  optional("PORT", "8000"),
			Env:   optional("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			PostgresURL: os.Getenv("PG_URL"),
		},
		Redis: RedisConfig{
			Address:  os.Getenv("RDS_ADDRESS"),
			Password: os.Getenv("RDS_PASSWORD"),
		},
		Elasticsearch: ElasticsearchConfig{
			URL: os.Getenv("ES_URL"),
		},
		Neo4j: Neo4jConfig{
			URL:      os.Getenv("N4J_URL"),
			User:     os.Getenv("N4J_USER"),
			Password: os.Getenv("N4J_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		},
		Monitoring: MonitoringConfig{
			LokiHost: os.Getenv("LOKI_HOST"),
		},
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
		FrontendURL:  os.Getenv("FRONTEND_URL"),
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var missing []string

	check := func(val string, name string) {
		if val == "" {
			missing = append(missing, name)
		}
	}

	check(c.Database.PostgresURL, "PG_URL")
	check(c.Redis.Address, "RDS_ADDRESS")
	check(c.JWT.Secret, "JWT_SECRET")
	check(c.OAuth.GoogleClientID, "GOOGLE_CLIENT_ID")
	check(c.OAuth.GoogleClientSecret, "GOOGLE_CLIENT_SECRET")
	check(c.OAuth.GoogleRedirectURL, "GOOGLE_REDIRECT_URL")
	check(c.FrontendURL, "FRONTEND_URL")

	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}
	return nil
}

func (c *Config) Print() {
	slog.Info("configuration loaded",
		"app.port", c.App.Port,
		"app.env", c.App.Env,
		"database.postgres_url", mask(c.Database.PostgresURL),
		"redis.address", c.Redis.Address,
		"redis.password_set", c.Redis.Password != "",
		"elasticsearch.url", c.Elasticsearch.URL,
		"neo4j.url", c.Neo4j.URL,
		"jwt.secret_set", c.JWT.Secret != "",
		"oauth.google_client_id", c.OAuth.GoogleClientID,
		"oauth.google_redirect_url", c.OAuth.GoogleRedirectURL,
		"monitoring.loki_host", c.Monitoring.LokiHost,
		"cookie_secure", c.CookieSecure,
		"frontend_url", c.FrontendURL,
	)
}

func optional(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mask(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "..." + s[len(s)-4:]
}
