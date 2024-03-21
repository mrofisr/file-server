package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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
	log.Println(filesDir)
	FileServer(r, "/files", filesDir)
	UploadFile(r, "/files", filesDir)
	DeleteFile(r, "/files", filesDir)
	log.Println("server is running on port: ", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters")
	}
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			http.Error(w, "filename is required", http.StatusBadRequest)
			return
		}
		file, err := root.Open(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, fi.Name(), fi.ModTime(), file)
	})
}

func UploadFile(r chi.Router, path string, root http.FileSystem) {
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		targetPath := filepath.Join(fmt.Sprintf("%v", root), handler.Filename)
		f, err := os.Create(targetPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("file uploaded successfully"))
	})
}

func DeleteFile(r chi.Router, path string, root http.FileSystem) {
	r.Delete(path, func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		err := os.Remove(filepath.Join(fmt.Sprintf("%v", root), filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("file deleted successfully"))
	})
}
