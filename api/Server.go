package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/janicaleksander/BeMotivated/storage"
	"github.com/janicaleksander/BeMotivated/types"
	"log"
	"net/http"
)

type APIServer struct {
	ListenAddress string
	Store         storage.Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func BuildServer(address string, storage storage.Storage) (s *APIServer) {
	return &APIServer{
		ListenAddress: address,
		Store:         storage,
	}
}
func WriteToJson(w http.ResponseWriter, statusCode int, val any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(val)
}
func prepareHandle(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.FuncProb})
		}
	}

}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/register", prepareHandle(s.handleRegister))
	router.HandleFunc("/api/login", prepareHandle(s.handleLogin))
	router.HandleFunc("/api/dashboard", prepareHandle(s.handleDashboard))
	router.HandleFunc("/api/logout", prepareHandle(s.handleLogOut))
	log.Println("Running on: ", s.ListenAddress)
	http.ListenAndServe(s.ListenAddress, router)

}
