package routes

import (
	"net/http"

	"chat-app/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, roomHandler *handler.RoomHandler, chatHandler *handler.ChatHandler, jwtMiddleware func(http.Handler) http.Handler) http.Handler {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.Use(gin.Logger(), gin.Recovery())

	router.POST("/api/login", gin.WrapF(authHandler.Login))
	router.POST("/api/register", gin.WrapF(authHandler.Register))
	router.GET("/api/users", gin.WrapF(userHandler.Users))
	router.GET("/api/rooms", gin.WrapF(roomHandler.Rooms))
	router.GET("/api/messages", gin.WrapH(jwtMiddleware(http.HandlerFunc(chatHandler.Messages))))

	router.GET("/login", func(c *gin.Context) {
		c.File("web/login.html")
	})
	router.GET("/register", func(c *gin.Context) {
		c.File("web/register.html")
	})
	router.GET("/chat", func(c *gin.Context) {
		c.File("web/chat.html")
	})
	router.GET("/ws", gin.WrapH(jwtMiddleware(http.HandlerFunc(chatHandler.Websocket))))

	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	return router
}
