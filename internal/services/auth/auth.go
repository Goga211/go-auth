package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Goga211/go-auth/internal/domain/model"
	"github.com/Goga211/go-auth/internal/lib/jwt"
	"github.com/Goga211/go-auth/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (model.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (model.App, error)
}

func New(log *slog.Logger, userSaver UserSaver, appProvider AppProvider, userProvider UserProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const operation = "auth.login"
	log := a.log.With(
		slog.String("operation", operation),
		slog.String("email", email))

	log.Info("login in progress")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", operation, ErrInvalidCredentials)
		}
		log.Warn("user not found")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid password")
		return "", fmt.Errorf("%s: %w", operation, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Warn("app not found")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("user logged in")
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const operation = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("operation", operation),
		slog.String("email", email))

	log.Info("Registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash password", slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	uid, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	log.Info("User registered successfully")
	return uid, nil
}
