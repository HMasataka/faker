package faker

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/HMasataka/gofiles"
)

func NewTables(dir string) (Tables, error) {
	filePaths, err := gofiles.ListFilesYield(dir)
	if err != nil {
		return Tables{}, err
	}

	var data string

	for filePath := range filePaths {
		b, err := gofiles.ReadFileAll(filePath)
		if err != nil {
			return Tables{}, err
		}

		data += string(b)
	}

	fmt.Println(data)

	var t Tables
	if _, err := toml.Decode(data, &t); err != nil {
		return Tables{}, err
	}

	return t, nil
}

type Tables struct {
	Tables []*Table `toml:"tables"`
}

type WantType string

type Table struct {
	Name    TableName   `toml:"name"`
	Want    int         `toml:"want"`
	Depends []TableName `toml:"depends"`
	Columns Columns     `toml:"columns"`
}

type Column struct {
	Name      ColumnName `toml:"name"`
	ValueType string     `toml:"valueType"`
	Value     string     `toml:"value"`
}

type Columns []Column

func (c Columns) ToColumnNames() ColumnNames {
	names := make(ColumnNames, len(c))

	for i := range c {
		names[i] = c[i].Name
	}

	return names
}
