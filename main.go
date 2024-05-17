package main

import (
	"mami/e-commerce/commons/logger"
	"mami/e-commerce/config"
	httpServer "mami/e-commerce/server/http"
	userModel "mami/e-commerce/user/model"

	"github.com/go-playground/validator/v10"
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

	err = db.AutoMigrate(&userModel.User{})
	if err != nil {
		logger.Fatal("Database migration fail", err)
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	httpSvr := httpServer.NewServer(db, validator)
	if err = httpSvr.Run(); err != nil {
		logger.Error(err)
	}

}
