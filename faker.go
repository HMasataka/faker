package faker

import (
	"errors"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
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
		value, err := f.newDummyValue(columns[i].ValueType, columns[i].Value)
		if err != nil {
			return nil, err
		}

		record[columns[i].Name] = value
	}

	f.db[tableName] = append(f.db[tableName], record)

	return record, nil
}

func (f Faker) newDummyValue(valueType, keyword string) (any, error) {
	switch valueType {
	case "fakeit":
		return gofakeit.Generate(keyword)
	case "fk":
		sp := strings.Split(keyword, ":")

		tableName, columnName := TableName(sp[0]), ColumnName(sp[1])
		value := f.db[tableName][0][columnName] // TODO 2個目以降のレコードのサポート

		return value, nil
	case "value":
		value := time.Now()
		return value, nil
	}

	return nil, errors.New("unsupported value type")
}
