package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Processed bool   `json:"processed"`
}

var db *pgxpool.Pool

func main() {
	cfg := loadDBConfig()

	var err error
	db, err = pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()

	if err := migrate(context.Background()); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	http.HandleFunc("/health", withCORS(healthHandler))
	http.HandleFunc("/posts", withCORS(postsHandler))
	http.HandleFunc("/posts/processed/", withCORS(markProcessedHandler))

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func migrate(ctx context.Context) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id        SERIAL PRIMARY KEY,
			title     TEXT        NOT NULL,
			text      TEXT        NOT NULL,
			processed BOOLEAN     NOT NULL DEFAULT FALSE
		)
	`)
	return err
}

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(r.Context()); err != nil {
		http.Error(w, "db unhealthy", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createPost(w, r)
	case http.MethodGet:
		listPosts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.Title == "" || input.Text == "" {
		http.Error(w, "Title and text are required", http.StatusBadRequest)
		return
	}

	var post Post
	err := db.QueryRow(r.Context(),
		`INSERT INTO posts (title, text) VALUES ($1, $2) RETURNING id, title, text, processed`,
		input.Title, input.Text,
	).Scan(&post.ID, &post.Title, &post.Text, &post.Processed)
	if err != nil {
		log.Printf("createPost: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	unprocessedOnly := r.URL.Query().Get("unprocessed") == "true"

	var query string
	if unprocessedOnly {
		query = `SELECT id, title, text, processed FROM posts WHERE processed = FALSE ORDER BY id`
	} else {
		query = `SELECT id, title, text, processed FROM posts WHERE processed = TRUE ORDER BY id`
	}

	rows, err := db.Query(r.Context(), query)
	if err != nil {
		log.Printf("listPosts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	result := make([]Post, 0)
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Text, &p.Processed); err != nil {
			log.Printf("listPosts scan: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		result = append(result, p)
	}
	if err := rows.Err(); err != nil {
		log.Printf("listPosts rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func markProcessedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/posts/processed/")
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var post Post
	err = db.QueryRow(r.Context(),
		`UPDATE posts SET processed = TRUE WHERE id = $1 RETURNING id, title, text, processed`,
		id,
	).Scan(&post.ID, &post.Title, &post.Text, &post.Processed)
	if err != nil {
		// pgx returns pgx.ErrNoRows when no row matched
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
