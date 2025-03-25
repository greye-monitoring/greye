package main

import (
	_ "greye/docs"
	"greye/pkg/factories"
	"greye/pkg/server"
	"log"
	"os"
	"strconv"
)

func main() {

	localStartup := os.Getenv("LOCAL")
	file := "./config/config.json"
	logMesage := "Starting greye"
	if len(localStartup) != 0 {
		file = "./config/env-" + os.Getenv("LOCAL") + ".json"
	}
	factory := factories.NewFactory(
		file,
	)
	configurator := factory.InitializeConfigurator()
	config, err := configurator.GetConfig()
	if err != nil {
		panic(err)
	}
	factory.InitializeRole()
	factory.InitializeScheduler()
	logger := factory.InitializeLogger()
	factory.InitializeHttpClient(logger)
	factory.InitializeImportService()
	logger.Info(logMesage)
	factory.InitializeNotification()
	appHandlers := factory.BuildAppHandlers()
	role := factory.InitializeRole()
	clHandlers := factory.BuildClusterHandlers()

	srv := server.NewServer(appHandlers, clHandlers, configurator, role)

	if err := srv.Run(strconv.Itoa(config.Server.Port)); err != nil {
		log.Fatal("HTTP server error:", err)
	}

	err = srv.Run(strconv.Itoa(config.Server.Port))
	if err != nil {
		panic(err)
	}
}
