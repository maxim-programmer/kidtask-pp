package child
import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kidtask/internal/auth"
	"kidtask/internal/middleware"
	"kidtask/internal/testutil"
)
func setupChildTest(t *testing.T) (*pgxpool.Pool, *mux.Router, string) {
	db := testutil.SetupTestDB(t)
	parentEmail := testutil.UniqueEmail("parent")
	parentHash, _ := auth.HashPassword("parentpass")
	parentID := testutil.CreateTestParent(t, db, parentEmail, parentHash, "Parent")
	storage := NewStorage(db)
	jwtSecret := "test-secret"
	handler := NewHandler(storage, jwtSecret, middleware.Auth(jwtSecret))
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	return db, r, parentID
}
func TestChildCreateIntegration(t *testing.T) {
	db, router, parentID := setupChildTest(t)
	defer db.Close()
	jwtSecret := "test-secret"
	parentToken, err := auth.GenerateToken(parentID, "parent", "", jwtSecret)
	require.NoError(t, err)
	existingUsername := testutil.UniqueUsername("existing")
	childHash, _ := auth.HashPassword("childpass")
	var existingChildID string
	err = db.QueryRow(context.Background(),
		`INSERT INTO children (parent_id, username, password_hash, name, age_group, balance) 
		 VALUES ($1, $2, $3, $4, $5, 0) RETURNING child_id`,
		parentID, existingUsername, childHash, "Existing Child", "junior",
	).Scan(&existingChildID)
	require.NoError(t, err)
	tests := []struct {
		name        string
		request     map[string]interface{}
		wantStatus  int
		checkAge    bool
		expectedAge string
	}{
		{
			name: "success - junior (age < 11)",
			request: map[string]interface{}{
				"name":     "Junior Child",
				"username": testutil.UniqueUsername("junior"),
				"password": "pass123",
				"birthday": time.Now().AddDate(-10, 0, 0),
			},
			wantStatus:  http.StatusCreated,
			checkAge:    true,
			expectedAge: "junior",
		},
		{
			name: "success - senior (age >= 11)",
			request: map[string]interface{}{
				"name":     "Senior Child",
				"username": testutil.UniqueUsername("senior"),
				"password": "pass123",
				"birthday": time.Now().AddDate(-12, 0, 0),
			},
			wantStatus:  http.StatusCreated,
			checkAge:    true,
			expectedAge: "senior",
		},
		{
			name: "duplicate username in family",
			request: map[string]interface{}{
				"name":     "Duplicate",
				"username": existingUsername,
				"password": "pass123",
				"birthday": time.Now().AddDate(-8, 0, 0),
			},
			wantStatus: http.StatusConflict,
			checkAge:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/children", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+parentToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusCreated && tt.checkAge {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				child := resp["child"].(map[string]interface{})
				assert.Equal(t, tt.expectedAge, child["age_group"])
			}
		})
	}
}
func TestChildLoginIntegration(t *testing.T) {
	db, router, parentID := setupChildTest(t)
	defer db.Close()
	childUsername := testutil.UniqueUsername("testchild")
	childHash, err := auth.HashPassword("childpass")
	require.NoError(t, err)
	var childID string
	err = db.QueryRow(context.Background(),
		`INSERT INTO children (parent_id, username, password_hash, name, age_group, balance) 
		 VALUES ($1, $2, $3, $4, $5, 0) RETURNING child_id`,
		parentID, childUsername, childHash, "Test Child", "junior",
	).Scan(&childID)
	require.NoError(t, err)
	tests := []struct {
		name       string
		request    map[string]string
		wantStatus int
	}{
		{
			name: "correct credentials",
			request: map[string]string{
				"username": childUsername,
				"password": "childpass",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			request: map[string]string{
				"username": childUsername,
				"password": "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/child/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "token")
				assert.Contains(t, resp, "child")
			}
		})
	}
}
