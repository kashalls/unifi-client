package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Host               string `env:"HOST,notEmpty"`
	Username           string `env:"USER" envDefault:""`
	Password           string `env:"PASSWORD" envDefault:""`
	ApiKey             string `env:"API_KEY" envDefault:""`
	Site               string `env:"SITE" envDefault:"default"`
	ExternalController bool   `env:"EXTERNAL_CONTROLLER" envDefault:"false"`
	SkipTLSVerify      bool   `env:"SKIP_TLS_VERIFY" envDefault:"false"`
	LongLogin          bool   `env:"LONG_LOGIN" envDefault:"true"`
}

func Parse() (*Config, error) {
	var config Config

	if err := env.ParseWithOptions(&config, env.Options{Prefix: "UNIFI_"}); err != nil {
		return nil, err
	}
	return &config, nil
}
