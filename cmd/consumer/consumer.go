package main

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/machariamarigi/notify/pkg/models"
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

// ===================NOTIFICATION STORAGE====================
type UserNotifications map[string][]models.Notification

type NotificationStore struct {
	data UserNotifications
	mu   sync.RWMutex
}

func (ns *NotificationStore) Add(userID string, notification models.Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) []models.Notification {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return ns.data[userID]
}
