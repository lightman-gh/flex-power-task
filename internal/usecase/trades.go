package usecase

import (
	"context"
	"flex/internal/dbschema/model"
	"flex/internal/dbschema/model/predicate"
	"flex/internal/dbschema/model/trade"
	"flex/internal/foundation/types/date"
	"fmt"
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
