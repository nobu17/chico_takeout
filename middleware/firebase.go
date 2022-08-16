package middleware

import (
	"context"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type firebaseApp struct {
	*firebase.App
}

type AuthService interface {
	VerifyIDToken(ctx context.Context, idToken string) (*AuthData, error)
}

type AuthData struct {
	IsAuthorized bool
	IsAdmin      bool
}

func NewFirebaseApp() (*firebaseApp, error) {
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_CREDENTIALS_JSON"))))
	if err != nil {
		return nil, err
	}
	return &firebaseApp{app}, nil
}

func (app *firebaseApp) VerifyIDToken(ctx context.Context, idToken string) (*AuthData, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return &AuthData{IsAdmin: false, IsAuthorized: false}, err
	}
	result := AuthData{IsAdmin: false, IsAuthorized: true}
	if role, ok := token.Claims["role"]; ok {
		if role.(string) == "Admin" {
			result.IsAdmin = true
		}
	}
	return &result, nil
}
