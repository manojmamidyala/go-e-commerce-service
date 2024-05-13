package main

import (
	"mami/e-commerce/commons/logger"
	"mami/e-commerce/config"
	httpServer "mami/e-commerce/server/http"
)

func main() {

	// load env variables just once in here so can be use in any other place
	env := config.LoadEnvVariables()
	logger.Initialize(env.Environment)

	// print the env variables
	logger.Infof("Port: %s\n", env.Environment)
	logger.Infof("Port: %d\n", env.HttpPort)

	db, err := config.NewDatabase(env.DatabaseURL)
	if err != nil {
		logger.Error("Database migration fail", err)
	}
	logger.Info(db.GetDB().Name())

	httpSvr := httpServer.NewServer(db)
	if err = httpSvr.Run(); err != nil {
		logger.Error(err)
	}

}
