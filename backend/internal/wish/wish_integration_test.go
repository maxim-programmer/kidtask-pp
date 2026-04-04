package wish
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
func setupWishTest(t *testing.T) (*pgxpool.Pool, *mux.Router, string, string) {
	db := testutil.SetupTestDB(t)
	parentEmail := testutil.UniqueEmail("wishparent")
	parentHash, _ := auth.HashPassword("parentpass")
	parentID := testutil.CreateTestParent(t, db, parentEmail, parentHash, "Parent")
	childUsername := testutil.UniqueUsername("wishchild")
	childHash, _ := auth.HashPassword("childpass")
	var childID string
	err := db.QueryRow(context.Background(),
		`INSERT INTO children (parent_id, username, password_hash, name, age_group, balance) 
		 VALUES ($1, $2, $3, $4, $5, 500) RETURNING child_id`,
		parentID, childUsername, childHash, "Child", "junior",
	).Scan(&childID)
	require.NoError(t, err)
	storage := NewStorage(db)
	jwtSecret := "test-secret"
	handler := NewHandler(storage, middleware.Auth(jwtSecret))
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	return db, r, parentID, childID
}
func TestWishFullCycle(t *testing.T) {
	db, router, parentID, childID := setupWishTest(t)
	defer db.Close()
	jwtSecret := "test-secret"
	parentToken, err := auth.GenerateToken(parentID, "parent", "", jwtSecret)
	require.NoError(t, err)
	childToken, err := auth.GenerateToken(childID, "child", parentID, jwtSecret)
	require.NoError(t, err)
	createReq := map[string]string{"title": "New Bike"}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/children/"+childID+"/wishes", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)
	wish := createResp["wish"].(map[string]interface{})
	wishID := wish["wish_id"].(string)
	priceReq := map[string]int{"price": 200}
	body, _ = json.Marshal(priceReq)
	req = httptest.NewRequest(http.MethodPatch, "/api/children/"+childID+"/wishes/"+wishID, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	req = httptest.NewRequest(http.MethodPost, "/api/children/"+childID+"/wishes/"+wishID+"/purchase", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	req = httptest.NewRequest(http.MethodPatch, "/api/children/"+childID+"/wishes/"+wishID+"/deliver", nil)
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	req = httptest.NewRequest(http.MethodGet, "/api/children/"+childID+"/wishes", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var listResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResp)
	require.NoError(t, err)
	wishes := listResp["wishes"].([]interface{})
	found := false
	for _, wItem := range wishes {
		wishItem := wItem.(map[string]interface{})
		if wishItem["wish_id"] == wishID {
			assert.Equal(t, StatusDelivered, wishItem["status"])
			found = true
			break
		}
	}
	assert.True(t, found, "Wish not found in list")
}
func TestWishPurchaseInsufficientBalance(t *testing.T) {
	db, router, parentID, childID := setupWishTest(t)
	defer db.Close()
	jwtSecret := "test-secret"
	parentToken, err := auth.GenerateToken(parentID, "parent", "", jwtSecret)
	require.NoError(t, err)
	childToken, err := auth.GenerateToken(childID, "child", parentID, jwtSecret)
	require.NoError(t, err)
	createReq := map[string]string{"title": "Expensive Bike"}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/children/"+childID+"/wishes", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+childToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)
	wish := createResp["wish"].(map[string]interface{})
	wishID := wish["wish_id"].(string)
	priceReq := map[string]int{"price": 1000}
	body, _ = json.Marshal(priceReq)
	req = httptest.NewRequest(http.MethodPatch, "/api/children/"+childID+"/wishes/"+wishID, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+parentToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	req = httptest.NewRequest(http.MethodPost, "/api/children/"+childID+"/wishes/"+wishID+"/purchase", nil)
	req.Header.Set("Authorization", "Bearer "+childToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var balance int
	err = db.QueryRow(context.Background(), "SELECT balance FROM children WHERE child_id = $1", childID).Scan(&balance)
	require.NoError(t, err)
	assert.Equal(t, 500, balance)
}
