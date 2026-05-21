# Chat App

Dự án Go mẫu với cấu trúc:

- Realtime chat qua WebSocket
- Worker queue xử lý task từ Redis list
- Redis Pub/Sub broadcast tin nhắn
- JWT authentication cho endpoint WebSocket
- Postgres / Redis cấu hình qua Docker Compose

## Cấu trúc dự án

- `cmd/server/main.go` - entrypoint ứng dụng
- `internal/handler` - HTTP handlers cho auth và chat
- `internal/usecase` - logic nghiệp vụ xác thực và chat
- `internal/repository` - lưu trữ user và message
- `internal/model` - định nghĩa user/message
- `internal/middleware` - logger, JWT auth
- `internal/websocket` - hub/connection WebSocket
- `internal/infrastructure` - kết nối Redis, Postgres, JWT
- `internal/routes` - định tuyến HTTP
- `pkg/response` - helper trả JSON
- `pkg/utils` - tải config từ file env

## Chạy dự án

1. Từ thư mục gốc:
   ```powershell
   go mod tidy
   go run ./cmd/server
   ```

2. Mở trình duyệt vào `http://localhost:8080`

3. Đăng nhập bằng POST `/login` với JSON:
   ```json
   {
     "username": "admin",
     "password": "password"
   }
   ```

4. Gửi `Authorization: Bearer <token>` khi kết nối WebSocket tới `/ws`

## Docker Compose

Khởi động Redis và Postgres:

```powershell
docker compose up -d
```

## Ghi chú

- Redis dùng cho Pub/Sub và queue processing
- Redis/Postgres lấy cấu hình kết nối từ file `.env` ở thư mục gốc
- `internal/usecase/chat_usecase.go` xử lý cả publish và worker queue
