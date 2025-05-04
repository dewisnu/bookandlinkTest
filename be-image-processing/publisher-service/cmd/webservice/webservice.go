package webservice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"os"
	"publisher-service/cmd/router"
	"publisher-service/internal/config"
	"publisher-service/internal/repository"
	"publisher-service/internal/service"
)

const logTagStartWebservice = "[Start]"

func Start(conf *config.Config) {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug, // Debug, Info, Warn, Error
		AddSource: false,           // Show file & line number
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger) // Set global logger

	err := config.InitializeDirectories()
	if err != nil {
		log.Fatalf("Failed to initialize directories: %v", err)
	}

	db, err := config.InitDB(&config.InitDatabaseParams{
		Conf: &conf.DatabaseConfig,
	})

	if err != nil {
		slog.Error(fmt.Sprintf("%s initializing  db: %+v", logTagStartWebservice, err))
	}

	rabbitmq, err := config.InitRabbitMQ(&config.InitRabbitMQParams{
		Conf: &conf.RabbitMQConfig,
	})

	if err != nil {
		slog.Error(fmt.Sprintf("Failed to initialize RabbitMQ: %v", err))
	}
	defer rabbitmq.Close()

	gin.SetMode(conf.GinMode)
	gn := gin.New()
	repo := repository.NewRepository(&repository.NewRepositoryParams{
		Database: db,
	})

	serv := service.NewService(&service.NewServiceParams{
		Repository: repo,
		RabbitMQ:   rabbitmq,
	})

	router.Init(&router.InitRouterParams{
		Service: serv,
		Gn:      gn,
		Conf:    conf,
	})

	slog.Info(fmt.Sprintf("%s Publisher service starting on port: %s", logTagStartWebservice, conf.ServicePort))

	if err = gn.Run(conf.ServicePort); err != nil {
		slog.Error(fmt.Sprintf("%s starting service, cause: %+v", logTagStartWebservice, err))
	}
}
