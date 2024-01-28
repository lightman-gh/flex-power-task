package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"flex/internal/foundation/types"
	"fmt"
	"net/http"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/sirupsen/logrus"
)

type CSV struct {
	client   *http.Client
	hostname string
	username string
	password string
}

func CSVServer(hostname, username, password string) *CSV {
	return &CSV{
		client:   http.DefaultClient,
		hostname: hostname,
		username: username,
		password: password,
	}
}

func (c *CSV) Process(ctx context.Context, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can not open file: %s, err: %v", filename, err)
	}

	var trades []*types.Trade

	if err := gocsv.UnmarshalFile(file, &trades); err != nil {
		return fmt.Errorf("can not open file: %s, err: %v", filename, err)
	}

	for _, t := range trades {
		err := c.upload(ctx, t)
		if err != nil {
			logrus.Errorf("can not upload record id: %s, err: %s. Skip", t.ID, err.Error())
		}
	}

	return nil
}

func (c *CSV) upload(ctx context.Context, trade *types.Trade) error {
	bodyBytes, err := json.Marshal(trade)
	if err != nil {
		return fmt.Errorf("can not marshal tarde, err: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.hostname, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("can not create a new http request, err: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)

	rsp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request error, err: %v", err)
	}

	defer func() {
		if err := rsp.Body.Close(); err != nil {
			logrus.Errorf("failed to close insuranceGig response body: %v", err)
		}
	}()

	switch rsp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		logrus.Debugf("uploaded record")
	default:
		return fmt.Errorf("error response received: %d", rsp.StatusCode)
	}

	return nil
}
