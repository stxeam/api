package infrastructure

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go-starter/internal/config/domain"
)

type ConfigAdapter struct{}

func NewConfigAdapter() domain.IConfig {
	return &ConfigAdapter{}
}

func (c *ConfigAdapter) Env() string           { return os.Getenv("ENV") }
func (c *ConfigAdapter) AppName() string       { return os.Getenv("APP_NAME") }
func (c *ConfigAdapter) AppProtocol() string   { return os.Getenv("APP_PROTOCOL") }
func (c *ConfigAdapter) AppDomain() string     { return os.Getenv("APP_DOMAIN") }
func (c *ConfigAdapter) APIVersion() string    { return os.Getenv("API_VERSION") }
func (c *ConfigAdapter) Port() string          { return os.Getenv("PORT") }
func (c *ConfigAdapter) AdminEmail() string    { return os.Getenv("ADMIN_EMAIL") }
func (c *ConfigAdapter) AdminPassword() string { return os.Getenv("ADMIN_PASSWORD") }
func (c *ConfigAdapter) DBHost() string        { return os.Getenv("DB_HOST") }
func (c *ConfigAdapter) DBPort() string        { return os.Getenv("DB_PORT") }
func (c *ConfigAdapter) DBUser() string        { return os.Getenv("DB_USER") }
func (c *ConfigAdapter) DBPassword() string    { return os.Getenv("DB_PASSWORD") }
func (c *ConfigAdapter) DBName() string        { return os.Getenv("DB_NAME") }
func (c *ConfigAdapter) JWTAccessTokenSecret() string {
	return os.Getenv("JWT_ACCESS_TOKEN_SECRET")
}
func (c *ConfigAdapter) JWTRefreshTokenSecret() string {
	return os.Getenv("JWT_REFRESH_TOKEN_SECRET")
}
func (c *ConfigAdapter) JWTAccessTokenExpiry() int {
	v := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY")
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return 0
}
func (c *ConfigAdapter) JWTRefreshTokenExpiry() int {
	v := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY")
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return 0
}
func (c *ConfigAdapter) JWTAlgo() string { return os.Getenv("JWT_ALGO") }
func (c *ConfigAdapter) CookiesSecure() bool {
	v := os.Getenv("COOKIES_SECURE")
	switch strings.ToLower(v) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	}
	return false
}
func (c *ConfigAdapter) CookiesSameSite() string { return os.Getenv("COOKIES_SAME_SITE") }
func (c *ConfigAdapter) CORSOrigins() []string {
	raw := os.Getenv("CORS_ORIGINS")
	raw = strings.Trim(raw, "[]")
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(strings.Trim(p, `"`))
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}
func (c *ConfigAdapter) CORSCredentials() bool {
	v := os.Getenv("CORS_CREDENTIALS")
	switch strings.ToLower(v) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	}
	return false
}
func (c *ConfigAdapter) DatabaseURL() string {
	ssl := c.SSLMode()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost(), c.DBPort(), c.DBUser(), c.DBPassword(), c.DBName(), ssl)
	return dsn
}
func (c *ConfigAdapter) SSLMode() string   { return os.Getenv("DB_SSLMODE") }
func (c *ConfigAdapter) JWTSecret() string { return c.JWTAccessTokenSecret() }
func (c *ConfigAdapter) Debug() bool       { return c.Env() == "dev" }

func (c *ConfigAdapter) NatsURL() string { return os.Getenv("NATS_URL") }

func (c *ConfigAdapter) S3Domain() string    { return os.Getenv("S3_DOMAIN") }
func (c *ConfigAdapter) S3Port() string      { return os.Getenv("S3_PORT") }
func (c *ConfigAdapter) S3Region() string    { return os.Getenv("S3_REGION") }
func (c *ConfigAdapter) S3AccessKey() string { return os.Getenv("S3_ACCESS_KEY") }
func (c *ConfigAdapter) S3SecretKey() string { return os.Getenv("S3_SECRET_KEY") }
func (c *ConfigAdapter) S3Bucket() string    { return os.Getenv("S3_BUCKET") }
func (c *ConfigAdapter) S3BucketExpiry() time.Duration {
	v := os.Getenv("S3_BUCKET_EXPIRY")
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return time.Duration(n) * time.Second
}
func (c *ConfigAdapter) S3PublicEndpoint() string {
	return fmt.Sprintf("%s://%s", c.AppProtocol(), c.AppDomain())
}
