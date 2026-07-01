package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sandeepshukla/golangrestproject1/handlers"
	"github.com/sandeepshukla/golangrestproject1/router"
	"github.com/sandeepshukla/golangrestproject1/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	userStore := store.NewUserStore()
	userHandler := handlers.NewUserHandler(userStore)
	r := router.NewRouter(userHandler)

	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /users       - list all users")
	fmt.Println("  POST   /users       - create a user")
	fmt.Println("  GET    /users/{id}  - get user by id")
	fmt.Println("  PUT    /users/{id}  - update user by id")
	fmt.Println("  DELETE /users/{id}  - delete user by id")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
