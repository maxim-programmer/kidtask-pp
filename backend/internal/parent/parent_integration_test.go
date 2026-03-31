//go:build integration
// +build integration

package parent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kidtask/internal/auth"
	"kidtask/internal/testutil"
)

func TestParentRegisterIntegration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	storage := NewStorage(db)
	jwtSecret := "test-secret"
	handler := NewHandler(storage, jwtSecret)

	// Создаём уникальный email для существующего родителя
	existingEmail := testutil.UniqueEmail("existing")
	hash, err := auth.HashPassword("testpass")
	require.NoError(t, err)
	testutil.CreateTestParent(t, db, existingEmail, hash, "Existing Parent")

	tests := []struct {
		name       string
		request    map[string]string
		wantStatus int
	}{
		{
			name: "success - new parent",
			request: map[string]string{
				"email":    testutil.UniqueEmail("new"),
				"password": "password123",
				"name":     "New Parent",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate email",
			request: map[string]string{
				"email":    existingEmail,
				"password": "password123",
				"name":     "Duplicate Parent",
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "invalid email - no @",
			request: map[string]string{
				"email":    "invalid-email",
				"password": "password123",
				"name":     "Invalid Parent",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "token")
				assert.Contains(t, resp, "user")
			}
		})
	}
}

func TestParentLoginIntegration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	storage := NewStorage(db)
	jwtSecret := "test-secret"
	handler := NewHandler(storage, jwtSecret)

	// Создаём тестового родителя с уникальным email
	email := testutil.UniqueEmail("login")
	hash, err := auth.HashPassword("correctpass")
	require.NoError(t, err)
	testutil.CreateTestParent(t, db, email, hash, "Login Parent")

	tests := []struct {
		name       string
		request    map[string]string
		wantStatus int
	}{
		{
			name: "correct credentials",
			request: map[string]string{
				"email":    email,
				"password": "correctpass",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			request: map[string]string{
				"email":    email,
				"password": "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-existent email",
			request: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "somepass",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "token")
				assert.Contains(t, resp, "user")
			}
		})
	}
}
