package user

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib" // SQL driver
)

type DBRepo struct {
	db *sql.DB
}

const (
	AddUser       = "INSERT INTO users(name, password) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id"
	DeleteUser    = "DELETE FROM users WHERE id = $1"
	GetUserByID   = "SELECT * FROM users WHERE id = $1"
	GetUserByName = "SELECT * FROM users WHERE name = $1"
)

func NewDBRepo(url string) (*DBRepo, error) {
	if url == "" {
		return &DBRepo{}, ErrDBMissingURL
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return &DBRepo{}, err
	}

	return &DBRepo{db: db}, err
}

func (r *DBRepo) AddUser(ctx context.Context, user User) (User, error) {
	if user.Name == "" || user.Password == "" {
		return User{}, ErrCredMissing
	}

	var id string
	if err := r.db.QueryRowContext(ctx, AddUser, user.Name, user.Password).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	user.ID = id
	return user, nil
}

func (r *DBRepo) DeleteUser(ctx context.Context, uid string) error {
	res, err := r.db.ExecContext(ctx, DeleteUser, uid)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *DBRepo) GetUserByID(ctx context.Context, uid string) (User, error) {
	if uid == "" {
		return User{}, ErrNotFound
	}
	return r.getUser(ctx, GetUserByID, uid)
}

func (r *DBRepo) GetUserByName(ctx context.Context, name string) (User, error) {
	if name == "" {
		return User{}, ErrNotFound
	}
	return r.getUser(ctx, GetUserByName, name)
}

func (r *DBRepo) getUser(ctx context.Context, query string, args ...any) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, nil
}
