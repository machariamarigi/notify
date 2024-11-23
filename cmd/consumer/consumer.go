package main

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	ConsumerPort       = ":8081"
	ConsumerTopic      = "notifications"
	ConsumerGroup      = "notifications-group"
	KafkaServerAddress = "localhost:9092"
)

// =====================HELPER FUNCTIONS=====================
var ErrNoMessageFound = errors.New("no message found")

func getUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessageFound
	}

	return userID, nil
}


