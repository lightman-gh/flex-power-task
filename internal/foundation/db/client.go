package db

import (
	"context"
	"flex/internal/dbschema/model"

	"entgo.io/ent/dialect"
	entSQL "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewPostgresClient(backendURL string) (dialect.Driver, *model.Client) {
	drv, err := entSQL.Open(dialect.Postgres, backendURL)
	if err != nil {
		logrus.Fatal(err)
	}

	opts := []model.Option{
		model.Log(logrus.Debug),
		model.Driver(drv),
	}

	client := model.NewClient(opts...)

	if err := client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}

	return drv, client
}
