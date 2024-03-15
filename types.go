package main

import (
	"time"
)

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
	MailID    string `json:"mail_id"`
}

type Dog struct {
	ID             int       `json:"id"`
	DogName        string    `json:"dog_name"`
	DogBreed       string    `json:"breed"`
	Location       string    `json:"location"`
	ImageURL       string    `json:"image_url"`
	ContactNumber  int64     `json:"contact_number"`
	Owner          int       `json:"owner_id"`
	IsActive       int       `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"lastmodified_at"`
}

type User struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	MailID         string    `json:"mail_id"`
	IsActive       int       `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"lastmodified_at"`
}

type Favorite struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	DogID  int `json:"dog_id"`
}
type Breed struct {
	ID        int    `json:"id"`
	BreedName string `json:"breed_name"`
}

func NewDog(dogName, breed, location, imageUrl string, owner int, contactNumber int64) *Dog {
	return &Dog{
		DogName:        dogName,
		DogBreed:       breed,
		Location:       location,
		ImageURL:       imageUrl,
		ContactNumber:  contactNumber,
		IsActive:       1,
		Owner:          owner,
		CreatedAt:      time.Now().UTC(),
		LastModifiedAt: time.Now().UTC(),
	}
}

func NewUser(userFirstName, userLastName, userMailId string) *User {
	return &User{
		FirstName:      userFirstName,
		LastName:       userLastName,
		MailID:         userMailId,
		IsActive:       1,
		CreatedAt:      time.Now().UTC(),
		LastModifiedAt: time.Now().UTC(),
	}
}
