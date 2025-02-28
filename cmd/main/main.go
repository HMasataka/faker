package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/HMasataka/faker"
	"github.com/HMasataka/gofiles"
	"github.com/rs/zerolog/log"
)

func main() {
	b, err := gofiles.ReadFileAll("data.toml")
	if err != nil {
		panic(err)
	}

	var t faker.Toml
	if _, err := toml.Decode(string(b), &t); err != nil {
		log.Fatal().Err(err).Send()
	}

	fmt.Printf("%+v", t)
}
