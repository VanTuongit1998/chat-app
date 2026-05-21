package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"chat-app/configs"
	"chat-app/internal/handler"
	"chat-app/internal/middleware"
	"chat-app/internal/repository"
	"chat-app/internal/routes"
	"chat-app/internal/service"
	"chat-app/internal/utils"
	"chat-app/internal/websocket"

	_ "github.com/lib/pq"
)

func main() {
	_ = configs.LoadEnv(".env")

	redisAddr := configs.GetEnv("REDIS_ADDR", "")
	if redisAddr == "" {
		redisAddr = fmt.Sprintf("%s:%s", requiredEnv("REDIS_HOST"), requiredEnv("REDIS_PORT"))
	}

	postgresDSN := configs.GetEnv("POSTGRES_DSN", "")
	if postgresDSN == "" {
		postgresDSN = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			requiredEnv("POSTGRES_USER"),
			requiredEnv("POSTGRES_PASSWORD"),
			requiredEnv("POSTGRES_HOST"),
			requiredEnv("POSTGRES_PORT"),
			requiredEnv("POSTGRES_DB"),
			requiredEnv("POSTGRES_SSLMODE"),
		)
	}

	jwtSecret := configs.GetEnv("JWT_SECRET", "secret123")
	appPort := configs.GetEnv("APP_PORT", "8080")

	db, err := configs.NewPostgresDB(postgresDSN)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	defer db.Close()

	redisClient := configs.NewClient(redisAddr)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("failed to initialize user repository: %v", err)
	}
	roomRepo, err := repository.NewRoomRepository(db)
	if err != nil {
		log.Fatalf("failed to initialize room repository: %v", err)
	}
	messageRepo, err := repository.NewMessageRepository(db)
	if err != nil {
		log.Fatalf("failed to initialize message repository: %v", err)
	}
	jwtService := utils.NewJwtService(jwtSecret, time.Hour*24)

	authUsecase := service.NewAuthUsecase(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	roomService := service.NewRoomService(roomRepo)
	chatUsecase := service.NewChatUsecase(redisClient, messageRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)
	hub := websocket.NewHub(chatUsecase, redisClient)
	chatHandler := handler.NewChatHandler(hub, chatUsecase)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)
	go chatUsecase.StartWorker(ctx)
	go chatUsecase.StartQueueMonitor(ctx)

	jwtMiddleware := middleware.JwtAuthentication(jwtService)
	loggerMiddleware := middleware.Logger
	router := routes.NewRouter(authHandler, userHandler, roomHandler, chatHandler, jwtMiddleware, loggerMiddleware)

	listenAddr := fmt.Sprintf(":%s", appPort)
	log.Printf("server listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func requiredEnv(key string) string {
	value := configs.GetEnv(key, "")
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}
