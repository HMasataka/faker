package faker

type DB struct {
	Name                 string `toml:"name"`
	User                 string `toml:"user"`
	Password             string `toml:"password"`
	Addr                 string `toml:"addr"`
	Net                  string `toml:"net"`
	ParseTime            bool   `toml:"parseTime"`
	Collation            string `toml:"collation"`
	Location             string `toml:"location"`
	AllowNativePasswords bool   `toml:"allowNativePasswords"`
}
