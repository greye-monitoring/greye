package factories

import (
	"fmt"
	"greye/internal/application/application"
	applicationHandler "greye/internal/application/infrastructure/handlers"
	application2 "greye/internal/cluster/application"
	clusterHandler "greye/internal/cluster/infrastructure/handlers"
	authenticationApp "greye/pkg/authentication/application"
	portsAuth "greye/pkg/authentication/domain/ports"
	clientApp "greye/pkg/client/application"
	ports2 "greye/pkg/client/domain/ports"
	configApp "greye/pkg/config/application"
	configRepo "greye/pkg/config/infrastructure/repositories"
	k8s "greye/pkg/importProcess/application"
	loggerApp "greye/pkg/logging/application"
	metricsApp "greye/pkg/metrics/application"
	metricsPort "greye/pkg/metrics/domain/ports"
	notificationApp "greye/pkg/notification/application"
	"greye/pkg/notification/domain/ports"
	"greye/pkg/role/domain/models"
	schedulerApp "greye/pkg/scheduler/application"
	"greye/pkg/server"
	models2 "greye/pkg/type/domain/models"
	valApp "greye/pkg/validator/application"
	"os"
	"regexp"
)

const MongoClientTimeout = 10

type Factory struct {
	//Variables
	configFilePath string
	role           models.Role

	importService  *k8s.ImportProcessApplication
	configurator   *configApp.ConfigService
	validator      *valApp.Validator
	logger         *loggerApp.Logger
	httpClient     *clientApp.HttpApplication
	scheduler      *schedulerApp.Job
	notification   map[string]ports.Sender
	protocol       map[string]ports2.MonitoringMethod
	metricApp      *metricsPort.MetricPorts
	metricCluster  *metricsPort.MetricPorts
	authentication map[string]portsAuth.Authentication
}

func NewFactory(configFilePath string) *Factory {
	return &Factory{
		configFilePath: configFilePath,
	}
}

func (f *Factory) InitializeValidator() *valApp.Validator {
	if f.validator == nil {
		app := valApp.NewValidator()
		f.validator = app
		return app
	}
	return f.validator
}

func (f *Factory) InitializeConfigurator() *configApp.ConfigService {
	if f.configurator == nil {
		validator := f.InitializeValidator()
		path := f.configFilePath

		repo := configRepo.NewJSONRepository(path)
		//log := loggerApp.NewLogger()
		app := configApp.NewConfigService(repo, validator)
		err := app.Config()
		if err != nil {
			panic(err)
		}
		f.configurator = app
		return app
	}
	return f.configurator
}

func (f *Factory) InitializeLogger() *loggerApp.Logger {
	if f.logger == nil {
		configurator := f.InitializeConfigurator()
		config, _ := configurator.GetConfig()
		logLevel := config.Server.LogLevel
		logs := loggerApp.NewLogger(logLevel)
		f.logger = logs
		return logs
	}
	return f.logger
}

func (f *Factory) InitializeHttpClient(logHandler *loggerApp.Logger) *clientApp.HttpApplication {
	if f.httpClient == nil {
		httpClient := clientApp.NewHttpApplication(logHandler)
		f.httpClient = httpClient
		return httpClient
	}
	return f.httpClient
}

func (f *Factory) InitializeScheduler() *schedulerApp.Job {
	if f.httpClient == nil {
		sched := schedulerApp.NewJob()
		f.scheduler = sched
		return sched
	}
	return f.scheduler
}

func (f *Factory) InitializeRole() models.Role {
	if f.role != "" {
		return f.role
	}
	configurator := f.InitializeConfigurator()
	log := f.InitializeLogger()
	config, _ := configurator.GetConfig()
	serverName := config.Server.ApplicationName
	hostname := os.Getenv("HOSTNAME")

	regexPattern := fmt.Sprintf(`^%s-0|%s(:[0-9]0[0-9]0)$`, serverName, serverName)

	r, err := regexp.MatchString(regexPattern, hostname)

	if err != nil || !r {
		nClusterMonitor := config.Server.NumberGreye
		regexPattern := fmt.Sprintf(`^%s(-([0-%d]))|%s(:808[0-3])$`, serverName, nClusterMonitor, serverName)

		r, err = regexp.MatchString(regexPattern, hostname)

		if err != nil || !r {
			fmt.Println(os.Environ())
			panic("the env variable 'HOSTNAME' must be set")
		}
		log.Info("I'm worker %s", hostname)
		var rt models.Role = "worker"
		f.role = rt
		return rt
	}
	log.Info("I'm controller %s", hostname)
	var rt models.Role = "controller"
	f.role = rt
	return rt

}

