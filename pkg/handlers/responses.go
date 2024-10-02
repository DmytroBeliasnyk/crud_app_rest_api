package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Message string `json:"message"`
}

func newErrResponse(ctx *gin.Context, status int, message string) {
	log.Printf("ERROR: %s\n", message)

	ctx.AbortWithStatusJSON(status, errResponse{message})
}
