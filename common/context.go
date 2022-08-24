package common

import "context"

type ctxKey int

const (
	authTokenKey ctxKey = iota
	authIsAdminKey
)

func SetIsAdmin(isAdmin bool, ctx context.Context) context.Context {
	return context.WithValue(ctx, authIsAdminKey, isAdmin)
}

func GetIsAdmin(ctx context.Context) bool {
	v := ctx.Value(authIsAdminKey)
	isAdmin, ok := v.(bool)
	if !ok {
		return false
	}
	return isAdmin
}

func SetAuthToken(token string, ctx context.Context) context.Context {
	return context.WithValue(ctx, authTokenKey, token)
}

func GetAuthToken(ctx context.Context) string {
	v := ctx.Value(authTokenKey)
	token, ok := v.(string)
	if !ok {
		return ""
	}
	return token
}
