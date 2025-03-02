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

type Table struct {
	Name    string   `toml:"name"`
	Depends []string `toml:"depends"`
	Column  []Column `toml:"column"`
}

type Column struct {
	Name      string `toml:"name"`
	ValueType string `toml:"valueType"`
	Value     string `toml:"value"`
}
