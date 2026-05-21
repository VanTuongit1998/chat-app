package routes

import (
	"net/http"

	"chat-app/internal/handler"
)

func NewRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, roomHandler *handler.RoomHandler, chatHandler *handler.ChatHandler, jwtMiddleware func(http.Handler) http.Handler, loggerMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/users", userHandler.Users)
	mux.HandleFunc("/api/rooms", roomHandler.Rooms)
	mux.Handle("/api/messages", jwtMiddleware(http.HandlerFunc(chatHandler.Messages)))
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "web/login.html")
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	// serve register page
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "web/register.html")
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "web/chat.html")
	})
	mux.Handle("/ws", jwtMiddleware(http.HandlerFunc(chatHandler.Websocket)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})
	return loggerMiddleware(mux)
}
