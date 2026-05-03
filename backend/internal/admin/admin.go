package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"kidtask/internal/respond"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

type FamilyStat struct {
	ParentID    string `json:"parent_id"`
	ParentName  string `json:"parent_name"`
	ParentEmail string `json:"parent_email"`
	IsBlocked   bool   `json:"is_blocked"`
	ChildCount  int    `json:"child_count"`
	TasksDone   int    `json:"tasks_done"`
}

type ChildRow struct {
	ChildID    string `json:"child_id"`
	ParentID   string `json:"parent_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Balance    int    `json:"balance"`
	IsBlocked  bool   `json:"is_blocked"`
	ParentName string `json:"parent_name"`
}

type BalanceLog struct {
	LogID     string    `json:"log_id"`
	ChildID   string    `json:"child_id"`
	Delta     int       `json:"delta"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type Complaint struct {
	ComplaintID string    `json:"complaint_id"`
	ParentID    string    `json:"parent_id"`
	ParentName  string    `json:"parent_name"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type ChatMessage struct {
	MessageID string    `json:"message_id"`
	ParentID  string    `json:"parent_id"`
	FromAdmin bool      `json:"from_admin"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type GlobalStats struct {
	TotalFamilies  int     `json:"total_families"`
	TotalChildren  int     `json:"total_children"`
	TotalTasks     int     `json:"total_tasks"`
	CompletedTasks int     `json:"completed_tasks"`
	AvgTasksDone   float64 `json:"avg_tasks_done"`
	TotalWishes    int     `json:"total_wishes"`
	TotalBalance   int     `json:"total_balance"`
}

type WishRow struct {
	WishID     string    `json:"wish_id"`
	ChildID    string    `json:"child_id"`
	ChildName  string    `json:"child_name"`
	ParentName string    `json:"parent_name"`
	Title      string    `json:"title"`
	Price      *int      `json:"price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *Storage) GetGlobalStats(ctx context.Context) (*GlobalStats, error) {
	st := &GlobalStats{}
	err := s.db.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM parents),
			(SELECT COUNT(*) FROM children),
			(SELECT COUNT(*) FROM tasks),
			(SELECT COUNT(*) FROM tasks WHERE status='completed'),
			(SELECT COUNT(*) FROM wishes),
			(SELECT COALESCE(SUM(balance), 0) FROM children)
	`).Scan(&st.TotalFamilies, &st.TotalChildren, &st.TotalTasks, &st.CompletedTasks, &st.TotalWishes, &st.TotalBalance)
	if err != nil {
		return nil, err
	}
	if st.TotalChildren > 0 {
		st.AvgTasksDone = float64(st.CompletedTasks) / float64(st.TotalChildren)
	}
	return st, nil
}

func (s *Storage) ListFamilies(ctx context.Context) ([]*FamilyStat, error) {
	rows, err := s.db.Query(ctx, `
		SELECT p.parent_id, p.name, p.email, p.is_blocked,
		       COUNT(DISTINCT c.child_id),
		       COUNT(t.task_id) FILTER (WHERE t.status = 'completed')
		FROM parents p
		LEFT JOIN children c ON c.parent_id = p.parent_id
		LEFT JOIN tasks t ON t.parent_id = p.parent_id
		GROUP BY p.parent_id
		ORDER BY p.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*FamilyStat
	for rows.Next() {
		f := &FamilyStat{}
		if err := rows.Scan(&f.ParentID, &f.ParentName, &f.ParentEmail, &f.IsBlocked, &f.ChildCount, &f.TasksDone); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, nil
}

func (s *Storage) SetParentBlocked(ctx context.Context, parentID string, blocked bool) error {
	_, err := s.db.Exec(ctx, `UPDATE parents SET is_blocked=$1 WHERE parent_id=$2`, blocked, parentID)
	return err
}

func (s *Storage) DeleteParent(ctx context.Context, parentID string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM parents WHERE parent_id=$1`, parentID)
	return err
}

func (s *Storage) ListChildren(ctx context.Context) ([]*ChildRow, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.child_id, c.parent_id, c.name, c.username, c.balance, c.is_blocked, p.name
		FROM children c JOIN parents p ON p.parent_id = c.parent_id
		ORDER BY p.name, c.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*ChildRow
	for rows.Next() {
		c := &ChildRow{}
		if err := rows.Scan(&c.ChildID, &c.ParentID, &c.Name, &c.Username, &c.Balance, &c.IsBlocked, &c.ParentName); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (s *Storage) SetChildBlocked(ctx context.Context, childID string, blocked bool) error {
	_, err := s.db.Exec(ctx, `UPDATE children SET is_blocked=$1 WHERE child_id=$2`, blocked, childID)
	return err
}

func (s *Storage) AdjustBalance(ctx context.Context, childID string, delta int, reason string) (int, error) {
	var newBalance int
	err := s.db.QueryRow(ctx,
		`UPDATE children SET balance = GREATEST(0, balance + $1) WHERE child_id=$2 RETURNING balance`,
		delta, childID,
	).Scan(&newBalance)
	if err != nil {
		return 0, err
	}
	_, err = s.db.Exec(ctx,
		`INSERT INTO balance_logs (child_id, delta, reason) VALUES ($1, $2, $3)`,
		childID, delta, reason,
	)
	return newBalance, err
}

func (s *Storage) GetBalanceLogs(ctx context.Context, childID string) ([]*BalanceLog, error) {
	rows, err := s.db.Query(ctx,
		`SELECT log_id, child_id, delta, reason, created_at FROM balance_logs WHERE child_id=$1 ORDER BY created_at DESC`,
		childID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*BalanceLog
	for rows.Next() {
		l := &BalanceLog{}
		if err := rows.Scan(&l.LogID, &l.ChildID, &l.Delta, &l.Reason, &l.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (s *Storage) ListComplaints(ctx context.Context) ([]*Complaint, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.complaint_id, c.parent_id, p.name, c.subject, c.body, c.status, c.created_at
		FROM complaints c JOIN parents p ON p.parent_id = c.parent_id
		ORDER BY c.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Complaint
	for rows.Next() {
		c := &Complaint{}
		if err := rows.Scan(&c.ComplaintID, &c.ParentID, &c.ParentName, &c.Subject, &c.Body, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (s *Storage) CreateComplaint(ctx context.Context, parentID, subject, body string) (*Complaint, error) {
	c := &Complaint{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO complaints (parent_id, subject, body) VALUES ($1,$2,$3) RETURNING complaint_id, parent_id, subject, body, status, created_at`,
		parentID, subject, body,
	).Scan(&c.ComplaintID, &c.ParentID, &c.Subject, &c.Body, &c.Status, &c.CreatedAt)
	return c, err
}

func (s *Storage) ResolveComplaint(ctx context.Context, complaintID string) error {
	_, err := s.db.Exec(ctx, `UPDATE complaints SET status='resolved' WHERE complaint_id=$1`, complaintID)
	return err
}

func (s *Storage) GetChatMessages(ctx context.Context, parentID string) ([]*ChatMessage, error) {
	rows, err := s.db.Query(ctx,
		`SELECT message_id, parent_id, from_admin, body, created_at FROM chat_messages WHERE parent_id=$1 ORDER BY created_at ASC`,
		parentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*ChatMessage
	for rows.Next() {
		m := &ChatMessage{}
		if err := rows.Scan(&m.MessageID, &m.ParentID, &m.FromAdmin, &m.Body, &m.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (s *Storage) SendChatMessage(ctx context.Context, parentID, body string, fromAdmin bool) (*ChatMessage, error) {
	m := &ChatMessage{}
	err := s.db.QueryRow(ctx,
		`INSERT INTO chat_messages (parent_id, from_admin, body) VALUES ($1,$2,$3) RETURNING message_id, parent_id, from_admin, body, created_at`,
		parentID, fromAdmin, body,
	).Scan(&m.MessageID, &m.ParentID, &m.FromAdmin, &m.Body, &m.CreatedAt)
	return m, err
}

func (s *Storage) ListWishes(ctx context.Context, sortBy string) ([]*WishRow, error) {
	orderClause := "w.created_at DESC"
	switch sortBy {
	case "price_asc":
		orderClause = "w.price ASC NULLS LAST"
	case "price_desc":
		orderClause = "w.price DESC NULLS LAST"
	case "oldest":
		orderClause = "w.created_at ASC"
	}
	rows, err := s.db.Query(ctx, `
		SELECT w.wish_id, w.child_id, c.name, p.name, w.title, w.price, w.status, w.created_at
		FROM wishes w
		JOIN children c ON c.child_id = w.child_id
		JOIN parents p ON p.parent_id = c.parent_id
		ORDER BY `+orderClause,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*WishRow
	for rows.Next() {
		w := &WishRow{}
		if err := rows.Scan(&w.WishID, &w.ChildID, &w.ChildName, &w.ParentName, &w.Title, &w.Price, &w.Status, &w.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	return list, nil
}

func (s *Storage) ListParentsForChat(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.db.Query(ctx, `
		SELECT p.parent_id, p.name, p.email,
		       COUNT(m.message_id) FILTER (WHERE m.from_admin = false) as unread
		FROM parents p
		LEFT JOIN chat_messages m ON m.parent_id = p.parent_id
		GROUP BY p.parent_id
		ORDER BY p.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var id, name, email string
		var unread int
		if err := rows.Scan(&id, &name, &email, &unread); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{
			"parent_id": id,
			"name":      name,
			"email":     email,
			"unread":    unread,
		})
	}
	return list, nil
}

type Handler struct {
	storage     *Storage
	mw          func(http.Handler) http.Handler
	adminSecret string
}

func NewHandler(storage *Storage, mw func(http.Handler) http.Handler, adminSecret string) *Handler {
	return &Handler{storage: storage, mw: mw, adminSecret: adminSecret}
}

func (h *Handler) adminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Admin-Secret") != h.adminSecret {
			respond.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid admin secret")
			return
		}
		next(w, r)
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	a := r.PathPrefix("/api/admin").Subrouter()

	a.HandleFunc("/stats", h.adminAuth(h.Stats)).Methods(http.MethodGet)

	a.HandleFunc("/families", h.adminAuth(h.ListFamilies)).Methods(http.MethodGet)
	a.HandleFunc("/families/{id}/block", h.adminAuth(h.BlockParent)).Methods(http.MethodPost)
	a.HandleFunc("/families/{id}/unblock", h.adminAuth(h.UnblockParent)).Methods(http.MethodPost)
	a.HandleFunc("/families/{id}", h.adminAuth(h.DeleteFamily)).Methods(http.MethodDelete)

	a.HandleFunc("/children", h.adminAuth(h.ListChildren)).Methods(http.MethodGet)
	a.HandleFunc("/children/{id}/block", h.adminAuth(h.BlockChild)).Methods(http.MethodPost)
	a.HandleFunc("/children/{id}/unblock", h.adminAuth(h.UnblockChild)).Methods(http.MethodPost)
	a.HandleFunc("/children/{id}/balance", h.adminAuth(h.AdjustBalance)).Methods(http.MethodPost)
	a.HandleFunc("/children/{id}/logs", h.adminAuth(h.BalanceLogs)).Methods(http.MethodGet)

	a.HandleFunc("/complaints", h.adminAuth(h.ListComplaints)).Methods(http.MethodGet)
	a.HandleFunc("/complaints/{id}/resolve", h.adminAuth(h.ResolveComplaint)).Methods(http.MethodPost)

	a.HandleFunc("/chat", h.adminAuth(h.ListParentsForChat)).Methods(http.MethodGet)
	a.HandleFunc("/chat/{parent_id}", h.adminAuth(h.GetChat)).Methods(http.MethodGet)
	a.HandleFunc("/chat/{parent_id}", h.adminAuth(h.AdminSendMessage)).Methods(http.MethodPost)

	a.HandleFunc("/wishes", h.adminAuth(h.ListWishes)).Methods(http.MethodGet)

	r.Handle("/api/support/chat", h.mw(http.HandlerFunc(h.UserGetChat))).Methods(http.MethodGet)
	r.Handle("/api/support/chat", h.mw(http.HandlerFunc(h.UserSendMessage))).Methods(http.MethodPost)
	r.Handle("/api/support/complaints", h.mw(http.HandlerFunc(h.UserCreateComplaint))).Methods(http.MethodPost)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	st, err := h.storage.GetGlobalStats(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, st)
}

func (h *Handler) ListFamilies(w http.ResponseWriter, r *http.Request) {
	list, err := h.storage.ListFamilies(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if list == nil {
		list = []*FamilyStat{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"families": list})
}

func (h *Handler) BlockParent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.storage.SetParentBlocked(r.Context(), id, true); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"message": "blocked"})
}

func (h *Handler) UnblockParent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.storage.SetParentBlocked(r.Context(), id, false); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"message": "unblocked"})
}

func (h *Handler) DeleteFamily(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.storage.DeleteParent(r.Context(), id); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"message": "deleted"})
}

