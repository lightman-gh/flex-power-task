package middleware

import (
	"flex/internal/dbschema/model"
	"flex/internal/foundation/context"
	"flex/internal/foundation/http/response"
	"flex/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WithBasicAuth(user *usecase.UsersService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, password, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("the api must authorised"))

			return
		}

		u, err := user.GetByPassword(ctx, &model.User{
			Username: username,
			Password: password,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Error("can not authorise user, err: %s", err.Error()))
		}

		context.SetUser(ctx, u)
	}
}
