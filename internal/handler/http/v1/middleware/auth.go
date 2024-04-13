package middleware

import (
	"banner/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	tokenHeader = "token"
	emptyToken  = ""
)

func CheckToken(userToken, adminToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(tokenHeader)
		if token == emptyToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		switch token {
		case userToken:
			c.Set(domain.UserStatusHeader, domain.User)
		case adminToken:
			c.Set(domain.UserStatusHeader, domain.Admin)
		default:
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
