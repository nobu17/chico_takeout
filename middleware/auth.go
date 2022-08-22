package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authTokenKey   = "authToken"
	authIsAdminKey = "isAdmin"
)

func SetAuthInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getAuthTokenFromHeader(c.Request)
		if err != nil {
			c.Next()
			return
		}
		setAuthToken(c, token)
		c.Next()
	}
}

func CheckAuthInfo(auth AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getAuthToken(c)
		result, err := auth.VerifyIDToken(c.Request.Context(), token)
		if err != nil {
			handleUnAuth(c)
			return
		}
		// set auth role
		setIsAdmin(c, result.IsAdmin)
		c.Next()
	}
}

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin := getIsAdmin(c)
		if !isAdmin {
			handleForbidden(c)
			return
		}
		c.Next()
	}
}

func setAuthToken(c *gin.Context, token string) {
	c.Set(authTokenKey, token)
}

func getAuthToken(c *gin.Context) string {
	return c.GetString(authTokenKey)
}

func setIsAdmin(c *gin.Context, isAdmin bool) {
	c.Set(authIsAdminKey, isAdmin)
}

func getIsAdmin(c *gin.Context) bool {
	result, ok := c.Get(authIsAdminKey)
	if !ok {
		return false
	}
	return result.(bool)
}

func handleUnAuth(c *gin.Context) {
	c.JSON(401, gin.H{"message": "invalid auth"})
	c.Abort()
}

func handleForbidden(c *gin.Context) {
	c.JSON(403, gin.H{"message": "invalid auth right"})
	c.Abort()
}

func getAuthTokenFromHeader(r *http.Request) (string, error) {
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		return "", errors.New("no Bearer headers")
	}
	headers := strings.Split(r.Header.Get("Authorization"), " ")

	if len(headers) < 2 {
		return "", errors.New("incorrect header format")
	}
	return headers[1], nil
}