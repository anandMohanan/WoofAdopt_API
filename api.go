package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error
type ApiError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func convertHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (api *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/getDogs", convertHttpHandlerFunc(api.handleGetDogs)).Methods("GET")
	router.HandleFunc("/getUsers", convertHttpHandlerFunc(api.handleGetUsers)).Methods("GET")
	router.HandleFunc("/createUser", convertHttpHandlerFunc(api.handleCreateUser)).Methods("POST")
	fmt.Println("server running")
	http.ListenAndServe(api.listenAddr, router)

}
func NewApiServer(listenAddr string, store Storage) APIServer {
	return APIServer{
		listenAddr,
		store,
	}
}

func (api *APIServer) handleDogs(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// dogName, breed, location, imageUrl, firstName, lastName, mailID string, contactNumber int64)
func (api *APIServer) handleGetDogs(w http.ResponseWriter, r *http.Request) error {
	// dog := NewDog("timmy", "german", "chennai", "https://github.com", "anand", "mohanan", 7448506511)
	return nil
}
func (api *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	resp, err := api.store.GetAllUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, resp)
}

func (api *APIServer) handleCreateDogs(w http.ResponseWriter, r *http.Request) error {
	createDogReq := new(CreateDogRequest)
	if err := json.NewDecoder(r.Body).Decode(createDogReq); err != nil {
		return err
	}
	dog := NewDog(createDogReq.DogName, createDogReq.Breed, createDogReq.Location, createDogReq.ImageURL, createDogReq.Owner, createDogReq.ContactNumber)
	if err := api.store.CreateDog(dog); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, dog)
}

func (api *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}
	user := NewUser(createUserReq.FirstName, createUserReq.LastName, createUserReq.MailID)
	if err := api.store.CreateUser(user); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, user)
}

func (api *APIServer) handleDeleteDogs(w http.ResponseWriter, r *http.Request) error {
	return nil
}
