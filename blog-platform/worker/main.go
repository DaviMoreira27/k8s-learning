package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

const apiURL = "http://localhost:8080"

func main() {
	log.Println("Worker started")

	for {
		processPosts()
		time.Sleep(10 * time.Second)
	}
}

func processPosts() {
	resp, err := http.Get(apiURL + "/posts?unprocessed=true")
	if err != nil {
		log.Println("Error calling API:", err)
		return
	}
	defer resp.Body.Close()

	var posts []Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		log.Println("Error decoding response:", err)
		return
	}

	log.Printf("Fetched %d posts\n", len(posts))

	for _, post := range posts {
		log.Printf("Processing post ID=%d Title=%s\n", post.ID, post.Title)
		time.Sleep(10 * time.Second)
		log.Printf("Updating post ID=%d\n", post.ID)
		updatePost(post.ID)
	}
}

func updatePost(id int) error {
	url := fmt.Sprintf("%s/posts/processed/%d", apiURL, id)

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error with the request: %s", resp.Status)
	}

	return nil
}
