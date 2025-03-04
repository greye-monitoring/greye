package main

import (
	_ "greye/docs"
	"greye/pkg/factories"
	"greye/pkg/server"
	"log"
	"os"
	"strconv"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	//hostname := os.Getenv("HOSTNAME")
	//r, _ := regexp.MatchString("-0$", hostname)
	//time.Sleep(1 * time.Minute)

	localStartup := os.Getenv("LOCAL")
	file := "./config/config.json"
	logMesage := "Starting ClusterMonitor"
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

	//msg := fmt.Sprintf("ClusterMonitor is working on %s server", networkInfo.LocalIp)
	//fmt.Println(msg)
	srv := server.NewServer(appHandlers, clHandlers, configurator, role)

	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetOutput(os.Stdout)
	//log.SetLevel(log.InfoLevel)
	//
	//requestLogger := log.WithFields(log.Fields{"request_id": networkInfo.LocalIp, "user_ip": networkInfo.LocalIp})
	//requestLogger.Info("something happened on that request")
	//requestLogger.Warn("something not great happened")

	//go func() {
	//log.Println("Starting HTTP server on port", config.Server.Port)
	if err := srv.Run(strconv.Itoa(config.Server.Port)); err != nil {
		log.Fatal("HTTP server error:", err)
	}
	//}()
	err = srv.Run(strconv.Itoa(config.Server.Port))
	if err != nil {
		panic(err)
		//go func() {
		//	log.Println("Starting HTTPS server on port", config.Server.TlsPort)
		//	if err := srv.RunTls(strconv.Itoa(config.Server.TlsPort)); err != nil {
		//		log.Fatal("HTTPS server error:", err)
		//	}
		//}()

		// Prevent main from exiting
		//select {}
	}
}
