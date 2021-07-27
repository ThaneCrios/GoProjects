package main

import (
	"flag"
	"gitlab.com/faemproject/backend/delivery/delivery.front/proto"
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"gitlab.com/faemproject/backend/delivery/delivery.front/repository"

	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"

	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/core/shared/os"
	"gitlab.com/faemproject/backend/core/shared/prometheus"
	"gitlab.com/faemproject/backend/core/shared/rabbit"
	"gitlab.com/faemproject/backend/core/shared/store"
	"gitlab.com/faemproject/backend/core/shared/web"
	"gitlab.com/faemproject/backend/core/shared/web/middleware"
	"gitlab.com/faemproject/backend/delivery/delivery.front/broker/publisher"
	"gitlab.com/faemproject/backend/delivery/delivery.front/broker/subscriber"
	"gitlab.com/faemproject/backend/delivery/delivery.front/config"
	"gitlab.com/faemproject/backend/delivery/delivery.front/handler"
	"gitlab.com/faemproject/backend/delivery/delivery.front/server"
)

const (
	defaultConfigPath     = "deployment/config/devlocal/deliveryfront.toml"
	maxRequestsAllowed    = 1000
	serverShutdownTimeout = 30 * time.Second
	brokerShutdownTimeout = 30 * time.Second
)

func main() {
	// Parse flags
	configPath := flag.String("config", defaultConfigPath, "configuration file path")
	flag.Parse()

	cfg, err := config.Parse(*configPath)
	if err != nil {
		log.Fatalf("failed to parse the config file: %v", err)
	}

	cfg.Print() // just for debugging

	if err := logs.SetLogLevel(cfg.Application.LogLevel); err != nil {
		log.Fatalf("Failed to set log level: %v", err)
	}
	if err := logs.SetLogFormat(cfg.Application.LogFormat); err != nil {
		log.Fatalf("Failed to set log format: %v", err)
	}
	logger := logs.Eloger

	// Connect to the db and remember to close it
	db, err := store.Connect(&pg.Options{
		Addr:     store.Addr(cfg.Database.Host, cfg.Database.Port),
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Db,
	})
	if err != nil {
		logger.Fatalf("failed to create a db instance: %v", err)
	}
	defer db.Close()

	// Connect to the broker and remember to close it
	rmq := &rabbit.Rabbit{
		Credits: rabbit.ConnCredits{
			URL:  cfg.Broker.UserURL,
			User: cfg.Broker.UserCredits,
		},
	}
	if err = rmq.Init(cfg.Broker.ExchangePrefix, cfg.Broker.ExchangePostfix); err != nil {
		logger.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.CloseRabbit()

	// Create a publisher
	pub := publisher.Publisher{
		Rabbit:  rmq,
		Encoder: &rabbit.JsonEncoder{},
	}
	if err = pub.Init(); err != nil {
		logger.Fatalf("failed to init the publisher: %v", err)
	}
	defer pub.Wait(brokerShutdownTimeout)

	// Create a service object
	hdlr := handler.Handler{
		DB:     repository.Pg{Db: db},
		Pub:    &pub,
		Config: cfg,
	}

	// Create a subscriber
	sub := subscriber.Subscriber{
		Rabbit:  rmq,
		Encoder: &rabbit.JsonEncoder{},
		Handler: &hdlr,
	}
	if err = sub.Init(); err != nil {
		logger.Fatalf("failed to start the subscriber: %v", err)
	}
	defer sub.Wait(brokerShutdownTimeout)

	// Create a rest gateway and handle http requests
	router := web.NewRouter(
		loggerOption(logger),
		prometheusmetric,
		throttler,
		cors,
	)
	rest := server.Rest{
		Router:  router,
		Handler: &hdlr,
	}
	rest.Route()

	// Start an http server and remember to shut it down
	go web.Start(router, cfg.Application.Port)
	defer web.Stop(router, serverShutdownTimeout)

	// Wait for program exit
	<-os.NotifyAboutExit()
}

func loggerOption(logger *logrus.Logger) web.Option {
	return func(e *echo.Echo) {
		e.Logger = &middleware.Logger{Logger: logger} // replace the original echo.Logger with the logrus one
		// Log the requests
		e.Use(middleware.LoggerWithSkipper(
			func(c echo.Context) bool {
				return strings.Contains(c.Request().RequestURI, "/api/v2/locations")
			}),
		)
	}
}

func cors(e *echo.Echo) {
	e.Use(echomiddleware.CORS())
}

func throttler(e *echo.Echo) {
	e.Use(middleware.Throttle(maxRequestsAllowed))
}

func prometheusmetric(e *echo.Echo) {
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
}

func init() {
	proto.InitCourierStatesVariables()
	proto.InitCourierTypesVariables()

	proto.InitOrderStatesVariables()
	proto.InitOrderTypesVariables()

	proto.InitTaskTypesVariables()

	proto.InitPaymentTypesVariables()
	proto.InitPaymentStatusesVariables()

	proto.LocalisationInit()
}
