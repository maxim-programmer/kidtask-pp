//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbURL := "postgres://postgres:postgres@db:5432/kidtask?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	rand.Seed(time.Now().UnixNano())

	log.Println("Seeding database...")

	_, err = pool.Exec(context.Background(), "TRUNCATE wishes, tasks, children, parents CASCADE")
	if err != nil {
		log.Fatal(err)
	}

	parents := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		email := "parent_" + strconv.Itoa(i) + "@example.com"
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		var parentID string
		err = pool.QueryRow(context.Background(),
			`INSERT INTO parents (email, password_hash, name) VALUES ($1, $2, $3) RETURNING parent_id`,
			email, hash, "Parent "+strconv.Itoa(i),
		).Scan(&parentID)
		if err != nil {
			log.Fatal(err)
		}
		parents[i] = parentID
		if (i+1)%1000 == 0 {
			log.Printf("Created %d parents", i+1)
		}
	}

	children := make([]string, 20000)
	for i := 0; i < 20000; i++ {
		parentID := parents[rand.Intn(len(parents))]
		username := "child_user_" + strconv.Itoa(i)
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		var childID string
		err = pool.QueryRow(context.Background(),
			`INSERT INTO children (parent_id, username, password_hash, name, age_group, balance) 
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING child_id`,
			parentID, username, hash, "Child "+strconv.Itoa(i), "junior", rand.Intn(10000),
		).Scan(&childID)
		if err != nil {
			log.Fatal(err)
		}
		children[i] = childID
		if (i+1)%2000 == 0 {
			log.Printf("Created %d children", i+1)
		}
	}

	for i := 0; i < 50000; i++ {
		parentID := parents[rand.Intn(len(parents))]
		childID := children[rand.Intn(len(children))]
		_, err = pool.Exec(context.Background(),
			`INSERT INTO tasks (parent_id, child_id, title, reward, status) 
			 VALUES ($1, $2, $3, $4, $5)`,
			parentID, childID, "Task "+strconv.Itoa(i), rand.Intn(500)+1, "active",
		)
		if err != nil {
			log.Fatal(err)
		}
		if (i+1)%5000 == 0 {
			log.Printf("Created %d tasks", i+1)
		}
	}

	for i := 0; i < 30000; i++ {
		childID := children[rand.Intn(len(children))]
		statuses := []string{"awaiting_price", "available", "purchased", "delivered"}
		status := statuses[rand.Intn(len(statuses))]
		price := rand.Intn(1000) + 1
		if status == "awaiting_price" {
			price = 0
		}
		_, err = pool.Exec(context.Background(),
			`INSERT INTO wishes (child_id, title, status, price) 
			 VALUES ($1, $2, $3, $4)`,
			childID, "Wish "+strconv.Itoa(i), status, price,
		)
		if err != nil {
			log.Fatal(err)
		}
		if (i+1)%3000 == 0 {
			log.Printf("Created %d wishes", i+1)
		}
	}

	log.Println("Seeding completed!")
}