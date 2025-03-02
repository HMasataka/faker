package faker

import (
	"github.com/BurntSushi/toml"
	"github.com/HMasataka/gofiles"
)

type Toml struct {
	DB     DB     `toml:"db"`
	Tables Tables `toml:"tables"`
}

func NewTomlFromPath(path string) (*Toml, error) {
	b, err := gofiles.ReadFileAll(path)
	if err != nil {
		return nil, err
	}

	var t Toml
	if _, err := toml.Decode(string(b), &t); err != nil {
		return nil, err
	}

	return &t, nil
}
