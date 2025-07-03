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
		Validation struct {
			// UsernamePattern defines the verification pattern of a valid username.
			UsernamePattern struct {
				MinLength int `yaml:"min_length"`
				MaxLength int `yaml:"max_length"`
			} `yaml:"username_pattern"`

			// PasswordPattern defines the verification pattern of a valid password.
			PasswordPattern struct {
				MinLength           int  `yaml:"min_length"`
				RequireUpperCase    bool `yaml:"require_upper_case"`
				RequireLowerCase    bool `yaml:"require_lower_case"`
				RequireDigit        bool `yaml:"require_digit"`
				RequireSpecialChars bool `yaml:"require_special_chars"`
			} `yaml:"password_pattern"`
		} `yaml:"validation"`

		// JWT defines the JWT configuration.
		JWT struct {
			// TokenExpireTime defines the expiration time of a token in seconds.
			TokenExpireTime int64 `yaml:"token_expire_time"`

			// TokenRefreshTime defines the refresh time of a token in seconds.
			TokenRefreshTime int64 `yaml:"token_refresh_time"`
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
		Debug    bool   `yaml:"debug"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DB       string `yaml:"db"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	}
	Store struct {
		Mysql map[string]*MySQL `yaml:"mysql"`
		Redis *Redis            `yaml:"redis"`
	}
	Config struct {
		App   *App   `yaml:"app"`
		Store *Store `yaml:"store"`
	}
)

func Load(file string) (*Config, error) {
	var conf Config

	bs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bs, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
