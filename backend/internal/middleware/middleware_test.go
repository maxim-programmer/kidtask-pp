//go:build integration
// +build integration

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"kidtask/internal/auth"
)

func TestAuthMiddlewareNoToken(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	called := false
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		claims := GetClaims(r)
		assert.NotNil(t, claims)
		assert.Equal(t, "test-user", claims.UserID)
		assert.Equal(t, "parent", claims.Role)
		w.WriteHeader(http.StatusOK)
	}))

	token, _ := auth.GenerateToken("test-user", "parent", "", secret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRoleChildOnParentRoute(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	childToken, _ := auth.GenerateToken("child-1", "child", "parent-1", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/parent-route", RequireRole("parent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/parent-route", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRoleParentOnChildRoute(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	parentToken, _ := auth.GenerateToken("parent-1", "parent", "", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/child-route", RequireRole("child", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/child-route", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRoleCorrectRole(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	called := false
	parentToken, _ := auth.GenerateToken("parent-1", "parent", "", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/parent-route", RequireRole("parent", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/parent-route", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRoleChildCorrectRole(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	called := false
	childToken, _ := auth.GenerateToken("child-1", "child", "parent-1", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/child-route", RequireRole("child", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/child-route", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}
