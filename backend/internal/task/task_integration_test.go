package task
import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kidtask/internal/auth"
	"kidtask/internal/middleware"
	"kidtask/internal/testutil"
)
func setupTaskTest(t *testing.T) (*pgxpool.Pool, *mux.Router, string, string) {
	db := testutil.SetupTestDB(t)
	parentEmail := testutil.UniqueEmail("taskparent")
	parentHash, _ := auth.HashPassword("parentpass")
	parentID := testutil.CreateTestParent(t, db, parentEmail, parentHash, "Parent")
	childUsername := testutil.UniqueUsername("taskchild")
	childHash, _ := auth.HashPassword("childpass")
	childID := testutil.CreateTestChild(t, db, parentID, childUsername, childHash, "Child", "junior")
	storage := NewStorage(db)
	jwtSecret := "test-secret"
	handler := NewHandler(storage, middleware.Auth(jwtSecret))
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	return db, r, parentID, childID
}
func TestTaskFullCycle(t *testing.T) {
	db, router, parentID, childID := setupTaskTest(t)
	defer db.Close()
	parentToken, err := auth.GenerateToken(parentID, "parent", "", "test-secret")
	require.NoError(t, err)
	childToken, err := auth.GenerateToken(childID, "child", parentID, "test-secret")
	require.NoError(t, err)
	createReq := map[string]interface{}{
		"child_id": childID,
		"title":    "Test Task",
		"reward":   100,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)
	task := createResp["task"].(map[string]interface{})
	taskID := task["task_id"].(string)
	req = httptest.NewRequest(http.MethodPost, "/api/tasks/"+taskID+"/submit", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	req = httptest.NewRequest(http.MethodPost, "/api/tasks/"+taskID+"/approve", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var balance int
	err = db.QueryRow(context.Background(), "SELECT balance FROM children WHERE child_id = $1", childID).Scan(&balance)
	require.NoError(t, err)
	assert.Equal(t, 100, balance)
}
func TestTaskRejectAndResubmit(t *testing.T) {
	db, router, parentID, childID := setupTaskTest(t)
	defer db.Close()
	parentToken, err := auth.GenerateToken(parentID, "parent", "", "test-secret")
	require.NoError(t, err)
	childToken, err := auth.GenerateToken(childID, "child", parentID, "test-secret")
	require.NoError(t, err)
	createReq := map[string]interface{}{
		"child_id": childID,
		"title":    "Test Task",
		"reward":   100,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)
	task := createResp["task"].(map[string]interface{})
	taskID := task["task_id"].(string)
	req = httptest.NewRequest(http.MethodPost, "/api/tasks/"+taskID+"/submit", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	comment := "Need to do better"
	rejectBody, _ := json.Marshal(map[string]string{"comment": comment})
	req = httptest.NewRequest(http.MethodPost, "/api/tasks/"+taskID+"/reject", bytes.NewBuffer(rejectBody))
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	req = httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var getResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &getResp)
	require.NoError(t, err)
	taskResp := getResp["task"].(map[string]interface{})
	assert.Equal(t, StatusNeedsRework, taskResp["status"])
	assert.Equal(t, comment, taskResp["rejection_comment"])
	req = httptest.NewRequest(http.MethodPost, "/api/tasks/"+taskID+"/submit", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	req = httptest.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &getResp)
	require.NoError(t, err)
	taskResp = getResp["task"].(map[string]interface{})
	assert.Equal(t, StatusPendingReview, taskResp["status"])
}
