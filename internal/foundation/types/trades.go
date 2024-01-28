package types

import (
	"flex/internal/foundation/types/date"
	"time"
)

type Trade struct {
	ID            string       `json:"id" binding:"required" csv:"id"`
	Price         int32        `json:"price" binding:"required,gte=1" csv:"price"`
	Quantity      int32        `json:"quantity" binding:"required,gte=1" csv:"quantity"`
	Direction     string       `json:"direction" binding:"required" csv:"direction"`
	DeliveryDay   date.ISO8601 `json:"delivery_day" binding:"required" csv:"delivery_day"`
	DeliveryHour  int32        `json:"delivery_hour" binding:"required,lte=23,gte=0" csv:"delivery_hour"`
	TraderId      string       `json:"trader_id" binding:"required" csv:"trader_id"`
	ExecutionTime time.Time    `json:"execution_time" binding:"required" csv:"execution_time"`
}
