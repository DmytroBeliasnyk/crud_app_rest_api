package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Message string `json:"message"`
}

func newErrResponse(ctx *gin.Context, status int, err string) {
	logrus.WithFields(logrus.Fields{
		"uri":    ctx.Request.RequestURI,
		"method": ctx.Request.Method,
	}).Error(err)

	ctx.AbortWithStatusJSON(status, errResponse{err})
}
