package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Goga211/go-auth/internal/domain/model"
	"github.com/Goga211/go-auth/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const operation = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const operation = "sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, email, passHash)

	if err != nil {
		if sqliteErr, ok := errors.AsType[sqlite3.Error](err); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", operation, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (model.User, error) {
	const operation = "sqlite.Get.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", operation, err)
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	var user model.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", operation, storage.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", operation, err)
	}

	return user, nil
}

func (s *Storage) App(ctx context.Context, appID int) (model.App, error) {
	const operation = "sqlite.Get.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return model.App{}, fmt.Errorf("%s: %w", operation, err)
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, appID)

	var app model.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.App{}, fmt.Errorf("%s: %w", operation, storage.ErrAppNotFound)
		}
		return model.App{}, fmt.Errorf("%s: %w", operation, err)
	}

	return app, nil
}
