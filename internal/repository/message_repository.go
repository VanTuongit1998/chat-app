package repository

import (
	"context"
	"database/sql"

	"chat-app/internal/model"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) (*MessageRepository, error) {
	repo := &MessageRepository{db: db}
	if err := repo.Migrate(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *MessageRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS messages (
			id BIGSERIAL PRIMARY KEY,
			sender TEXT NOT NULL,
			recipient TEXT NOT NULL,
			body TEXT NOT NULL,
			sent_at TIMESTAMPTZ NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_messages_conversation
			ON messages (sender, recipient, sent_at);
	`)
	return err
}

func (r *MessageRepository) Save(ctx context.Context, message *model.Message) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO messages (sender, recipient, body, sent_at) VALUES ($1, $2, $3, $4)`,
		message.Sender,
		message.To,
		message.Message,
		message.SentAt,
	)
	return err
}

func (r *MessageRepository) Latest(ctx context.Context, limit int64) ([]*model.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT sender, recipient, body, sent_at
		FROM messages
		ORDER BY sent_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reversed := make([]*model.Message, 0)
	for rows.Next() {
		msg := &model.Message{}
		if err := rows.Scan(&msg.Sender, &msg.To, &msg.Message, &msg.SentAt); err != nil {
			return nil, err
		}
		reversed = append(reversed, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	messages := make([]*model.Message, 0, len(reversed))
	for i := len(reversed) - 1; i >= 0; i-- {
		messages = append(messages, reversed[i])
	}

	return messages, nil
}

func (r *MessageRepository) Conversation(ctx context.Context, userA, userB string, limit int64) ([]*model.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT sender, recipient, body, sent_at
		FROM messages
		WHERE (sender = $1 AND recipient = $2)
			OR (sender = $2 AND recipient = $1)
		ORDER BY sent_at DESC
		LIMIT $3
	`, userA, userB, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reversed := make([]*model.Message, 0)
	for rows.Next() {
		msg := &model.Message{}
		if err := rows.Scan(&msg.Sender, &msg.To, &msg.Message, &msg.SentAt); err != nil {
			return nil, err
		}
		reversed = append(reversed, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]*model.Message, 0, len(reversed))
	for i := len(reversed) - 1; i >= 0; i-- {
		result = append(result, reversed[i])
	}

	return result, nil
}
