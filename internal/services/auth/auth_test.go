package auth

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Goga211/go-auth/internal/domain/model"
	"github.com/Goga211/go-auth/internal/services/auth/mocks"
	"github.com/Goga211/go-auth/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_UserNotFound(t *testing.T) {
	userProvider := mocks.NewMockUserProvider(t)

	userProvider.EXPECT().
		User(mock.Anything, mock.Anything).
		Return(model.User{}, storage.ErrUserNotFound)

	a := New(slog.Default(), nil, nil, userProvider, time.Hour)

	_, err := a.Login(context.Background(), "ghost@mail.com", "pass", 1)

	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_HappyPath(t *testing.T) {
	const (
		email = "user@mail.com"
		pass  = "secret"
		appID = 1
	)

	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)

	userProvider := mocks.NewMockUserProvider(t)
	userProvider.EXPECT().
		User(mock.Anything, email).
		Return(model.User{ID: 1, Email: email, PassHash: hash}, nil)

	appProvider := mocks.NewMockAppProvider(t)
	appProvider.EXPECT().
		App(mock.Anything, appID).
		Return(model.App{ID: appID, Secret: "test-secret"}, nil)

	a := New(slog.Default(), nil, appProvider, userProvider, time.Hour)

	token, err := a.Login(context.Background(), email, pass, appID)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}
