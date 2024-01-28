package trades

import (
	"flex/internal/foundation/http/router"
	service "flex/internal/usecase"
)

type TradeHandler struct {
	router *router.Router

	user  *service.UsersService
	trade *service.TradesService
}

func NewTradeHandler(
	router *router.Router,
	user *service.UsersService,
	trade *service.TradesService,
) *TradeHandler {
	return &TradeHandler{
		router: router,
		user:   user,
		trade:  trade,
	}
}

func (handler *TradeHandler) RegisterRouters() {
	handler.registerUserRoutes()
	handler.registerTradesRoutes()
}
