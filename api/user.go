package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anandMohanan/WoofAdopt_API/models"
	"github.com/anandMohanan/WoofAdopt_API/shared"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (api *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	LoginReq := new(models.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(LoginReq); err != nil {
		return err
	}
	resp, err := api.Store.GetUserByUsername(LoginReq.Username)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(resp.EncryptedPassword), []byte(LoginReq.Password))

	if err != nil {
		return fmt.Errorf("Incorrect Password")
	}
	token, err := shared.CreateJwt(resp)
	if err != nil {
		return err
	}
	loginResponse := models.LoginResponse{
		User:  *resp,
		Token: token,
	}
	return shared.WriteJSON(w, http.StatusOK, loginResponse)
}

func (api *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	userid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	resp, err := api.Store.GetUserById(userid)
	if err != nil {
		return err
	}
	return shared.WriteJSON(w, http.StatusOK, resp)
}
func (api *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	resp, err := api.Store.GetAllUsers()
	if err != nil {
		return err
	}
	return shared.WriteJSON(w, http.StatusOK, resp)
}
func (api *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	deleteUserReq := new(models.DeleteUserRequest)
	if err := json.NewDecoder(r.Body).Decode(deleteUserReq); err != nil {
		return err
	}
	if err := api.Store.DeleteUser(deleteUserReq.UserID); err != nil {
		return err
	}
	return shared.WriteJSON(w, http.StatusOK, "user deleted")
}

func (api *APIServer) handleCreateDogs(w http.ResponseWriter, r *http.Request) error {
	createDogReq := new(models.CreateDogRequest)
	if err := json.NewDecoder(r.Body).Decode(createDogReq); err != nil {
		return err
	}
	dog := models.NewDog(createDogReq.DogName, createDogReq.Breed, createDogReq.Location, createDogReq.ImageURL, createDogReq.Owner, createDogReq.ContactNumber)
	if err := api.Store.CreateDog(dog); err != nil {
		return err
	}
	return shared.WriteJSON(w, http.StatusOK, dog)
}

func (api *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := new(models.CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}
	user := models.NewUser(createUserReq.FirstName, createUserReq.LastName, createUserReq.MailID, createUserReq.Password, createUserReq.Username)
	if err := api.Store.CreateUser(user); err != nil {
		return err
	}
	tokenString, err := shared.CreateJwt(user)
	if err != nil {
		return err
	}
	fmt.Println(tokenString)
	return shared.WriteJSON(w, http.StatusOK, user)
}
