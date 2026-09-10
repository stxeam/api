package domain

import "time"

type IConfig interface {
	Env() string
	AppName() string
	AppProtocol() string
	AppDomain() string
	APIVersion() string
	Port() string

	DBHost() string
	DBPort() string
	DBUser() string
	DBPassword() string
	DBName() string

	JWTAccessTokenSecret() string
	JWTRefreshTokenSecret() string
	JWTAccessTokenExpiry() int
	JWTRefreshTokenExpiry() int
	JWTAlgo() string

	CookiesSecure() bool
	CookiesSameSite() string

	CORSOrigins() []string
	CORSCredentials() bool

	DatabaseURL() string
	SSLMode() string

	JWTSecret() string
	Debug() bool

	S3Domain() string
	S3Port() string
	S3Region() string
	S3AccessKey() string
	S3SecretKey() string
	S3Bucket() string
	S3BucketExpiry() time.Duration
	S3PublicEndpoint() string

	AdminEmail() string
	AdminPassword() string
	NatsURL() string
}
