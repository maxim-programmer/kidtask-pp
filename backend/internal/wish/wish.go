package wish

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"kidtask/internal/middleware"
	"kidtask/internal/respond"
)

const (
	StatusAwaitingPrice = "awaiting_price"
	StatusAvailable     = "available"
	StatusPurchased     = "purchased"
	StatusDelivered     = "delivered"
)

type Wish struct {
	WishID      string  `json:"wish_id"`
	ChildID     string  `json:"child_id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Price       *int    `json:"price,omitempty"`
	Status      string  `json:"status"`
}

func (w *Wish) Validate() error {
	if strings.TrimSpace(w.Title) == "" {
		return errors.New("title is required")
	}
	return nil
}

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(ctx context.Context, w *Wish) error {
	return s.db.QueryRow(ctx,
		`INSERT INTO wishes (child_id, title, description)
		 VALUES ($1, $2, $3) RETURNING wish_id, status`,
		w.ChildID, w.Title, w.Description,
	).Scan(&w.WishID, &w.Status)
}

func (s *Storage) GetByID(ctx context.Context, id string) (*Wish, error) {
	w := &Wish{}
	err := s.db.QueryRow(ctx,
		`SELECT wish_id, child_id, title, description, price, status FROM wishes WHERE wish_id = $1`, id,
	).Scan(&w.WishID, &w.ChildID, &w.Title, &w.Description, &w.Price, &w.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return w, err
}

func (s *Storage) ListByChild(ctx context.Context, childID, status string) ([]*Wish, error) {
	query := `SELECT wish_id, child_id, title, description, price, status FROM wishes WHERE child_id = $1`
	args := []any{childID}

	if status != "" {
		args = append(args, status)
		query += ` AND status = $2`
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishes []*Wish
	for rows.Next() {
		w := &Wish{}
		if err := rows.Scan(&w.WishID, &w.ChildID, &w.Title, &w.Description, &w.Price, &w.Status); err != nil {
			return nil, err
		}
		wishes = append(wishes, w)
	}
	return wishes, nil
}

func (s *Storage) Update(ctx context.Context, w *Wish) error {
	_, err := s.db.Exec(ctx,
		`UPDATE wishes SET title=$1, description=$2 WHERE wish_id=$3`,
		w.Title, w.Description, w.WishID,
	)
	return err
}

func (s *Storage) SetPrice(ctx context.Context, id string, price int) error {
	_, err := s.db.Exec(ctx,
		`UPDATE wishes SET price=$1, status='available' WHERE wish_id=$2`,
		price, id,
	)
	return err
}

func (s *Storage) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := s.db.Exec(ctx, `UPDATE wishes SET status=$1 WHERE wish_id=$2`, status, id)
	return err
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM wishes WHERE wish_id = $1`, id)
	return err
}

