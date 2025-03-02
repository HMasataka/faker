package faker

import (
	"github.com/BurntSushi/toml"
	"github.com/HMasataka/gofiles"
)

func NewTables(path string) (Tables, error) {
	b, err := gofiles.ReadFileAll(path)
	if err != nil {
		return Tables{}, err
	}

	var t Tables
	if _, err := toml.Decode(string(b), &t); err != nil {
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
