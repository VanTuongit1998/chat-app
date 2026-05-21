package repository

import (
	"context"
	"database/sql"

	"chat-app/internal/model"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) (*RoomRepository, error) {
	repo := &RoomRepository{db: db}
	if err := repo.Migrate(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *RoomRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS rooms (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

func (r *RoomRepository) FindAll(ctx context.Context) ([]*model.Room, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, created_at FROM rooms ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]*model.Room, 0)
	for rows.Next() {
		room := &model.Room{}
		if err := rows.Scan(&room.ID, &room.Name, &room.CreatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}
