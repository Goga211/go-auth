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

func TestLogin_FailCases(t *testing.T) {
	const (
		email = "user@mail.com"
		pass  = "secret"
		appID = 1
	)

	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)

	tests := []struct {
		name     string
		password string
		setup    func(up *mocks.MockUserProvider, ap *mocks.MockAppProvider)
		wantErr  error
	}{
		{
			name:     "user not found",
			password: pass,
			setup: func(up *mocks.MockUserProvider, ap *mocks.MockAppProvider) {
				up.EXPECT().User(mock.Anything, email).Return(model.User{}, storage.ErrUserNotFound)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name:     "wrong password",
			password: "wrong-pass",
			setup: func(up *mocks.MockUserProvider, ap *mocks.MockAppProvider) {
				up.EXPECT().User(mock.Anything, email).Return(model.User{ID: 1, Email: email, PassHash: hash}, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name:     "app not found",
			password: pass,
			setup: func(up *mocks.MockUserProvider, ap *mocks.MockAppProvider) {
				up.EXPECT().User(mock.Anything, email).Return(model.User{ID: 1, Email: email, PassHash: hash}, nil)
				ap.EXPECT().App(mock.Anything, appID).Return(model.App{}, storage.ErrAppNotFound)
			},
			wantErr: storage.ErrAppNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			up := mocks.NewMockUserProvider(t)
			ap := mocks.NewMockAppProvider(t)
			tt.setup(up, ap)

			a := New(slog.Default(), nil, ap, up, time.Hour)

			_, err := a.Login(context.Background(), email, tt.password, appID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
