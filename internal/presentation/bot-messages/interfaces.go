package presentation

import (
	"github.com/go-redis/redis/v8"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type ChatSetup struct {
	client        *whatsmeow.Client
	evt           *events.Message
	redisClient   *redis.Client
	currentChatId string
}
