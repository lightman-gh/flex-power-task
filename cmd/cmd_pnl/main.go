package main

import (
	"context"
	"flex/internal/foundation/db"
	"flex/internal/foundation/types/date"
	service "flex/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DatabaseURL string `env:"database_url"`
	TraderId    string `env:"trader_id"`
	DeliveryDay string `env:"delivery_day"`
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

	tradeService := service.NewTradesService(dbClient)

	delivery := date.ISO8601{}
	if err := delivery.Scan(cfg.DeliveryDay); err != nil {
		logrus.Fatal(err)
	}

	records, err := tradeService.ComputePNL(context.Background(), &service.ListQuery{
		DeliveryDate: &delivery,
		TraderID:     &cfg.TraderId,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("  Hour   | Number of Trades | Total BUY [MW] | Total Sell [MW ] | PnL [Eur]")
	logrus.Infof("=========================================================================")

	total := service.Record{}

	for _, r := range records {
		logrus.Infof(" %d - %d |        %2.0d        |        %d       |        %d      | %d", r.Hour-1, r.Hour, r.Count, r.Buy, r.Sell, r.PNL)

		total.Buy += r.Buy
		total.Sell += r.Sell
		total.Count += r.Count
		total.PNL += r.PNL
	}

	logrus.Infof("=========================================================================")
	logrus.Infof(" Total   |        %2.0d        |        %d       |        %d      | %d", total.Count, total.Buy, total.Sell, total.PNL)
}
