package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lozhkindm/otus-banners-rotation/internal/app"
	httphandlers "github.com/lozhkindm/otus-banners-rotation/internal/handlers"
	"github.com/lozhkindm/otus-banners-rotation/internal/logger"
	internalhttp "github.com/lozhkindm/otus-banners-rotation/internal/server/http"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", ".env", "Path to configuration file")
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// loading .env
	if err := godotenv.Load(configFile); err != nil {
		log.Fatal(err)
	}

	// populating config
	config := NewConfig()
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}

	// creating logger
	logg, err := logger.New(config.Logger.Level, config.Logger.Development)
	if err != nil {
		logg.Fatal(err.Error())
	}

	// creating storage
	storage, closeFunc, err := NewStorage(ctx, config)
	if err != nil {
		logg.Fatal("failed to create a storage: " + err.Error())
	}
	defer closeFunc(ctx)

	// creating queue
	queue, err := NewQueue(ctx, config)
	if err != nil {
		logg.Fatal("failed to create a queue client: " + err.Error())
	}
	defer queue.Close(ctx)

	// creating application
	application := app.New(logg, storage, queue, config.RabbitMQ.Exchange.Name)

	// creating handlers
	handlers := httphandlers.NewHandlers(application, logg)

	// creating router
	router, err := NewRouter(handlers, config)
	if err != nil {
		logg.Fatal("failed to create a router: " + err.Error())
	}

	// creating http server
	httpServer := internalhttp.NewServer(
		config.Server.Host,
		config.Server.Port,
		config.Server.ReadTimeout,
		config.Server.WriteTimeout,
		config.Server.IdleTimeout,
		logg,
		router,
	)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logg.Fatal("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info(config.App.Name + " is running...")

	if err := httpServer.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
