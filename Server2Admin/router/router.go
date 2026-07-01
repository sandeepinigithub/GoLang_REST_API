package router

import (
	"net/http"

	"github.com/sandeepshukla/server2admin/handlers"
)

func NewRouter(h *handlers.AdminHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetAllUsers(w, r)
		case http.MethodPost:
			h.CreateUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/admin/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetUser(w, r)
		case http.MethodPut:
			h.UpdateUser(w, r)
		case http.MethodDelete:
			h.DeleteUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
