package main

import (
	"github.com/BurntSushi/toml"
	"github.com/HMasataka/faker"
	"github.com/HMasataka/gofiles"
	"github.com/rs/zerolog/log"
)

func main() {
	b, err := gofiles.ReadFileAll("data.toml")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	var t faker.Toml
	if _, err := toml.Decode(string(b), &t); err != nil {
		log.Fatal().Err(err).Send()
	}

	conn, err := faker.NewConnection(&t.DB)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer conn.Close()

	if err = conn.Ping(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
