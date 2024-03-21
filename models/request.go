package models 

type LoginRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}
type CreateDogRequest struct {
    DogName       string `json:"dog_name"`
    Breed         string `json:"breed"`
    Location      string `json:"location"`
    ImageURL      string `json:"image_url"`
    ContactNumber int64  `json:"contact_number"`
    Owner         int    `json:"owner_id"`
}
type CreateUserRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Password  string `json:"password"`
    Username  string `json:"user_name"`
    MailID    string `json:"mail_id"`
}
type DeleteUserRequest struct {
    UserID int `json:"user_id"`
}
