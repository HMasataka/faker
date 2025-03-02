package faker

import (
	"database/sql"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/HMasataka/gofiles"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func NewDataBaseConfig(path string) (DataBaseConfig, error) {
	b, err := gofiles.ReadFileAll(path)
	if err != nil {
		return DataBaseConfig{}, err
	}

	var db DataBaseConfig
	if _, err := toml.Decode(string(b), &db); err != nil {
		return DataBaseConfig{}, err
	}

	return db, nil
}

type DataBaseConfig struct {
	Name                 string `toml:"name"`
	User                 string `toml:"user"`
	Password             string `toml:"password"`
	Addr                 string `toml:"addr"`
	Net                  string `toml:"net"`
	ParseTime            bool   `toml:"parseTime"`
	Collation            string `toml:"collation"`
	Location             string `toml:"location"`
	AllowNativePasswords bool   `toml:"allowNativePasswords"`
}

func NewConnection(cfg *DataBaseConfig) (*sql.DB, error) {
	jst, err := time.LoadLocation(cfg.Location)
	if err != nil {
		return nil, err
	}

	c := mysql.Config{
		DBName:               cfg.Name,
		User:                 cfg.User,
		Passwd:               cfg.Password,
		Addr:                 cfg.Addr,
		Net:                  cfg.Net,
		ParseTime:            cfg.ParseTime,
		Collation:            cfg.Collation,
		Loc:                  jst,
		AllowNativePasswords: cfg.AllowNativePasswords,
	}

	log.Info().Str("DSN", c.FormatDSN()).Send()

	conn, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