func (s *Storage) DeductBalance(ctx context.Context, childID string, amount int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE children SET balance = balance - $1 WHERE child_id = $2 AND balance >= $1`,
		amount, childID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("insufficient coins")
	}
	return nil
}

func (s *Storage) GetChildParentID(ctx context.Context, childID string) (string, error) {
	var parentID string
	err := s.db.QueryRow(ctx, `SELECT parent_id FROM children WHERE child_id = $1`, childID).Scan(&parentID)
	return parentID, err
}

type Handler struct {
	storage *Storage
	mw      func(http.Handler) http.Handler
}

func NewHandler(storage *Storage, mw func(http.Handler) http.Handler) *Handler {
	return &Handler{storage: storage, mw: mw}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.Handle("/api/children/{childId}/wishes", h.mw(http.HandlerFunc(h.List))).Methods(http.MethodGet)
	r.Handle("/api/children/{childId}/wishes", h.mw(http.HandlerFunc(
		middleware.RequireRole("child", h.Create)))).Methods(http.MethodPost)
	r.Handle("/api/children/{childId}/wishes/{id}", h.mw(http.HandlerFunc(h.Update))).Methods(http.MethodPatch)
	r.Handle("/api/children/{childId}/wishes/{id}", h.mw(http.HandlerFunc(h.Delete))).Methods(http.MethodDelete)
	r.Handle("/api/children/{childId}/wishes/{id}/purchase", h.mw(http.HandlerFunc(
		middleware.RequireRole("child", h.Purchase)))).Methods(http.MethodPost)
	r.Handle("/api/children/{childId}/wishes/{id}/deliver", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Deliver)))).Methods(http.MethodPatch)
}

func (h *Handler) checkAccess(w http.ResponseWriter, r *http.Request, childID string) bool {
	claims := middleware.GetClaims(r)
	if claims.Role == "parent" {
		parentID, err := h.storage.GetChildParentID(r.Context(), childID)
		if err != nil || parentID != claims.UserID {
			respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
			return false
		}
	}
	if claims.Role == "child" && claims.UserID != childID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return false
	}
	return true
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	childID := mux.Vars(r)["childId"]
	if !h.checkAccess(w, r, childID) {
		return
	}

	wishes, err := h.storage.ListByChild(r.Context(), childID, r.URL.Query().Get("status"))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if wishes == nil {
		wishes = []*Wish{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"wishes": wishes})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	childID := mux.Vars(r)["childId"]

	if claims.UserID != childID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}

	wsh := &Wish{ChildID: childID, Title: req.Title, Description: req.Description}
	if err := wsh.Validate(); err != nil {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.storage.Create(r.Context(), wsh); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]any{"wish": wsh})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	childID := mux.Vars(r)["childId"]
	id := mux.Vars(r)["id"]

	if !h.checkAccess(w, r, childID) {
		return
	}

	wsh, err := h.storage.GetByID(r.Context(), id)
	if err != nil || wsh == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "wish not found")
		return
	}
	if wsh.ChildID != childID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}
	if wsh.Status == StatusPurchased || wsh.Status == StatusDelivered {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "cannot edit purchased wish")
		return
	}

	if claims.Role == "parent" {
		var req struct {
			Price *int `json:"price"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
			return
		}
		if req.Price == nil || *req.Price <= 0 {
			respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "price must be greater than 0")
			return
		}
		if err := h.storage.SetPrice(r.Context(), id, *req.Price); err != nil {
			respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
			return
		}
		wsh.Price = req.Price
		wsh.Status = StatusAvailable
	} else {
		var req struct {
			Title       *string `json:"title"`
			Description *string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
			return
		}
		if req.Title != nil {
			wsh.Title = *req.Title
		}
		if req.Description != nil {
			wsh.Description = req.Description
		}
		if err := h.storage.Update(r.Context(), wsh); err != nil {
			respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
			return
		}
	}

	respond.JSON(w, http.StatusOK, map[string]any{"wish": wsh})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	childID := mux.Vars(r)["childId"]
	id := mux.Vars(r)["id"]

	if !h.checkAccess(w, r, childID) {
		return
	}

	wsh, err := h.storage.GetByID(r.Context(), id)
	if err != nil || wsh == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "wish not found")
		return
	}
	if wsh.ChildID != childID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}
	if wsh.Status == StatusPurchased || wsh.Status == StatusDelivered {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "cannot delete purchased wish")
		return
	}

	if err := h.storage.Delete(r.Context(), id); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"message": "wish deleted"})
}

func (h *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	childID := mux.Vars(r)["childId"]
	id := mux.Vars(r)["id"]

	if claims.UserID != childID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	wsh, err := h.storage.GetByID(r.Context(), id)
	if err != nil || wsh == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "wish not found")
		return
	}
	if wsh.Status != StatusAvailable {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "wish is not available")
		return
	}
	if wsh.Price == nil {
		respond.Error(w, http.StatusBadRequest, "NO_PRICE", "price not set")
		return
	}

	if err := h.storage.DeductBalance(r.Context(), childID, *wsh.Price); err != nil {
		respond.Error(w, http.StatusBadRequest, "INSUFFICIENT_COINS", "not enough coins")
		return
	}

	if err := h.storage.UpdateStatus(r.Context(), id, StatusPurchased); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	wsh.Status = StatusPurchased
	respond.JSON(w, http.StatusOK, map[string]any{"wish": wsh})
}

func (h *Handler) Deliver(w http.ResponseWriter, r *http.Request) {
	childID := mux.Vars(r)["childId"]
	id := mux.Vars(r)["id"]

	if !h.checkAccess(w, r, childID) {
		return
	}

	wsh, err := h.storage.GetByID(r.Context(), id)
	if err != nil || wsh == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "wish not found")
		return
	}
	if wsh.Status != StatusPurchased {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "wish is not purchased")
		return
	}

	if err := h.storage.UpdateStatus(r.Context(), id, StatusDelivered); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	wsh.Status = StatusDelivered
	respond.JSON(w, http.StatusOK, map[string]any{"wish": wsh})
}