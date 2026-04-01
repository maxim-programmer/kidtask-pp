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

func TestParentCannotSeeOtherChild(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	parentToken, _ := auth.GenerateToken("parent-1", "parent", "", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/children/{id}", RequireRole("parent", func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		vars := mux.Vars(r)
		if vars["id"] != claims.UserID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/children/child-2", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChildCannotSeeOtherTask(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	childToken, _ := auth.GenerateToken("child-1", "child", "parent-1", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/tasks/{id}", RequireRole("child", func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		vars := mux.Vars(r)
		if vars["id"] != claims.UserID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/tasks/task-2", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChildCannotPurchaseOtherWish(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	childToken, _ := auth.GenerateToken("child-1", "child", "parent-1", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/children/{childId}/wishes/{wishId}/purchase", RequireRole("child", func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		vars := mux.Vars(r)
		if vars["childId"] != claims.UserID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/children/child-2/wishes/wish-1/purchase", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestParentCannotApproveOtherTask(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	parentToken, _ := auth.GenerateToken("parent-1", "parent", "", secret)

	router := mux.NewRouter()
	router.Use(mw)

	router.Handle("/tasks/{id}/approve", RequireRole("parent", func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		vars := mux.Vars(r)
		if vars["id"] != claims.UserID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/tasks/task-2/approve", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestParentCanAccessOwnChild(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	parentToken, _ := auth.GenerateToken("parent-1", "parent", "", secret)

	router := mux.NewRouter()
	router.Use(mw)

	called := false
	router.Handle("/children/{id}", RequireRole("parent", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/children/child-1", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChildCanAccessOwnTask(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	childToken, _ := auth.GenerateToken("child-1", "child", "parent-1", secret)

	router := mux.NewRouter()
	router.Use(mw)

	called := false
	router.Handle("/tasks/{id}", RequireRole("child", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/tasks/task-1", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChildCanPurchaseOwnWish(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	childToken, _ := auth.GenerateToken("child-1", "child", "parent-1", secret)

	router := mux.NewRouter()
	router.Use(mw)

	called := false
	router.Handle("/children/{childId}/wishes/{wishId}/purchase", RequireRole("child", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/children/child-1/wishes/wish-1/purchase", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestParentCanApproveOwnTask(t *testing.T) {
	secret := "test-secret"
	mw := Auth(secret)

	parentToken, _ := auth.GenerateToken("parent-1", "parent", "", secret)

	router := mux.NewRouter()
	router.Use(mw)

	called := false
	router.Handle("/tasks/{id}/approve", RequireRole("parent", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/tasks/task-1/approve", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}
