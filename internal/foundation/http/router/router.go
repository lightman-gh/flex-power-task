package router

import (
	"errors"
	"flex/internal/foundation/http/middleware"
	service "flex/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Router struct {
	engine *gin.Engine

	User   *gin.RouterGroup
	Trades *gin.RouterGroup
}

func NewRouter(userService *service.UsersService) *Router {
	router := &Router{
		engine: gin.Default(),
	}

	router.User = router.engine.Group("users")

	router.Trades = router.engine.Group("trades")
	router.Trades.Use(middleware.WithBasicAuth(userService))

	return router
}

func (router *Router) Run(addr string) error {
	if err := http.ListenAndServe(addr, router.engine); errors.Is(err, http.ErrServerClosed) {
		logrus.Infof("the server closed")
	} else if err != nil {
		logrus.Fatalf("the server closed unpredictably: %s", err.Error())
	}

	return nil
}
