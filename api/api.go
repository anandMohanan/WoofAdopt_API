package api

import (
	"fmt"

	"net/http"

	"github.com/anandMohanan/WoofAdopt_API/shared"
	"github.com/anandMohanan/WoofAdopt_API/storage"
	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddr string
	Store      storage.Storage
}

func (api *APIServer) Run() {
	router := mux.NewRouter()

	// router.HandleFunc("/getDogs", convertHttpHandlerFunc(api.handleGetDogs)).Methods("GET")
	router.HandleFunc("/getUsers", convertHttpHandlerFunc(api.handleGetUsers)).Methods("GET")
	router.HandleFunc("/user/{id}", shared.WithJWTAuth(convertHttpHandlerFunc(api.handleGetUserById), api.Store)).Methods("GET")
	router.HandleFunc("/createUser", convertHttpHandlerFunc(api.handleCreateUser)).Methods("POST")
	router.HandleFunc("/deleteUser", convertHttpHandlerFunc(api.handleDeleteUser)).Methods("DELETE")
	router.HandleFunc("/login", convertHttpHandlerFunc(api.handleLogin)).Methods("POST")
	fmt.Println("server running")
	http.ListenAndServe(api.ListenAddr, router)

}
func NewApiServer(listenAddr string, store storage.Storage) APIServer {
	return APIServer{
		listenAddr,
		store,
	}
}

func convertHttpHandlerFunc(f shared.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			shared.WriteJSON(w, http.StatusBadRequest, shared.ApiError{Error: err.Error()})
		}
	}
}
