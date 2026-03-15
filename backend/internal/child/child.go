package child

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

type Child struct {
	ChildID      string `json:"child_id"`
	ParentID     string `json:"parent_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Name         string `json:"name"`
	Balance      int    `json:"balance"`
	AgeGroup     string `json:"age_group"`
}

func (c *Child) Validate() error {
	if strings.TrimSpace(c.Username) == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if c.AgeGroup != "junior" && c.AgeGroup != "senior" {
		return errors.New("age_group must be junior or senior")
	}
	return nil
}

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(ctx context.Context, c *Child) error {
	return s.db.QueryRow(ctx,
		`INSERT INTO children (parent_id, username, password_hash, name, age_group)
		 VALUES ($1, $2, $3, $4, $5) RETURNING child_id, balance`,
		c.ParentID, c.Username, c.PasswordHash, c.Name, c.AgeGroup,
	).Scan(&c.ChildID, &c.Balance)
}

func (s *Storage) GetByID(ctx context.Context, id string) (*Child, error) {
	c := &Child{}
	err := s.db.QueryRow(ctx,
		`SELECT child_id, parent_id, username, password_hash, name, balance, age_group
		 FROM children WHERE child_id = $1`, id,
	).Scan(&c.ChildID, &c.ParentID, &c.Username, &c.PasswordHash, &c.Name, &c.Balance, &c.AgeGroup)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return c, err
}

func (s *Storage) GetByUsername(ctx context.Context, username string) (*Child, error) {
	c := &Child{}
	err := s.db.QueryRow(ctx,
		`SELECT child_id, parent_id, username, password_hash, name, balance, age_group
		 FROM children WHERE username = $1`, username,
	).Scan(&c.ChildID, &c.ParentID, &c.Username, &c.PasswordHash, &c.Name, &c.Balance, &c.AgeGroup)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return c, err
}

func (s *Storage) ListByParent(ctx context.Context, parentID string) ([]*Child, error) {
	rows, err := s.db.Query(ctx,
		`SELECT child_id, parent_id, username, password_hash, name, balance, age_group
		 FROM children WHERE parent_id = $1`, parentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []*Child
	for rows.Next() {
		c := &Child{}
		if err := rows.Scan(&c.ChildID, &c.ParentID, &c.Username, &c.PasswordHash, &c.Name, &c.Balance, &c.AgeGroup); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, nil
}

func (s *Storage) Update(ctx context.Context, c *Child) error {
	_, err := s.db.Exec(ctx,
		`UPDATE children SET username=$1, password_hash=$2, name=$3, age_group=$4 WHERE child_id=$5`,
		c.Username, c.PasswordHash, c.Name, c.AgeGroup, c.ChildID,
	)
	return err
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM children WHERE child_id = $1`, id)
	return err
}

func (s *Storage) UsernameExists(ctx context.Context, username, parentID, excludeID string) (bool, error) {
	var exists bool
	var err error
	if excludeID == "" {
		err = s.db.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM children WHERE username=$1 AND parent_id=$2)`,
			username, parentID,
		).Scan(&exists)
	} else {
		err = s.db.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM children WHERE username=$1 AND parent_id=$2 AND child_id != $3)`,
			username, parentID, excludeID,
		).Scan(&exists)
	}
	return exists, err
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
	r.HandleFunc("/api/auth/child/login", h.Login).Methods(http.MethodPost)

	r.Handle("/api/children", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.List)))).Methods(http.MethodGet)
	r.Handle("/api/children", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Create)))).Methods(http.MethodPost)
	r.Handle("/api/children/{id}", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Get)))).Methods(http.MethodGet)
	r.Handle("/api/children/{id}", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Update)))).Methods(http.MethodPatch)
	r.Handle("/api/children/{id}", h.mw(http.HandlerFunc(
		middleware.RequireRole("parent", h.Delete)))).Methods(http.MethodDelete)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	if req.Username == "" || req.Password == "" {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "username and password required")
		return
	}

	c, err := h.storage.GetByUsername(r.Context(), req.Username)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if c == nil || !auth.CheckPassword(c.PasswordHash, req.Password) {
		respond.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password")
		return
	}

	token, err := auth.GenerateToken(c.ChildID, "child", c.ParentID, h.jwtSecret)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"token": token, "child": c})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	children, err := h.storage.ListByParent(r.Context(), claims.UserID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if children == nil {
		children = []*Child{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"children": children})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	var req struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
		AgeGroup string `json:"age_group"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	if req.AgeGroup == "" {
		req.AgeGroup = "junior"
	}

	c := &Child{ParentID: claims.UserID, Username: req.Username, Name: req.Name, AgeGroup: req.AgeGroup}
	if err := c.Validate(); err != nil {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if len(req.Password) < 4 {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "password too short")
		return
	}

	exists, err := h.storage.UsernameExists(r.Context(), req.Username, claims.UserID, "")
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if exists {
		respond.Error(w, http.StatusConflict, "USERNAME_TAKEN", "username already taken")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	c.PasswordHash = hash

	if err := h.storage.Create(r.Context(), c); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]any{"child": c})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	c, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if c == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "child not found")
		return
	}
	if c.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"child": c})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	c, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if c == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "child not found")
		return
	}
	if c.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	var req struct {
		Name     *string `json:"name"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		AgeGroup *string `json:"age_group"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.AgeGroup != nil {
		c.AgeGroup = *req.AgeGroup
	}
	if req.Username != nil && *req.Username != c.Username {
		exists, err := h.storage.UsernameExists(r.Context(), *req.Username, claims.UserID, c.ChildID)
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
			return
		}
		if exists {
			respond.Error(w, http.StatusConflict, "USERNAME_TAKEN", "username already taken")
			return
		}
		c.Username = *req.Username
	}
	if req.Password != nil {
		if len(*req.Password) < 4 {
			respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "password too short")
			return
		}
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
			return
		}
		c.PasswordHash = hash
	}

	if err := h.storage.Update(r.Context(), c); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"child": c})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := mux.Vars(r)["id"]

	c, err := h.storage.GetByID(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if c == nil {
		respond.Error(w, http.StatusNotFound, "NOT_FOUND", "child not found")
		return
	}
	if c.ParentID != claims.UserID {
		respond.Error(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	if err := h.storage.Delete(r.Context(), id); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{"message": "child deleted"})
}