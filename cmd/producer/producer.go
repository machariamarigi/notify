package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/machariamarigi/notify/pkg/models"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

// =====================HELPER FUNCTIONS=====================
var ErrUserNotFoundInProducer = errors.New("user not found")

func findUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFoundInProducer
}

func getIDFromRequest(formValue string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.PostForm(formValue))
	if err != nil {
		return 0, fmt.Errorf(
			"failed to parse ID from form value %s: %w", formValue, err)
	}

	return id, nil
}

// =====================KAFKA FUNCTIONS=====================
func sendKafkaMessage(producer sarama.SyncProducer, users []models.User, ctx gin.Context, fromID, toID int) error {
	message := ctx.PostForm("message")

	fromUser, err := findUserByID(fromID, users)
	if err != nil {
		return err
	}

	toUser, err := findUserByID(toID, users)
	if err != nil {
		return err
	}

	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key: sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}
