package etc

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Certs struct {
		CAFile   string `yaml:"ca_file"`
		KeyFile  string `yaml:"key_file"`
		CertFile string `yaml:"cert_file"`
	}

	Settings struct {
		// UsernamePattern defines the verification pattern of a valid username.
		// It must be a valid regular expression.
		UsernamePattern string `yaml:"username_pattern"`

		// PasswordPattern defines the verification pattern of a valid password.
		// It must be a valid regular expression.
		PasswordPattern string `yaml:"password_pattern"`

		// JWT defines the JWT configuration.
		JWT struct {
			// TokenExpireTime defines the expiration time of a token in seconds.
			TokenExpireTime int `yaml:"token_expire_time"`

			// TokenRefreshTime defines the refresh time of a token in seconds.
			TokenRefreshTime int `yaml:"token_refresh_time"`
		} `yaml:"jwt"`

		// DefaultValue defines the default values for some configurations.
		DefaultValue struct {
			// UserState defines the default user state when a user is created.
			UserState int8 `yaml:"user_state"`
		} `yaml:"default_value"`
	}

	App struct {
		SnowflakeID int64     `yaml:"snowflake_id"`
		Certs       *Certs    `yaml:"certs"`
		Settings    *Settings `yaml:"settings"`
	}
	MySQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DB       string `yaml:"db"`
	}
	Store struct {
		Mysql map[string]*MySQL `yaml:"mysql"`
	}
	Config struct {
		App   *App   `yaml:"app"`
		Store *Store `yaml:"store"`
	}
)

var conf Config

func GetConfig() Config {
	return conf
}

func Load(file string) error {
	bs, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, &conf)
}
