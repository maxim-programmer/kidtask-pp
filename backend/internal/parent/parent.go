package parent

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"kidtask/internal/auth"
	"kidtask/internal/middleware"
	"kidtask/internal/respond"
)

type Parent struct {
	ParentID     string `json:"parent_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Name         string `json:"name"`
}

func (p *Parent) Validate() error {
	if strings.TrimSpace(p.Email) == "" || !strings.Contains(p.Email, "@") {
		return errors.New("invalid email")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(ctx context.Context, p *Parent) error {
	return s.db.QueryRow(ctx,
		`INSERT INTO parents (email, password_hash, name)
		 VALUES ($1, $2, $3) RETURNING parent_id`,
		p.Email, p.PasswordHash, p.Name,
	).Scan(&p.ParentID)
}

func (s *Storage) GetByEmail(ctx context.Context, email string) (*Parent, error) {
	p := &Parent{}
	err := s.db.QueryRow(ctx,
		`SELECT parent_id, email, password_hash, name FROM parents WHERE email = $1`, email,
	).Scan(&p.ParentID, &p.Email, &p.PasswordHash, &p.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (s *Storage) GetByID(ctx context.Context, id string) (*Parent, error) {
	p := &Parent{}
	err := s.db.QueryRow(ctx,
		`SELECT parent_id, email, password_hash, name FROM parents WHERE parent_id = $1`, id,
	).Scan(&p.ParentID, &p.Email, &p.PasswordHash, &p.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (s *Storage) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM parents WHERE email = $1)`, email,
	).Scan(&exists)
	return exists, err
}

func (s *Storage) EmailExistsExcluding(ctx context.Context, email, excludeID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM parents WHERE email = $1 AND parent_id != $2)`, email, excludeID,
	).Scan(&exists)
	return exists, err
}

func (s *Storage) Update(ctx context.Context, p *Parent) error {
	_, err := s.db.Exec(ctx,
		`UPDATE parents SET email=$1, password_hash=$2, name=$3 WHERE parent_id=$4`,
		p.Email, p.PasswordHash, p.Name, p.ParentID,
	)
	return err
}

type Handler struct {
	storage   *Storage
	jwtSecret string
	mw        func(http.Handler) http.Handler
}

func NewHandler(storage *Storage, jwtSecret string, mw func(http.Handler) http.Handler) *Handler {
	return &Handler{storage: storage, jwtSecret: jwtSecret, mw: mw}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/auth/register", h.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login", h.Login).Methods(http.MethodPost)
	r.Handle("/api/me", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.UpdateMe)))).Methods(http.MethodPatch)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}

	p := &Parent{Email: req.Email, Name: req.Name}
	if err := p.Validate(); err != nil {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if len(req.Password) < 6 {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "password too short")
		return
	}

	exists, err := h.storage.EmailExists(r.Context(), req.Email)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if exists {
		respond.Error(w, http.StatusConflict, "EMAIL_TAKEN", "email already taken")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	p.PasswordHash = hash

	if err := h.storage.Create(r.Context(), p); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	token, err := auth.GenerateToken(p.ParentID, "parent", "", h.jwtSecret)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]any{"token": token, "user": p})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	if req.Email == "" || req.Password == "" {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "email and password required")
		return
	}

	p, err := h.storage.GetByEmail(r.Context(), req.Email)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if p == nil || !auth.CheckPassword(p.PasswordHash, req.Password) {
		respond.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
		return
	}

	token, err := auth.GenerateToken(p.ParentID, "parent", "", h.jwtSecret)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"token": token, "user": p})
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	p, err := h.storage.GetByID(r.Context(), claims.UserID)
	if err != nil || p == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "parent not found")
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "name is required")
			return
		}
		p.Name = *req.Name
	}

	if req.Email != nil {
		if strings.TrimSpace(*req.Email) == "" || !strings.Contains(*req.Email, "@") {
			respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid email")
			return
		}
		if *req.Email != p.Email {
			exists, err := h.storage.EmailExistsExcluding(r.Context(), *req.Email, p.ParentID)
			if err != nil {
				respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
				return
			}
			if exists {
				respond.Error(w, http.StatusConflict, "EMAIL_TAKEN", "email already taken")
				return
			}
		}
		p.Email = *req.Email
	}

	if req.Password != nil {
		if len(*req.Password) < 6 {
			respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "password too short")
			return
		}
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
			return
		}
		p.PasswordHash = hash
	}

	if err := h.storage.Update(r.Context(), p); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"user": p})
}