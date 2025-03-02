package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/HMasataka/faker"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo/mutable"
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

	queue := tables.Tables
	fake := faker.NewFaker()

	for len(queue) > 0 {
		var deletable []int

		for i, table := range queue {
			if !fake.HasTables(table.Depends) {
				continue
			}

			record, err := fake.NewDummyRecord(table.Name, table.Column)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			columnNames := make(faker.ColumnNames, len(record))
			values := make([]any, len(record))

			keyValueIndex(record, func(i int, columnName faker.ColumnName, value any) {
				columnNames[i] = columnName
				values[i] = value
			})

			questions := repeat(len(table.Column), "?")
			question := strings.Join(questions, ",")

			query := fmt.Sprintf("INSERT INTO `%v` (%v) VALUES (%v)", table.Name, strings.Join(columnNames.ToStrings(), ","), question)

			log.Info().Str("query", query).Any("values", values).Send()

			if _, err := conn.ExecContext(context.Background(), query, values...); err != nil {
				log.Fatal().Err(err).Send()
			}

			deletable = append(deletable, i)
		}

		mutable.Reverse(deletable)

		for _, d := range deletable {
			queue = remove(queue, d)
		}
	}
}

func keyValueIndex[T comparable, V any](m map[T]V, fn func(idx int, key T, value V)) {
	idx := 0
	for key, value := range m {
		fn(idx, key, value)
		idx++
	}
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func repeat[T any](count int, v ...T) []T {
	l := count * len(v)
	xs := make([]T, l)

	for i := 0; i < l; i += len(v) {
		copy(xs[i:], v)
	}

	return xs
}
