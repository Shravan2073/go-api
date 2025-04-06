package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"math/rand"
	"time"
	"github.com/gorilla/mux"
)

// Function to get all files in a specific directory (images and GIFs)
func getMemesFromDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if file is an image or GIF
		if !info.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg" || filepath.Ext(path) == ".png" || filepath.Ext(path) == ".gif") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

// API handler to serve a random meme
func getRandomMeme(w http.ResponseWriter, r *http.Request) {
	// Directory where your memes (images and GIFs) are stored
	memeDirectory := "./memes"
	files, err := getMemesFromDirectory(memeDirectory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading memes directory: %v", err), http.StatusInternalServerError)
		return
	}

	// If there are no memes found, return an error message
	if len(files) == 0 {
		http.Error(w, "No memes found!", http.StatusNotFound)
		return
	}

	// Select a random meme from the list of files
	rand.Seed(time.Now().UnixNano())
	randomMeme := files[rand.Intn(len(files))]

	// Serve the random meme as a response
	http.ServeFile(w, r, randomMeme)
}

func main() {
	r := mux.NewRouter()

	// Route to get a random meme
	r.HandleFunc("/random-meme", getRandomMeme).Methods("GET")

	// Start the server
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
