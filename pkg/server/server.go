package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	_ "greye/docs"
	apiExposedApp "greye/internal/application/domain/ports"
	apiExposedCl "greye/internal/cluster/domain/ports"
	configPort "greye/pkg/config/domain/ports"
	"greye/pkg/role/domain/models"
	"os"
)

type Server struct {
	app          *fiber.App
	networkInfo  NetworkInfo
	application  apiExposedApp.ApiExposed
	cluster      apiExposedCl.ApiExposed
	configurator configPort.ConfigApplication
	role         models.Role
}

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)
)

func init() {
	// Registra le metriche con il registratore Prometheus
	prometheus.MustRegister(httpRequestsTotal)
}

func NewServer(applicationServer apiExposedApp.ApiExposed, clusterServer apiExposedCl.ApiExposed, configurator configPort.ConfigApplication, role models.Role) *Server {

	return &Server{application: applicationServer, cluster: clusterServer, configurator: configurator, role: role}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func (s *Server) Run(port string) error {
	app := fiber.New()
	app.Use(cors.New())
	//app.Use(logger.New())
	app.Use(logger.New(logger.Config{
		Format:     `{"time":"${time}", "ip":"${ip}", "port":"${port}", "status":${status}, "method":"${method}", "path":"${path}", "latency":"${latency}"}` + "\n",
		Output:     os.Stdout,
		TimeFormat: "2006-01-02T15:04:05Z07:00", // ISO 8601 format for time
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/metrics", func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(c.Context())
		return nil
	})
	v1Application := app.Group("/api/v1/application")

	// User Endpoints
	//v1Application.Post("", s.application.AddApplicationBySvc)

	v1Application.Put("/monitor", s.application.MonitoringApplication)
	v1Application.Get("/monitor", s.application.GetApplicationMonitored)
	v1Application.Get("/monitor/pod", s.application.GetApplicationMonitoredByPod)

	v1Application.Delete("/monitor/:service", s.application.UnscheduleApplication)
	//v1Application.Post("/check", s.application.Check)

	if s.role == models.Controller {
		v1Cluster := app.Group("/api/v1/cluster")

		// User Endpoints

		v1Cluster.Get("/status", s.cluster.Status)
		v1Cluster.Put("/status", s.cluster.UpdateStatus)
		v1Cluster.Put("/suspend", s.cluster.UpdateSingleStatus)
		v1Cluster.Delete("", s.cluster.Remove)

	}
	s.app = app
	err := app.Listen(":" + port)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	err := s.app.Shutdown()
	if err != nil {
		return err
	}
	return nil
}
