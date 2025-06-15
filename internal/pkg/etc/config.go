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

	App struct {
		SnowflakeID int64  `yaml:"snowflake_id"`
		Certs       *Certs `yaml:"certs"`
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
