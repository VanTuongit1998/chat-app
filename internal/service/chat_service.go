package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"chat-app/internal/model"
	"chat-app/internal/repository"

	"github.com/go-redis/redis/v8"
)

type ChatUsecase struct {
	redisClient       *redis.Client
	messageRepository *repository.MessageRepository
	channelName       string
	queueName         string
}

func NewChatUsecase(redisClient *redis.Client, messageRepository *repository.MessageRepository) *ChatUsecase {
	return &ChatUsecase{
		redisClient:       redisClient,
		messageRepository: messageRepository,
		channelName:       "chat_messages",
		queueName:         "chat_task_queue",
	}
}

func (u *ChatUsecase) PublishMessage(ctx context.Context, msg *model.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := u.redisClient.Publish(ctx, u.channelName, payload).Err(); err != nil {
		return err
	}

	if err := u.redisClient.RPush(ctx, u.queueName, payload).Err(); err != nil {
		return err
	}

	return u.messageRepository.Save(ctx, msg)
}

func (u *ChatUsecase) Conversation(ctx context.Context, userA, userB string, limit int64) ([]*model.Message, error) {
	return u.messageRepository.Conversation(ctx, userA, userB, limit)
}

func (u *ChatUsecase) StartWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("chat worker stopping")
			return
		default:
			u.processNext(ctx)
		}
	}
}

func (u *ChatUsecase) processNext(ctx context.Context) {
	result, err := u.redisClient.BLPop(ctx, 0, u.queueName).Result()
	if err != nil {
		if err == context.Canceled {
			return
		}
		log.Printf("worker BLPop error: %v", err)
		time.Sleep(500 * time.Millisecond)
		return
	}

	if len(result) != 2 {
		return
	}

	var msg model.Message
	if err := json.Unmarshal([]byte(result[1]), &msg); err != nil {
		log.Printf("worker parse fail: %v", err)
		return
	}

	log.Printf("[worker] process task: %s -> %s", msg.Sender, msg.Message)
	if err := u.redisClient.RPush(ctx, "processed_chat_tasks", result[1]).Err(); err != nil {
		log.Printf("worker save processed task error: %v", err)
	}
}

func (u *ChatUsecase) StartQueueMonitor(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("queue monitor stopped")
			return
		case <-ticker.C:
			length, err := u.redisClient.LLen(ctx, u.queueName).Result()
			if err != nil {
				log.Printf("queue monitor error: %v", err)
				continue
			}
			log.Printf("[monitor] queue=%s length=%d", u.queueName, length)
		}
	}
}
