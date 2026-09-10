package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go-starter/internal/config/domain"
)

type EnvValidator struct{}

func NewEnvValidator() *EnvValidator {
	return &EnvValidator{}
}

func (v *EnvValidator) Validate(cfg domain.IConfig) error {
	var errs []string

	if strings.TrimSpace(cfg.Env()) == "" {
		errs = append(errs, "ENV must not be empty")
	}
	if strings.TrimSpace(cfg.AppName()) == "" {
		errs = append(errs, "APP_NAME must not be empty")
	}
	proto := strings.TrimSpace(cfg.AppProtocol())
	if proto == "" {
		errs = append(errs, "APP_PROTOCOL must not be empty")
	} else if proto != "http" && proto != "https" {
		errs = append(errs, fmt.Sprintf("APP_PROTOCOL must be http or https, got '%s'", proto))
	}
	if strings.TrimSpace(cfg.AppDomain()) == "" {
		errs = append(errs, "APP_DOMAIN must not be empty")
	} else {
		u := fmt.Sprintf("%s://%s", proto, cfg.AppDomain())
		if _, err := url.Parse(u); err != nil {
			errs = append(errs, fmt.Sprintf("APP_DOMAIN is not valid '%s': %v", cfg.AppDomain(), err))
		}
	}
	if strings.TrimSpace(cfg.APIVersion()) == "" {
		errs = append(errs, "API_VERSION must not be empty")
	}
	if strings.TrimSpace(cfg.Port()) == "" {
		errs = append(errs, "PORT must not be empty")
	} else if _, err := strconv.Atoi(cfg.Port()); err != nil {
		errs = append(errs, fmt.Sprintf("PORT must be numeric, got '%s'", cfg.Port()))
	}
	if strings.TrimSpace(cfg.NatsURL()) == "" {
		errs = append(errs, "NATS_URL must not be empty")
	} else if _, err := url.Parse(cfg.NatsURL()); err != nil {
		errs = append(errs, fmt.Sprintf("NATS_URL is not a valid URL '%s': %v", cfg.NatsURL(), err))
	} else if !strings.HasPrefix(cfg.NatsURL(), "nats://") {
		errs = append(errs, fmt.Sprintf("NATS_URL must start with nats://, got '%s'", cfg.NatsURL()))
	}

	if strings.TrimSpace(cfg.S3Domain()) == "" {
		errs = append(errs, "S3_DOMAIN must not be empty")
	}
	if strings.TrimSpace(cfg.S3Port()) == "" {
		errs = append(errs, "S3_PORT must not be empty")
	} else if _, err := strconv.Atoi(cfg.S3Port()); err != nil {
		errs = append(errs, fmt.Sprintf("S3_PORT must be numeric, got '%s'", cfg.S3Port()))
	}
	if strings.TrimSpace(cfg.S3Bucket()) == "" {
		errs = append(errs, "S3_BUCKET must not be empty")
	}
	if strings.TrimSpace(cfg.S3AccessKey()) == "" {
		errs = append(errs, "S3_ACCESS_KEY must not be empty")
	}
	if strings.TrimSpace(cfg.S3SecretKey()) == "" {
		errs = append(errs, "S3_SECRET_KEY must not be empty")
	}
	if strings.TrimSpace(cfg.S3Region()) == "" {
		errs = append(errs, "S3_REGION must not be empty")
	}
	if cfg.S3BucketExpiry() <= 0 {
		errs = append(errs, fmt.Sprintf("S3_BUCKET_EXPIRY must be positive, got '%v'", cfg.S3BucketExpiry()))
	}
	if s3Endpoint := cfg.S3PublicEndpoint(); strings.TrimSpace(s3Endpoint) != "" {
		if _, err := url.Parse(s3Endpoint); err != nil {
			errs = append(errs, fmt.Sprintf("S3_PUBLIC_ENDPOINT is not a valid URL '%s': %v", s3Endpoint, err))
		}
	}

	if strings.TrimSpace(cfg.DBHost()) == "" {
		errs = append(errs, "DB_HOST must not be empty")
	}
	if strings.TrimSpace(cfg.DBPort()) == "" {
		errs = append(errs, "DB_PORT must not be empty")
	} else if _, err := strconv.Atoi(cfg.DBPort()); err != nil {
		errs = append(errs, fmt.Sprintf("DB_PORT must be numeric, got '%s'", cfg.DBPort()))
	}
	if strings.TrimSpace(cfg.DBUser()) == "" {
		errs = append(errs, "DB_USER must not be empty")
	}
	if strings.TrimSpace(cfg.DBPassword()) == "" {
		errs = append(errs, "DB_PASSWORD must not be empty")
	}
	if strings.TrimSpace(cfg.DBName()) == "" {
		errs = append(errs, "DB_NAME must not be empty")
	}
	if strings.TrimSpace(cfg.SSLMode()) == "" {
		errs = append(errs, "DB_SSLMODE must not be empty")
	}

	if strings.TrimSpace(cfg.JWTAccessTokenSecret()) == "" {
		errs = append(errs, "JWT_ACCESS_TOKEN_SECRET must not be empty")
	}
	if strings.TrimSpace(cfg.JWTRefreshTokenSecret()) == "" {
		errs = append(errs, "JWT_REFRESH_TOKEN_SECRET must not be empty")
	}
	if cfg.JWTAccessTokenExpiry() <= 0 {
		errs = append(errs, fmt.Sprintf("JWT_ACCESS_TOKEN_EXPIRY must be positive, got '%d'", cfg.JWTAccessTokenExpiry()))
	}
	if cfg.JWTRefreshTokenExpiry() <= 0 {
		errs = append(errs, fmt.Sprintf("JWT_REFRESH_TOKEN_EXPIRY must be positive, got '%d'", cfg.JWTRefreshTokenExpiry()))
	}
	if strings.TrimSpace(cfg.JWTAlgo()) == "" {
		errs = append(errs, "JWT_ALGO must not be empty")
	}
	if strings.TrimSpace(cfg.CookiesSameSite()) == "" {
		errs = append(errs, "COOKIES_SAME_SITE must not be empty")
	}
	if strings.TrimSpace(cfg.AdminEmail()) == "" {
		errs = append(errs, "ADMIN_EMAIL must not be empty")
	}
	if strings.TrimSpace(cfg.AdminPassword()) == "" {
		errs = append(errs, "ADMIN_PASSWORD must not be empty")
	}

	if len(errs) > 0 {
		return fmt.Errorf("env validation failed: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (v *EnvValidator) ValidateEnv(cfg domain.IConfig) error {
	return v.Validate(cfg)
}
