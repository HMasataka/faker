package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/HMasataka/faker"
	"github.com/HMasataka/ruin"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func newConn(cfg *faker.Config) (*sql.DB, error) {
	db, err := faker.NewDataBaseConfig(cfg.DataBaseConfigFile)
	if err != nil {
		return nil, err
	}

	conn, err := faker.NewConnection(&db)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		log.Fatal().Err(err).Send()
	}

	return conn, nil
}

func main() {
	cfg, err := faker.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	conn, err := newConn(cfg)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer conn.Close()

	tables, err := faker.NewTables(cfg.TablesDirectory)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	queue := ruin.New(tables.Tables)
	fake := faker.NewFaker()

	for !queue.IsEmpty() {
		table, err := queue.Pop()
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		if !fake.HasTables(table.Depends) {
			queue.Push(table)
			continue
		}

		record, err := fake.NewDummyRecords(table)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		columnNames := record.ColumnNames.ToStrings()

		query := fmt.Sprintf("INSERT INTO `%v` (%v) VALUES %v", table.Name, strings.Join(columnNames, ","), faker.BuildQuestionMarks(len(record.Values), len(columnNames)))

		log.Info().Str("query", query).Any("values", record.Values).Send()

		if _, err := conn.ExecContext(context.Background(), query, lo.Flatten(record.Values)...); err != nil {
			log.Fatal().Err(err).Send()
		}
	}
}
