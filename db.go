package faker

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func NewConnection(cfg *DB) (*sql.DB, error) {
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
