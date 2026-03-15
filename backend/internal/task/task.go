package task

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
	StatusActive        = "active"
	StatusPendingReview = "pending_review"
	StatusNeedsRework   = "needs_rework"
	StatusCompleted     = "completed"
)

type Task struct {
	TaskID           string  `json:"task_id"`
	ParentID         string  `json:"parent_id"`
	ChildID          string  `json:"child_id"`
	Title            string  `json:"title"`
	Description      *string `json:"description,omitempty"`
	Reward           int     `json:"reward"`
	Status           string  `json:"status"`
	RejectionComment *string `json:"rejection_comment,omitempty"`
}

func (t *Task) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title is required")
	}
	if t.Reward <= 0 {
		return errors.New("reward must be greater than 0")
	}
	return nil
}

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(ctx context.Context, t *Task) error {
	return s.db.QueryRow(ctx,
		`INSERT INTO tasks (parent_id, child_id, title, description, reward)
		 VALUES ($1, $2, $3, $4, $5) RETURNING task_id, status`,
		t.ParentID, t.ChildID, t.Title, t.Description, t.Reward,
	).Scan(&t.TaskID, &t.Status)
}

func (s *Storage) GetByID(ctx context.Context, id string) (*Task, error) {
	t := &Task{}
	err := s.db.QueryRow(ctx,
		`SELECT task_id, parent_id, child_id, title, description, reward, status, rejection_comment
		 FROM tasks WHERE task_id = $1`, id,
	).Scan(&t.TaskID, &t.ParentID, &t.ChildID, &t.Title, &t.Description, &t.Reward, &t.Status, &t.RejectionComment)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (s *Storage) List(ctx context.Context, parentID, childID, status string) ([]*Task, error) {
	query := `SELECT task_id, parent_id, child_id, title, description, reward, status, rejection_comment
	          FROM tasks WHERE parent_id = $1`
	args := []any{parentID}

	if childID != "" {
		args = append(args, childID)
		query += ` AND child_id = $` + itoa(len(args))
	}
	if status != "" {
		args = append(args, status)
		query += ` AND status = $` + itoa(len(args))
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.TaskID, &t.ParentID, &t.ChildID, &t.Title, &t.Description, &t.Reward, &t.Status, &t.RejectionComment); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Storage) ListByChild(ctx context.Context, childID, status string) ([]*Task, error) {
	query := `SELECT task_id, parent_id, child_id, title, description, reward, status, rejection_comment
	          FROM tasks WHERE child_id = $1`
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

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.TaskID, &t.ParentID, &t.ChildID, &t.Title, &t.Description, &t.Reward, &t.Status, &t.RejectionComment); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Storage) UpdateStatus(ctx context.Context, id, status string, comment *string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE tasks SET status=$1, rejection_comment=$2 WHERE task_id=$3`,
		status, comment, id,
	)
	return err
}

func (s *Storage) Update(ctx context.Context, t *Task) error {
	_, err := s.db.Exec(ctx,
		`UPDATE tasks SET title=$1, description=$2, reward=$3 WHERE task_id=$4`,
		t.Title, t.Description, t.Reward, t.TaskID,
	)
	return err
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM tasks WHERE task_id = $1`, id)
	return err
}

func (s *Storage) AddBalance(ctx context.Context, childID string, amount int) error {
	_, err := s.db.Exec(ctx,
		`UPDATE children SET balance = balance + $1 WHERE child_id = $2`,
		amount, childID,
	)
	return err
}

func itoa(n int) string {
	return string(rune('0' + n))
}

type Handler struct {
	storage *Storage
	mw      func(http.Handler) http.Handler
}

func NewHandler(storage *Storage, mw func(http.Handler) http.Handler) *Handler {
	return &Handler{storage: storage, mw: mw}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.Handle("/api/tasks", h.mw(http.HandlerFunc(h.List))).Methods(http.MethodGet)
	r.Handle("/api/tasks", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Create)))).Methods(http.MethodPost)
	r.Handle("/api/tasks/{id}", h.mw(http.HandlerFunc(h.Get))).Methods(http.MethodGet)
	r.Handle("/api/tasks/{id}", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Update)))).Methods(http.MethodPatch)
	r.Handle("/api/tasks/{id}", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Delete)))).Methods(http.MethodDelete)
	r.Handle("/api/tasks/{id}/submit", h.mw(http.HandlerFunc(
		middleware.RequireRole("child", h.Submit)))).Methods(http.MethodPost)
	r.Handle("/api/tasks/{id}/approve", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Approve)))).Methods(http.MethodPost)
	r.Handle("/api/tasks/{id}/reject", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Reject)))).Methods(http.MethodPost)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	var tasks []*Task
	var err error

	if claims.Role == "parent" {
		tasks, err = h.storage.List(r.Context(), claims.UserID,
			r.URL.Query().Get("child_id"), r.URL.Query().Get("status"))
	} else {
		tasks, err = h.storage.ListByChild(r.Context(), claims.UserID, r.URL.Query().Get("status"))
	}

	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	var req struct {
		ChildID     string  `json:"child_id"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Reward      int     `json:"reward"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	if req.ChildID == "" {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "child_id is required")
		return
	}

	t := &Task{ParentID: claims.UserID, ChildID: req.ChildID, Title: req.Title, Description: req.Description, Reward: req.Reward}
	if err := t.Validate(); err != nil {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.storage.Create(r.Context(), t); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]any{"task": t})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	t, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if t == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}
	if claims.Role == "parent" && t.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}
	if claims.Role == "child" && t.ChildID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"task": t})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	t, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if t == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}
	if t.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}
	if t.Status != StatusActive {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "can only edit active tasks")
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Reward      *int    `json:"reward"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}

	if req.Title != nil {
		t.Title = *req.Title
	}
	if req.Description != nil {
		t.Description = req.Description
	}
	if req.Reward != nil {
		t.Reward = *req.Reward
	}

	if err := t.Validate(); err != nil {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.storage.Update(r.Context(), t); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"task": t})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	t, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if t == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}
	if t.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	if err := h.storage.Delete(r.Context(), id); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"message": "task deleted"})
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	t, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if t == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}
	if t.ChildID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}
	if t.Status != StatusActive && t.Status != StatusNeedsRework {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "task cannot be submitted")
		return
	}

	if err := h.storage.UpdateStatus(r.Context(), id, StatusPendingReview, nil); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	t.Status = StatusPendingReview
	respond.JSON(w, http.StatusOK, map[string]any{"task": t})
}

func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	t, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if t == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}
	if t.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}
	if t.Status != StatusPendingReview {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "task is not pending review")
		return
	}

	if err := h.storage.UpdateStatus(r.Context(), id, StatusCompleted, nil); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	if err := h.storage.AddBalance(r.Context(), t.ChildID, t.Reward); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	t.Status = StatusCompleted
	respond.JSON(w, http.StatusOK, map[string]any{"task": t})
}

func (h *Handler) Reject(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	t, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if t == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "task not found")
		return
	}
	if t.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}
	if t.Status != StatusPendingReview {
		respond.Error(w, http.StatusBadRequest, "INVALID_STATUS", "task is not pending review")
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	if req.Comment == "" {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "comment is required")
		return
	}

	if err := h.storage.UpdateStatus(r.Context(), id, StatusNeedsRework, &req.Comment); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	t.Status = StatusNeedsRework
	t.RejectionComment = &req.Comment
	respond.JSON(w, http.StatusOK, map[string]any{"task": t})
}