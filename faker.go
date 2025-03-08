package faker

import (
	"errors"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
)

const (
	FakeIt = "fakeit"
	FK     = "fk"
	Value  = "value"
)

func NewFaker() Faker {
	return Faker{
		db: make(DB),
	}
}

type Faker struct {
	db DB
}

func (f Faker) HasTable(tableName TableName) bool {
	return f.db.Has(tableName)
}

func (f Faker) HasTables(tableNames []TableName) bool {
	return f.db.HasAll(tableNames)
}

func (f Faker) GetTables() []TableName {
	return lo.Keys(f.db)
}

func (f Faker) NewDummyRecord(tableName TableName, columns []Column) (Record, error) {
	record := make(Record)

	for i := range columns {
		value, err := f.newDummyValue(tableName, columns[i])
		if err != nil {
			return nil, err
		}

		record[columns[i].Name] = value
	}

	f.db[tableName] = append(f.db[tableName], record)

	return record, nil
}

func (f Faker) newDummyValue(tableName TableName, column Column) (any, error) {
	switch column.ValueType {
	case FakeIt:
		return gofakeit.Generate(column.Value)
	case FK:
		sp := strings.Split(column.Value, ":")

		tn, columnName := TableName(sp[0]), ColumnName(sp[1])
		value := f.db[tn][0][columnName] // TODO 2個目以降のレコードのサポート

		return value, nil
	case Value:
		return f.buildValue(tableName, column)
	}

	return nil, errors.New("unsupported value type")
}

func (f Faker) buildValue(tableName TableName, column Column) (any, error) {
	switch column.Value {
	case "{now}":
		value := time.Now()
		return value, nil
	case "{increment}":
		if _, ok := f.db[tableName]; !ok {
			return 0, nil
		}

		if len(f.db[tableName]) == 0 {
			return 0, nil
		}

		record := f.db[tableName][len(f.db[tableName])-1]

		return record[column.Name].(int) + 1, nil
	}

	return column.Value, nil
}
