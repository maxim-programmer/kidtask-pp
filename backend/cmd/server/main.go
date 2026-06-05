package main

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"kidtask/internal/admin"
	"kidtask/internal/child"
	"kidtask/internal/config"
	"kidtask/internal/middleware"
	"kidtask/internal/parent"
	"kidtask/internal/respond"
	"kidtask/internal/task"
	"kidtask/internal/wish"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
	log.Println("migrations applied")

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	mw := middleware.Auth(cfg.JWTSecret)

	parentStorage := parent.NewStorage(db)
	childStorage := child.NewStorage(db)
	taskStorage := task.NewStorage(db)
	wishStorage := wish.NewStorage(db)
	adminStorage := admin.NewStorage(db)

	parentHandler := parent.NewHandler(parentStorage, cfg.JWTSecret, mw)
	childHandler := child.NewHandler(childStorage, cfg.JWTSecret, mw)
	taskHandler := task.NewHandler(taskStorage, mw)
	wishHandler := wish.NewHandler(wishStorage, mw)
	adminHandler := admin.NewHandler(adminStorage, mw, cfg.AdminSecret)

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Admin-Secret")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	r.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		var parents, children int
		db.QueryRow(context.Background(), `SELECT COUNT(*) FROM parents`).Scan(&parents)
		db.QueryRow(context.Background(), `SELECT COUNT(*) FROM children`).Scan(&children)
		respond.JSON(w, http.StatusOK, map[string]any{
			"parents":  parents,
			"children": children,
			"total":    parents + children,
		})
	}).Methods(http.MethodGet)

	r.Handle("/api/me", mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r)
		if claims.Role == "parent" {
			p, err := parentStorage.GetByID(r.Context(), claims.UserID)
			if err != nil || p == nil {
				respond.Error(w, http.StatusNotFound, "NOT_FOUND", "not found")
				return
			}
			respond.JSON(w, http.StatusOK, map[string]any{"role": "parent", "user": p})
			return
		}
		c, err := childStorage.GetByID(r.Context(), claims.UserID)
		if err != nil || c == nil {
			respond.Error(w, http.StatusNotFound, "NOT_FOUND", "not found")
			return
		}
		progress, err := childStorage.GetProgress(r.Context(), claims.UserID)
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, "SERVER_ERROR", "server error")
			return
		}
		respond.JSON(w, http.StatusOK, map[string]any{"role": "child", "user": c, "progress": progress})
	}))).Methods(http.MethodGet)

	parentHandler.RegisterRoutes(r)
	childHandler.RegisterRoutes(r)
	taskHandler.RegisterRoutes(r)
	wishHandler.RegisterRoutes(r)
	adminHandler.RegisterRoutes(r)

	log.Printf("server starting on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
		log.Fatal(err)
	}
}