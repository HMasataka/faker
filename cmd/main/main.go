package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/HMasataka/faker"
	"github.com/HMasataka/gofiles"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/rs/zerolog/log"
)

func main() {
	b, err := gofiles.ReadFileAll("data.toml")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	var t faker.Toml
	if _, err := toml.Decode(string(b), &t); err != nil {
		log.Fatal().Err(err).Send()
	}

	conn, err := faker.NewConnection(&t.DB)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer conn.Close()

	if err = conn.Ping(); err != nil {
		log.Fatal().Err(err).Send()
	}

	seen := make(map[string]map[string]any)

	for _, table := range t.Tables.Table {
		columnNames := make([]string, len(table.Column))
		columnValues := make([]any, len(table.Column))
		fmt.Println(table.Name, table.Depends)

		seenTable := make(map[string]any)
		seen[table.Name] = seenTable

		for i, column := range table.Column {
			columnNames[i] = column.Name

			switch column.ValueType {
			case "fakeit":
				value, err := gofakeit.Generate(column.Value)
				if err != nil {
					log.Fatal().Err(err).Send()
				}

				columnValues[i] = value
				seen[table.Name][column.Name] = value
			case "fk":
				sp := strings.Split(column.Value, ":")
				tableName, columnName := sp[0], sp[1]
				value := seen[tableName][columnName]

				columnValues[i] = value
				seen[table.Name][column.Name] = value
			case "value":
				value := time.Now()
				columnValues[i] = value
				seen[table.Name][column.Name] = value
			}
		}

		questions := repeat(len(table.Column), "?")
		question := strings.Join(questions, ",")

		query := fmt.Sprintf("INSERT INTO `%v` (%v) VALUES (%v)", table.Name, strings.Join(columnNames, ","), question)

		log.Info().Str("query", query).Any("values", columnValues).Send()

		_, err := conn.ExecContext(context.Background(), query, columnValues...)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}
}

func repeat[T any](count int, v ...T) []T {
	l := count * len(v)
	xs := make([]T, l)

	for i := 0; i < l; i += len(v) {
		copy(xs[i:], v)
	}

	return xs
}
