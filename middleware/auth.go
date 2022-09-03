package middleware

import (
	"errors"
	"net/http"
	"strings"

	"chico/takeout/common"

	"github.com/gin-gonic/gin"
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
		// set auth role and userId
		setIsAdmin(c, result.IsAdmin)
		setUserId(c, result.UserId)
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
	ctx := common.SetAuthToken(token, c.Request.Context())
	c.Request = c.Request.WithContext(ctx)
}

func getAuthToken(c *gin.Context) string {
	return common.GetAuthToken(c.Request.Context())
}

func setIsAdmin(c *gin.Context, isAdmin bool) {
	ctx := common.SetIsAdmin(isAdmin, c.Request.Context())
	c.Request = c.Request.WithContext(ctx)
}

func getIsAdmin(c *gin.Context) bool {
	return common.GetIsAdmin(c.Request.Context())
}

func setUserId(c *gin.Context, userId string) {
	ctx := common.SetUserId(userId, c.Request.Context())
	c.Request = c.Request.WithContext(ctx)
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
