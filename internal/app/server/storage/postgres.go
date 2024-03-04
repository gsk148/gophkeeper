package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	DeleteSession    = `DELETE FROM sessions WHERE cid = $1`
	GetSession       = "SELECT token FROM sessions WHERE cid = $1"
	StoreSession     = "INSERT INTO sessions(cid, token) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING token"
	AddUser          = "INSERT INTO users(name, password) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id"
	DeleteUser       = "DELETE FROM users WHERE id = $1"
	GetUserByID      = "SELECT * FROM users WHERE id = $1"
	GetUserByName    = "SELECT * FROM users WHERE name = $1"
	DeleteData       = "DELETE FROM storage WHERE uid = $1 AND id = $2"
	GetAllDataByType = "SELECT * FROM storage WHERE uid = $1 AND type = $2"
	GetDataByID      = "SELECT * FROM storage WHERE uid = $1 AND id = $2"
	StoreData        = `
		INSERT INTO storage(uid, data, type) VALUES($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id
	`
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(url string) (*DBStorage, error) {
	if url == "" {
		return &DBStorage{}, ErrDBMissingURL
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		return &DBStorage{}, err
	}

	if err = initDB(db); err != nil {
		return &DBStorage{}, err
	}
	return &DBStorage{db: db}, nil
}

func initDB(db *sql.DB) error {
	const op = "storage.Init"
	if err := goose.SetDialect("postgres"); err != nil {
		log.Errorf("Init DB: failed while goose set dialect, %s", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		log.Errorf("Init DB: failed while goose up, %s", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *DBStorage) DeleteSession(ctx context.Context, cid string) error {
	_, err := r.db.ExecContext(ctx, DeleteSession, cid)
	return err
}

func (r *DBStorage) GetSession(ctx context.Context, cid string) (string, error) {
	var token string
	err := r.db.QueryRowContext(ctx, GetSession, cid).Scan(&token)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return token, err
}

func (r *DBStorage) StoreSession(ctx context.Context, cid, token string) error {
	_, err := r.db.ExecContext(ctx, StoreSession, cid, token)
	return err
}

func (r *DBStorage) AddUser(ctx context.Context, u User) (User, error) {
	var id string
	if err := r.db.QueryRowContext(ctx, AddUser, u.Name, u.Password).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	u.ID = id
	return u, nil
}

func (r *DBStorage) DeleteUser(ctx context.Context, uid string) error {
	_, err := r.db.ExecContext(ctx, DeleteUser, uid)
	return err
}

func (r *DBStorage) GetUserByID(ctx context.Context, uid string) (User, error) {
	return r.getUser(ctx, GetUserByID, uid)
}

func (r *DBStorage) GetUserByName(ctx context.Context, name string) (User, error) {
	return r.getUser(ctx, GetUserByName, name)
}

func (r *DBStorage) DeleteData(ctx context.Context, uid, id string) error {
	_, err := r.db.ExecContext(ctx, DeleteData, uid, id)
	return err
}

func (r *DBStorage) GetAllDataByType(ctx context.Context, uid string, t Type) ([]SecureData, error) {
	rows, err := r.db.QueryContext(ctx, GetAllDataByType, uid, t)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer closeRows(rows)

	var data []SecureData
	for rows.Next() {
		var piece SecureData
		if err = rows.Scan(&piece.ID, &piece.UID, &piece.Data, &piece.Type); err != nil {
			return nil, err
		}
		data = append(data, piece)
	}
	return data, nil
}

func (r *DBStorage) GetDataByID(ctx context.Context, uid, id string) (SecureData, error) {
	var data SecureData
	err := r.db.QueryRowContext(ctx, GetDataByID, uid, id).Scan(&data.ID, &data.UID, &data.Data, &data.Type)
	return data, err
}

func (r *DBStorage) StoreData(ctx context.Context, data SecureData) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, StoreData, data.UID, data.Data, data.Type).Scan(&id)
	return id, err
}

func (r *DBStorage) getUser(ctx context.Context, query string, args ...any) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, nil
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Error(err)
	}
}
