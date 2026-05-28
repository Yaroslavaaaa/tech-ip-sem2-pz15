package main

import (
	"log"
	"net/http"
	"os"

	"tasks-app/internal/handlers"
	"tasks-app/internal/storage"
)

func main() {
	port := os.Getenv("TASKS_PORT")

	if port == "" {
		port = "8082"
	}

	store := storage.NewMemoryStorage()
	handler := handlers.NewTaskHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/tasks", handler.Tasks)
	mux.HandleFunc("/tasks/", handler.TaskByID)

	log.Println("Tasks service started on port:", port)

	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		log.Fatal(err)
	}
}
