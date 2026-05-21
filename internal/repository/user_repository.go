package repository

import (
	"context"
	"database/sql"
	"errors"

	"chat-app/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	repo := &UserRepository{db: db}
	if err := repo.Migrate(context.Background()); err != nil {
		return nil, err
	}
	if err := repo.SeedDefaults(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *UserRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

func (r *UserRepository) SeedDefaults(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username, password)
		VALUES
			('admin', 'password'),
			('guest', 'guest')
		ON CONFLICT (username) DO NOTHING;
	`)
	return err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRowContext(
		context.Background(),
		`SELECT id, username, password FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindAll() []*model.User {
	rows, err := r.db.QueryContext(context.Background(), `
		SELECT id, username, password
		FROM users
		ORDER BY username ASC
	`)
	if err != nil {
		return []*model.User{}
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err == nil {
			users = append(users, user)
		}
	}
	return users
}

func (r *UserRepository) Create(user *model.User) error {
	err := r.db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
		user.Username,
		user.Password,
	).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
