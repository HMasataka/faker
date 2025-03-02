package faker

type Tables struct {
	Table []Table `toml:"table"`
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
