package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mrofisr/files-server/lib/files"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("alhamdulillah server is running!"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, os.Getenv("DIRECTORY")))
	files.FileServer(r, "/files", filesDir)
	files.UploadFile(r, "/files", filesDir)
	files.DeleteFile(r, "/files", filesDir)
	log.Println("server is running on port: ", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
