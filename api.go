package main

import (
	// "crypto/ecdsa"
	// "crypto/elliptic"
	// "crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"

	// "log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func (api *APIServer) Run() {
	router := mux.NewRouter()

	// router.HandleFunc("/getDogs", convertHttpHandlerFunc(api.handleGetDogs)).Methods("GET")
	router.HandleFunc("/getUsers", convertHttpHandlerFunc(api.handleGetUsers)).Methods("GET")
	router.HandleFunc("/user/{id}", withJWTAuth(convertHttpHandlerFunc(api.handleGetUserById), api.store)).Methods("GET")
	router.HandleFunc("/createUser", convertHttpHandlerFunc(api.handleCreateUser)).Methods("POST")
	router.HandleFunc("/deleteUser", convertHttpHandlerFunc(api.handleDeleteUser)).Methods("DELETE")
	router.HandleFunc("/login", convertHttpHandlerFunc(api.handleLogin)).Methods("POST")
	fmt.Println("server running")
	http.ListenAndServe(api.listenAddr, router)

}
func NewApiServer(listenAddr string, store Storage) APIServer {
	return APIServer{
		listenAddr,
		store,
	}
}

func (api *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	LoginReq := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(LoginReq); err != nil {
		return err
	}
	resp, err := api.store.GetUserByUsername(LoginReq.Username)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(resp.EncryptedPassword), []byte(LoginReq.Password))

	if err != nil {
		return fmt.Errorf("Incorrect Password")
	}
	token, err := createJwt(resp)
	if err != nil {
		return err
	}
	loginResponse := LoginResponse{
		User:  *resp,
		Token: token,
	}
	return WriteJSON(w, http.StatusOK, loginResponse)
}

func (api *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	userid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	resp, err := api.store.GetUserById(userid)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, resp)
}
func (api *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	resp, err := api.store.GetAllUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, resp)
}
func (api *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	deleteUserReq := new(DeleteUserRequest)
	if err := json.NewDecoder(r.Body).Decode(deleteUserReq); err != nil {
		return err
	}
	if err := api.store.DeleteUser(deleteUserReq.UserID); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "user deleted")
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
	user := NewUser(createUserReq.FirstName, createUserReq.LastName, createUserReq.MailID, createUserReq.Password, createUserReq.Username)
	if err := api.store.CreateUser(user); err != nil {
		return err
	}
	tokenString, err := createJwt(user)
	if err != nil {
		return err
	}
	fmt.Println(tokenString)
	return WriteJSON(w, http.StatusOK, user)
}
func validateJwt(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	fmt.Println(secret)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
}
func createJwt(user *User) (string, error) {
	// key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	secret := os.Getenv("JWT_SECRET")
	claims := &jwt.MapClaims{
		"expiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"userId":    user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}
func withJWTAuth(handlerFunc http.HandlerFunc, api Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("with jwt middleware")
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJwt(tokenString)
		if err != nil {
			fmt.Println("validate jwt error")
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "access denied"})
			return
		}
		if !token.Valid {

			fmt.Println("validat token ")
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "access denied"})
			return
		}

		id := mux.Vars(r)["id"]
		userid, err := strconv.Atoi(id)
		if err != nil {

			fmt.Println("user ID ")
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "access denied"})
			return
		}
		user, err := api.GetUserById(userid)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "access denied"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(reflect.TypeOf(claims["userId"]), reflect.TypeOf(user.ID))
		if user.ID != int(claims["userId"].(float64)) {
			fmt.Println("error here ")
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "access denied"})
			return
		}
		fmt.Println(token)
		handlerFunc(w, r)
	}
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
