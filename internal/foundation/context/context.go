package context

import (
	"context"
	"flex/internal/dbschema/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetUser(ctx context.Context, u *model.User) {
	gContext, ok := ctx.(*gin.Context)
	if !ok {
		logrus.Fatal("can not convert context to gin context")
	}

	gContext.Set(KeyUser, u)
}

func GetUser(ctx context.Context) *model.User {
	gContext, ok := ctx.(*gin.Context)
	if !ok {
		logrus.Fatal("can not convert context to gin context")
	}

	ctxUser, exist := gContext.Get(KeyUser)
	if !exist {
		return nil
	}

	u, ok := ctxUser.(*model.User)
	if !ok {
		logrus.Fatal("can not convert user")
	}

	return u
}
