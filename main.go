package main

import (
	"fmt"
	"log"
	"net/http"

	"todo-app/handlers"
	"todo-app/store"
)

func main() {
	todoStore := store.NewTodoStore()
	todoHandler := handlers.NewTodoHandler(todoStore)

	mux := http.NewServeMux()
	todoHandler.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("To-Do API server running on http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop.")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
