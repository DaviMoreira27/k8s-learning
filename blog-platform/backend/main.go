package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Processed bool   `json:"processed"`
}

var (
	posts  []Post
	nextID = 1
	mutex  sync.Mutex
)

func main() {
	http.HandleFunc("/health", withCORS(healthHandler))
	http.HandleFunc("/posts", withCORS(postsHandler))
	http.HandleFunc("/posts/", withCORS(markProcessedHandler))

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Title == "" || input.Text == "" {
		http.Error(w, "Title and text are required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	post := Post{
		ID:        nextID,
		Title:     input.Title,
		Text:      input.Text,
		Processed: false,
	}
	nextID++
	posts = append(posts, post)
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	if r.URL.Query().Get("unprocessed") == "true" {
		var result []Post
		for _, p := range posts {
			if !p.Processed {
				result = append(result, p)
			}
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func markProcessedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i := range posts {
		if posts[i].ID == id {
			posts[i].Processed = true
			json.NewEncoder(w).Encode(posts[i])
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}
