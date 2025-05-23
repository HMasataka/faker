package faker

import "github.com/caarlos0/env/v9"

type Config struct {
	DataBaseConfigFile string `env:"DATA_BASE_CONFIG_FILE" envDefault:"_example/db.toml"`
	TablesDirectory    string `env:"TABLES_DIRECTORY" envDefault:"_example/tables"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