func (f *Factory) BuildAppHandlers() *applicationHandler.ApplicationHdl {
	logHandler := f.InitializeLogger()
	configurator := f.InitializeConfigurator()
	clientHandler := f.InitializeHttpClient(logHandler)
	schedulerHandler := f.InitializeScheduler()
	roleHandler := f.InitializeRole()
	importData := f.InitializeImportService()

	notification := f.InitializeNotification()
	protocol := f.InitializeProtocol()
	metrics := f.initializeMetrics(models2.Application)
	auth := f.InitializeAuthentication()
	appSchedulers := application.NewScheduler(clientHandler, configurator, roleHandler, logHandler, notification, protocol, metrics, importData, auth)
	validatorApp := f.InitializeValidator()
	appHandlers := applicationHandler.NewApiExposedHdl(validatorApp, logHandler, clientHandler, schedulerHandler, appSchedulers, configurator)
	return appHandlers
}

func (f *Factory) BuildClusterHandlers() *clusterHandler.ClusterHdl {
	if f.role == models.Worker {
		return nil
	}
	logHandler := f.InitializeLogger()
	configurator := f.InitializeConfigurator()
	clientHandler := f.InitializeHttpClient(logHandler)
	notification := f.InitializeNotification()
	networkInfo := server.NetworkInfo{}
	networkInfo.GetLocalIp()
	metrics := f.initializeMetrics(models2.Cluster)
	schedulerHandler := f.InitializeScheduler()
	clhandler := application2.NewCluster(clientHandler, configurator, logHandler, notification, metrics)
	clusterService := clusterHandler.NewClusterHandler(clhandler, networkInfo, logHandler, clientHandler, schedulerHandler)
	return clusterService
}

func (f *Factory) InitializeImportService() *k8s.ImportProcessApplication {
	if f.importService == nil {
		//k8sRepo := k8sRepo.NewKubernetesRepository("localhost")
		configurator, _ := f.InitializeConfigurator().GetConfig()
		appname := configurator.Server.ApplicationName
		importApp := k8s.NewImportProcessApplication(appname)
		f.importService = importApp
	}

	return f.importService
}

func (f *Factory) InitializeNotification() map[string]ports.Sender {
	if f.notification != nil {
		return f.notification
	}
	configurator := f.InitializeConfigurator()
	config, err := configurator.GetConfig()
	if err != nil {
		panic(err)
	}
	notificationMap := make(map[string]ports.Sender)

	for notificationConfigName, notificationConfig := range config.Notification {
		notification, _ := notificationApp.NotificationSenderFactory(notificationConfigName, notificationConfig)
		notificationMap[notificationConfigName] = notification
		//notification.Send("titolo", "messaggio")
	}
	f.notification = notificationMap
	return f.notification
}

func (f *Factory) InitializeProtocol() map[string]ports2.MonitoringMethod {
	if f.protocol != nil {
		return f.protocol
	}
	configurator := f.InitializeConfigurator()
	log := f.InitializeLogger()
	config, err := configurator.GetConfig()
	if err != nil {
		panic(err)
	}
	protocolMap := make(map[string]ports2.MonitoringMethod)

	for _, protocolConfig := range config.Protocol {
		prot, _ := clientApp.PrtocolFactory(protocolConfig, log)
		protocolMap[protocolConfig] = prot
	}
	f.protocol = protocolMap
	return f.protocol
}

func (f *Factory) initializeMetrics(rt models2.RoleType) *metricsPort.MetricPorts {
	if rt == models2.Application {
		if f.metricApp != nil {
			return f.metricApp
		}
		metricApp := metricsApp.MetricFactory(rt)
		f.metricApp = &metricApp
		return f.metricApp
	} else {
		if f.metricCluster != nil {
			return f.metricCluster
		}
		metricCluster := metricsApp.MetricFactory(rt)
		f.metricCluster = &metricCluster
		return f.metricCluster
	}
}

func (f *Factory) InitializeAuthentication() map[string]portsAuth.Authentication {
	if f.authentication != nil {
		return f.authentication
	}
	configurator := f.InitializeConfigurator()
	config, err := configurator.GetConfig()
	if err != nil {
		panic(err)
	}
	authenticationMap := make(map[string]portsAuth.Authentication)

	for _, authenticationConfig := range config.Application.Authentication {
		authMethod := authenticationApp.AuthFactory(authenticationConfig)
		authenticationMap[authenticationConfig] = authMethod
	}

	f.authentication = authenticationMap
	return f.authentication
}
