package tests

import (
	"testing"
	"time"

	"github.com/Goga211/go-auth/tests/suite"
	"github.com/Goga211/proto/gen/go/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID          = 1
	appSecret      = "secret_test"
	passDefaultLen = 10
	deltaSec       = 1
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	RegResponse, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, RegResponse.GetUserId())

	LoginResponse, err := st.AuthClient.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: pass,
		AppID:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := LoginResponse.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["appID"].(float64)))
	assert.Equal(t, RegResponse.GetUserId(), int64(claims["uid"].(float64)))

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSec)
}

func TestRegister_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	resp, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetUserId())

	resp, err = st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, resp.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password required",
		},
		{
			name:        "Register with empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email required",
		},
		{
			name:        "Register with both empty",
			email:       "",
			password:    "",
			expectedErr: "password required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Login with empty password",
			email:       email,
			password:    "",
			expectedErr: "password required",
		},
		{
			name:        "Login with empty email",
			email:       "",
			password:    pass,
			expectedErr: "email required",
		},
		{
			name:        "Login with wrong password",
			email:       email,
			password:    randomFakePassword(),
			expectedErr: "invalid credentials",
		},
		{
			name:        "Login with non-existing user",
			email:       gofakeit.Email(),
			password:    pass,
			expectedErr: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &auth.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppID:    appID,
			})
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}
