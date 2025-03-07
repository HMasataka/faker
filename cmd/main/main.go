package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/HMasataka/faker"
	"github.com/HMasataka/perseus"
	"github.com/HMasataka/ruin"
	"github.com/rs/zerolog/log"
)

func newConn() (*sql.DB, error) {
	db, err := faker.NewDataBaseConfig("db.toml")
	if err != nil {
		return nil, err
	}

	conn, err := faker.NewConnection(&db)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func main() {
	conn, err := newConn()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer conn.Close()

	if err = conn.Ping(); err != nil {
		log.Fatal().Err(err).Send()
	}

	tables, err := faker.NewTables("tables")
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

		record, err := fake.NewDummyRecord(table.Name, table.Column)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		columnNames := make(faker.ColumnNames, len(record))
		values := make([]any, len(record))

		withIndex := perseus.WithIndex(record)
		for kv := range withIndex {
			columnNames[kv.Index] = kv.Key
			values[kv.Index] = kv.Value
		}

		query := fmt.Sprintf("INSERT INTO `%v` (%v) VALUES (%v)", table.Name, strings.Join(columnNames.ToStrings(), ","), faker.BuildQuestionMarks(len(columnNames)))

		log.Info().Str("query", query).Any("values", values).Send()

		if _, err := conn.ExecContext(context.Background(), query, values...); err != nil {
			log.Fatal().Err(err).Send()
		}
	}
}
