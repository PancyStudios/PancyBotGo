package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLocalOnlyMiddlewareRejectsRemoteConnections(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/admin", localOnlyMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "http://localhost/admin", nil)
	request.Host = "localhost"
	request.RemoteAddr = "203.0.113.10:1234"
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("remote request was not rejected: got status %d", recorder.Code)
	}
}

func TestLocalOnlyMiddlewareAllowsLoopback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/admin", localOnlyMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "http://localhost/admin", nil)
	request.Host = "localhost"
	request.RemoteAddr = "127.0.0.1:1234"
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("loopback request was rejected: got status %d", recorder.Code)
	}
}

func TestLocalOnlyMiddlewareRejectsCrossSiteMutation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/admin/api/stop", localOnlyMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodPost, "http://localhost/admin/api/stop", nil)
	request.Host = "localhost"
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("Origin", "https://attacker.example")
	request.Header.Set("Sec-Fetch-Site", "cross-site")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("cross-site mutation was not rejected: got status %d", recorder.Code)
	}
}
