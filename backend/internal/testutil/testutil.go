package testutil

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupTestDB создаёт тестовую базу данных
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/kidtask_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to test database: %v", err)
	}

	// Очищаем таблицы перед тестами
	CleanDB(t, pool)

	return pool
}

// CleanDB очищает все таблицы
func CleanDB(t *testing.T, db *pgxpool.Pool) {
	tables := []string{"wishes", "tasks", "children", "parents"}
	for _, table := range tables {
		_, err := db.Exec(context.Background(), "TRUNCATE TABLE "+table+" CASCADE")
		if err != nil {
			t.Logf("Warning: could not truncate %s: %v", table, err)
		}
	}
}

// UniqueEmail генерирует уникальный email для тестов
func UniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%d@example.com", prefix, time.Now().UnixNano())
}

// UniqueUsername генерирует уникальный username для тестов
func UniqueUsername(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// CreateTestParent создаёт тестового родителя
func CreateTestParent(t *testing.T, db *pgxpool.Pool, email, password, name string) string {
	var parentID string
	err := db.QueryRow(context.Background(),
		`INSERT INTO parents (email, password_hash, name) 
		 VALUES ($1, $2, $3) RETURNING parent_id`,
		email, password, name,
	).Scan(&parentID)
	if err != nil {
		t.Fatalf("Failed to create test parent: %v", err)
	}
	return parentID
}

// CreateTestChild создаёт тестового ребёнка
func CreateTestChild(t *testing.T, db *pgxpool.Pool, parentID, username, passwordHash, name, ageGroup string) string {
	var childID string
	err := db.QueryRow(context.Background(),
		`INSERT INTO children (parent_id, username, password_hash, name, age_group, balance) 
		 VALUES ($1, $2, $3, $4, $5, 0) RETURNING child_id`,
		parentID, username, passwordHash, name, ageGroup,
	).Scan(&childID)
	if err != nil {
		t.Fatalf("Failed to create test child: %v", err)
	}
	return childID
}

// SetupTestRouter создаёт роутер с middleware для тестов
func SetupTestRouter(handlerFunc func(*mux.Router)) *mux.Router {
	r := mux.NewRouter()
	handlerFunc(r)
	return r
}
