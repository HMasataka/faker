package faker

type Toml struct {
	DB     DB     `toml:"db"`
	Tables Tables `toml:"tables"`
}