func (h *Handler) ListChildren(w http.ResponseWriter, r *http.Request) {
	list, err := h.storage.ListChildren(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if list == nil {
		list = []*ChildRow{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"children": list})
}

func (h *Handler) BlockChild(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.storage.SetChildBlocked(r.Context(), id, true); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"message": "blocked"})
}

func (h *Handler) UnblockChild(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.storage.SetChildBlocked(r.Context(), id, false); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"message": "unblocked"})
}

func (h *Handler) AdjustBalance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Delta  int    `json:"delta"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	if req.Delta == 0 {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "delta must not be zero")
		return
	}
	if req.Reason == "" {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "reason is required")
		return
	}
	newBalance, err := h.storage.AdjustBalance(r.Context(), id, req.Delta, req.Reason)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"balance": newBalance})
}

func (h *Handler) BalanceLogs(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	logs, err := h.storage.GetBalanceLogs(r.Context(), id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if logs == nil {
		logs = []*BalanceLog{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"logs": logs})
}

func (h *Handler) ListComplaints(w http.ResponseWriter, r *http.Request) {
	list, err := h.storage.ListComplaints(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if list == nil {
		list = []*Complaint{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"complaints": list})
}

func (h *Handler) ResolveComplaint(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.storage.ResolveComplaint(r.Context(), id); err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"message": "resolved"})
}

func (h *Handler) ListParentsForChat(w http.ResponseWriter, r *http.Request) {
	list, err := h.storage.ListParentsForChat(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if list == nil {
		list = []map[string]any{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"parents": list})
}

func (h *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
	parentID := mux.Vars(r)["parent_id"]
	msgs, err := h.storage.GetChatMessages(r.Context(), parentID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if msgs == nil {
		msgs = []*ChatMessage{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"messages": msgs})
}

func (h *Handler) AdminSendMessage(w http.ResponseWriter, r *http.Request) {
	parentID := mux.Vars(r)["parent_id"]
	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Body == "" {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "body is required")
		return
	}
	msg, err := h.storage.SendChatMessage(r.Context(), parentID, req.Body, true)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]any{"message": msg})
}

func (h *Handler) ListWishes(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	list, err := h.storage.ListWishes(r.Context(), sortBy)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if list == nil {
		list = []*WishRow{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"wishes": list})
}

func (h *Handler) UserGetChat(w http.ResponseWriter, r *http.Request) {
	parentID := r.URL.Query().Get("parent_id")
	if parentID == "" {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "parent_id required")
		return
	}
	msgs, err := h.storage.GetChatMessages(r.Context(), parentID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	if msgs == nil {
		msgs = []*ChatMessage{}
	}
	respond.JSON(w, http.StatusOK, map[string]any{"messages": msgs})
}

func (h *Handler) UserSendMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParentID string `json:"parent_id"`
		Body     string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Body == "" || req.ParentID == "" {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "parent_id and body are required")
		return
	}
	msg, err := h.storage.SendChatMessage(r.Context(), req.ParentID, req.Body, false)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]any{"message": msg})
}

func (h *Handler) UserCreateComplaint(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParentID string `json:"parent_id"`
		Subject  string `json:"subject"`
		Body     string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "BAD_REQUEST", "invalid body")
		return
	}
	if req.ParentID == "" || req.Subject == "" || req.Body == "" {
		respond.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "parent_id, subject and body are required")
		return
	}
	c, err := h.storage.CreateComplaint(r.Context(), req.ParentID, req.Subject, req.Body)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
		return
	}
	respond.JSON(w, http.StatusCreated, map[string]any{"complaint": c})
}