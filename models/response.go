package models 

type LoginResponse struct {
	User  User
	Token string `json:"token"`
}
