package etc

import "github.com/caarlos0/env/v11"

type Env struct {

	// JWTSecretKey is the secret key used to sign the jwt token.
	JWTSecretKey string `env:"JWT_SECRET_KEY,notEmpty"`

	// JWTRefreshSecretKey is the secret key used to sign the refresh token.
	JWTRefreshSecretKey string `env:"JWT_REFRESH_SECRET_KEY,notEmpty"`
}

func LoadEnv() (*Env, error) {
	var envCfg Env
	err := env.Parse(&envCfg)
	if err != nil {
		return nil, err
	}

	return &envCfg, nil
}
