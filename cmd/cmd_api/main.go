package main

import (
	"context"
	"flex/internal/api/trades"
	"flex/internal/foundation/db"
	"flex/internal/foundation/http/router"
	service "flex/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DatabaseURL string `env:"database_url"`
	APIHost     string `env:"api_host, default=0.0.0.0:8085"`
}

func main() {
	gin.SetMode(gin.DebugMode)
	logrus.SetLevel(logrus.TraceLevel)

	var cfg Config

	if err := envconfig.Process(context.TODO(), &cfg); err != nil {
		logrus.Fatal(err)
	}

	_, dbClient := db.NewPostgresClient(cfg.DatabaseURL)
	defer func() { _ = dbClient.Close() }()

	usrService := service.NewUsersService(dbClient)
	tradeService := service.NewTradesService(dbClient)

	r := router.NewRouter(usrService)

	handler := trades.NewTradeHandler(r, usrService, tradeService)
	handler.RegisterRouters()

	r.Run(cfg.APIHost)
}
