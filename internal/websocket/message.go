package websocket

type IncomingMessage struct {
	To      string `json:"to"`
	Message string `json:"message"`
}
