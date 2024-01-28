package usecase

import (
	"context"
	"flex/internal/dbschema/model"
	"flex/internal/dbschema/model/predicate"
	"flex/internal/dbschema/model/trade"
	"flex/internal/foundation/types/date"
	"fmt"
	"slices"
)

type TradesService struct {
	client *model.Client
}

func NewTradesService(client *model.Client) *TradesService {
	return &TradesService{
		client: client,
	}
}

func (t *TradesService) Create(ctx context.Context, usr *model.User, trade *model.Trade) (*model.Trade, error) {
	return t.client.Trade.Create().
		SetID(trade.ID).
		SetPrice(trade.Price).
		SetQuantity(trade.Quantity).
		SetDirection(trade.Direction).
		SetDeliveryDay(trade.DeliveryDay).
		SetDeliveryHour(trade.DeliveryHour).
		SetTraderID(trade.TraderID).
		SetExecutionTime(trade.ExecutionTime).
		SetCreator(usr).
		Save(ctx)
}

type ListQuery struct {
	DeliveryDate *date.ISO8601
	TraderID     *string
}

func (t *TradesService) List(ctx context.Context, query *ListQuery) ([]*model.Trade, error) {
	var pd []predicate.Trade

	if query.DeliveryDate != nil {
		pd = append(pd, trade.DeliveryDayEQ(*query.DeliveryDate))
	}

	if query.TraderID != nil {
		pd = append(pd, trade.TraderID(*query.TraderID))
	}

	trades, err := t.client.Trade.Query().
		Where(pd...).
		All(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return []*model.Trade{}, nil
		}

		return nil, fmt.Errorf("can not select tardes, err: %v", err)
	}

	return trades, nil
}

type Record struct {
	Hour  int32 `json:"hour"`
	Count int   `json:"count"`
	Buy   int64 `json:"buy"`
	Sell  int64 `json:"sell"`
	PNL   int64 `json:"pnl"`
}

func (t *TradesService) ComputePNL(ctx context.Context, query *ListQuery) ([]*Record, error) {
	lst, err := t.List(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("can not get list of trades, err: %v", err)
	}

	hours := map[int32][]*model.Trade{}

	// 1. collect all trades by time
	for _, t := range lst {
		arr, ok := hours[t.DeliveryHour]
		if !ok {
			arr = make([]*model.Trade, 0, 5)
		}

		arr = append(arr, t)
		hours[t.DeliveryHour] = arr
	}

	// 2. compute amount by an hour
	records := make([]*Record, 0, len(hours))

	for key, arr := range hours {
		rc := &Record{
			Hour:  key,
			Count: len(arr),
			Buy:   0,
			Sell:  0,
			PNL:   0,
		}

		for _, hour := range arr {
			switch hour.Direction {
			case "sell":
				rc.Sell += int64(hour.Quantity * hour.Price)
			case "buy":
				rc.Buy -= int64(hour.Quantity * hour.Price)
			}
		}

		rc.PNL = rc.Sell + rc.Buy
		records = append(records, rc)
	}

	slices.SortFunc(records, func(a, b *Record) int {
		if a.Hour > b.Hour {
			return 1
		}

		return -1
	})

	return records, nil
}
