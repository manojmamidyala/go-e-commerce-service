package main

import (
	"fmt"

	"mami/e-commerce/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// load env variables just once in here so can be use in any other place
	env := config.LoadEnvVariables()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Info().Msg(env.Environment)

	// print the env variables
	fmt.Printf("Port: %s\n", env.Environment)
	fmt.Printf("Port: %d\n", env.HttpPort)

}
