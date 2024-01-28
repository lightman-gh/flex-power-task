package trades

import (
	"flex/internal/dbschema/model"
	"flex/internal/foundation/context"
	"flex/internal/foundation/http/response"
	"flex/internal/foundation/types"
	"flex/internal/foundation/types/date"
	service "flex/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *TradeHandler) registerTradesRoutes() {
	handler.router.Trades.POST("", handler.handlePostTrade)
	handler.router.Trades.GET("", handler.handleGetTrades)
	handler.router.Trades.GET("pnl", handler.handleGetPnl)
}

func (handler *TradeHandler) handlePostTrade(ctx *gin.Context) {
	var trade types.Trade

	if err := ctx.ShouldBindJSON(&trade); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))

		return
	}

	usr := context.GetUser(ctx)

	t, err := handler.trade.Create(ctx, usr, &model.Trade{
		ID:            trade.ID,
		Price:         trade.Price,
		Quantity:      trade.Quantity,
		Direction:     trade.Direction,
		DeliveryDay:   trade.DeliveryDay,
		DeliveryHour:  trade.DeliveryHour,
		TraderID:      trade.TraderId,
		ExecutionTime: trade.ExecutionTime,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error("can not create a new trade, %s", err.Error()))

		return
	}

	ctx.JSON(http.StatusCreated, response.OK(t))
}

//type ListQuery struct {
//	DeliveryDate *date.ISO8601 `form:"delivery_day,omitempty"`
//	TraderID     *string       `form:"trader_id,omitempty"`
//}

func (handler *TradeHandler) handleGetTrades(ctx *gin.Context) {
	// For some reason this does not work with custom date. Interesting ....
	// Fallback to simple query and parse

	//var query ListQuery
	//
	//if err := ctx.BindQuery(&query); err != nil {
	//	ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error("can not read query, err: %s", err.Error()))
	//
	//	return
	//}

	serviceQuery := &service.ListQuery{}
	if tID := ctx.Query("trader_id"); tID != "" {
		serviceQuery.TraderID = &tID
	}

	if dDate := ctx.Query("delivery_day"); dDate != "" {
		pDate := date.ISO8601{}
		if err := pDate.Scan(dDate); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error("incorrect date format: %s", err.Error()))

			return
		}

		serviceQuery.DeliveryDate = &pDate
	}

	trades, err := handler.trade.List(ctx, serviceQuery)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))

		return
	}

	ctx.JSON(http.StatusOK, trades)
}

func (handler *TradeHandler) handleGetPnl(ctx *gin.Context) {
	serviceQuery := &service.ListQuery{}
	if tID := ctx.Query("trader_id"); tID != "" {
		serviceQuery.TraderID = &tID
	}

	if dDate := ctx.Query("delivery_day"); dDate != "" {
		pDate := date.ISO8601{}
		if err := pDate.Scan(dDate); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error("incorrect date format: %s", err.Error()))

			return
		}

		serviceQuery.DeliveryDate = &pDate
	}

	trades, err := handler.trade.ComputePNL(ctx, serviceQuery)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))

		return
	}

	ctx.JSON(http.StatusOK, trades)
}
