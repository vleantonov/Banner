package middleware_test

import (
	"banner/internal/domain"
	"banner/internal/handler/http/v1/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testUserToken  = "user"
	testAdminToken = "admin"
)

func setupMockEngine() *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.CheckToken(
			testUserToken,
			testAdminToken,
		),
	)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, c.GetString(domain.UserStatusHeader))
	})

	return r
}

func TestAuthPassion(t *testing.T) {
	testMap := []struct {
		name         string
		token        string
		expectStatus int
		expectUser   string
	}{
		{
			name:         "Auth with admin token",
			token:        testAdminToken,
			expectStatus: http.StatusOK,
			expectUser:   domain.Admin,
		},
		{
			name:         "Auth with user token",
			token:        testUserToken,
			expectStatus: http.StatusOK,
			expectUser:   domain.User,
		},
		{
			name:         "Auth with invalid token",
			token:        "InvalidToken",
			expectStatus: http.StatusForbidden,
			expectUser:   "",
		},
		{
			name:         "Auth with empty token",
			token:        "",
			expectStatus: http.StatusUnauthorized,
			expectUser:   "",
		},
	}

	r := setupMockEngine()

	for _, test := range testMap {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("token", test.token)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, test.expectStatus, res.Code)
		if test.expectUser != "" {
			assert.Equal(t, test.expectUser, res.Body.String())
		}
	}
}

func TestWithoutToken(t *testing.T) {
	r := setupMockEngine()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
