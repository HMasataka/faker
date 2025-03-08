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

func (f Faker) NewDummyRecords(table *Table) (*Records, error) {
	records := Records{
		ColumnNames: table.Columns.ToColumnNames(),
	}

	for i := range 1 {
		values := make([]any, len(table.Columns))
		records.Values = append(records.Values, values)

		for columIndex := range table.Columns {
			value, err := f.newDummyValue(i, table.Name, table.Columns[columIndex])
			if err != nil {
				return nil, err
			}

			records.Values[i][columIndex] = value
		}
	}

	f.db[table.Name] = records

	return &records, nil
}

func (f Faker) newDummyValue(idx int, tableName TableName, column Column) (any, error) {
	switch column.ValueType {
	case FakeIt:
		return gofakeit.Generate(column.Value)
	case FK:
		sp := strings.Split(column.Value, ":")

		tn, columnName := TableName(sp[0]), ColumnName(sp[1])
		value := f.db[tn].GetByColumnName(f.getDependsRecordFromIndex(idx), columnName)

		return value, nil
	case Value:
		return f.buildValue(tableName, column)
	}

	return nil, errors.New("unsupported value type")
}

func (f Faker) getDependsRecordFromIndex(idx int) int {
	return 0
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

		if f.db[tableName].Len() == 0 {
			return 0, nil
		}

		record, err := f.db[tableName].GetLast()
		if err != nil {
			return 0, err
		}

		columnNames := f.db[tableName].ColumnNames

		return record[columnNames.IndexOf(column.Name)].(int) + 1, nil
	}

	return column.Value, nil
}
