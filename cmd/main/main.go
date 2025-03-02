package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/HMasataka/faker"
	"github.com/brianvoe/gofakeit/v7"
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

	seen := make(faker.DB)
	queue := tables.Tables

	for len(queue) > 0 {
		var deletable []int

		for i, table := range queue {
			columnNames := make(faker.ColumnNames, len(table.Column))
			columnValues := make([]any, len(table.Column))

			if !seen.HasAll(table.Depends) {
				continue
			}

			for i, column := range table.Column {
				columnNames[i] = column.Name
				record := make(faker.Record)

				switch column.ValueType {
				case "fakeit":
					value, err := gofakeit.Generate(column.Value)
					if err != nil {
						log.Fatal().Err(err).Send()
					}

					columnValues[i] = value
					record[column.Name] = value
				case "fk":
					sp := strings.Split(column.Value, ":")
					tableName, columnName := faker.TableName(sp[0]), faker.ColumnName(sp[1])
					value := seen[tableName][0][columnName]

					columnValues[i] = value
					record[column.Name] = value
				case "value":
					value := time.Now()
					columnValues[i] = value
					record[column.Name] = value
				}

				seen[table.Name] = append(seen[table.Name], record)
			}

			questions := repeat(len(table.Column), "?")
			question := strings.Join(questions, ",")

			query := fmt.Sprintf("INSERT INTO `%v` (%v) VALUES (%v)", table.Name, strings.Join(columnNames.ToStrings(), ","), question)

			log.Info().Str("query", query).Any("values", columnValues).Send()

			_, err := conn.ExecContext(context.Background(), query, columnValues...)
			if err != nil {
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

func remove(slice []*faker.Table, s int) []*faker.Table {
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
