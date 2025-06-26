package etc

import "github.com/caarlos0/env/v11"

type Env struct {

	// JWTSecretKey is the secret key used to sign the jwt token.
	JWTSecretKey string `env:"JWT_SECRET_KEY,notEmpty"`

	// JWTRefreshSecretKey is the secret key used to sign the refresh token.
	JWTRefreshSecretKey string `env:"JWT_REFRESH_SECRET_KEY,notEmpty"`
}

var envCfg Env

func GetEnv() *Env {
	return &envCfg
}

func LoadEnv() error {
	return env.Parse(&envCfg)
}
