package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sandeepshukla/server2admin/client"
	"github.com/sandeepshukla/server2admin/handlers"
	"github.com/sandeepshukla/server2admin/router"
)

func main() {
	// Address of Server1User — override with SERVER1_URL env var if needed
	server1URL := os.Getenv("SERVER1_URL")
	if server1URL == "" {
		server1URL = "http://localhost:8080"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port

	// Wire up: HTTP client → handler → router
	userClient := client.NewUserClient(server1URL)
	adminHandler := handlers.NewAdminHandler(userClient)
	r := router.NewRouter(adminHandler)

	fmt.Printf("Server2Admin starting on http://localhost%s\n", addr)
	fmt.Printf("Forwarding requests to Server1User at %s\n", server1URL)
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /admin/users       - list all users  (via Server1User)")
	fmt.Println("  POST   /admin/users       - create a user   (via Server1User)")
	fmt.Println("  GET    /admin/users/{id}  - get user by id  (via Server1User)")
	fmt.Println("  PUT    /admin/users/{id}  - update user     (via Server1User)")
	fmt.Println("  DELETE /admin/users/{id}  - delete user     (via Server1User)")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server2Admin failed to start: %v", err)
	}
}
