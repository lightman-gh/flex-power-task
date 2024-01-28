package main

import (
	"context"
	"flex/internal/poll"
	service "flex/internal/usecase"

	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ApiHostname string `env:"hostname, default=http://localhost:8081/trades"`
	ApiUsername string `env:"username, default=lightman"`
	ApiPassword string `env:"password, default=password"`
	Location    string `env:"location, default=./examples"`
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	var cfg Config
	if err := envconfig.Process(context.TODO(), &cfg); err != nil {
		logrus.Fatal(err)
	}

	poller := poll.NewCSVPoll(cfg.Location)
	csv := service.CSVServer(cfg.ApiHostname, cfg.ApiUsername, cfg.ApiPassword)

	poller.Poll(context.Background(), csv.Process)
}
