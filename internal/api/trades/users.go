package trades

import (
	"flex/internal/dbschema/model"
	"flex/internal/foundation/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *TradeHandler) registerUserRoutes() {
	handler.router.User.POST("", handler.handlePostUser)
}

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *TradeHandler) handlePostUser(ctx *gin.Context) {
	var usr User

	if err := ctx.ShouldBindJSON(&usr); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error(err.Error()))

		return
	}

	err := handler.user.Create(ctx, &model.User{
		Username: usr.Username,
		Password: usr.Password,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error("can not create a new user, %s", err.Error()))
	}

	ctx.JSON(http.StatusOK, response.OK("Created"))
}
