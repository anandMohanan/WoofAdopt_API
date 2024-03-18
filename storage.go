package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Storage interface {
	CreateDog(*Dog) error
	DeleteUser(int) error
	UpdateDog(*Dog) error
	GetDogById(int) (*Dog, error)
	CreateUser(*User) error
	GetAllUsers() ([]*User, error)
	GetUserById(int) (*User, error)
	GetUserByUsername(user_name string) (*User, error)
}

type SqlLiteStore struct {
	db *sql.DB
}

func NewStore() (*SqlLiteStore, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	url := os.Getenv("DB_URL")
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	return &SqlLiteStore{
		db,
	}, nil
}
func (s *SqlLiteStore) Init() {
	if err := s.CreateDogTable(); err != nil {
		log.Fatal(err)
	}
	if err := s.CreateUserTable(); err != nil {
		log.Fatal(err)
	}
	if err := s.CreateBreedTable(); err != nil {
		log.Fatal(err)
	}
	if err := s.CreateFavoriteTable(); err != nil {
		log.Fatal(err)
	}
}

func (s *SqlLiteStore) CreateDogTable() error {
	_, err := s.db.Exec(CreateDogTableQuery)
	return err
}

func (s *SqlLiteStore) CreateBreedTable() error {
	_, err := s.db.Exec(CreateBreedTableQuery)
	return err
}
func (s *SqlLiteStore) CreateUserTable() error {
	_, err := s.db.Exec(CreateUserTableQuery)
	return err
}
func (s *SqlLiteStore) CreateFavoriteTable() error {
	_, err := s.db.Exec(CreateFavouriteTableQuery)
	return err
}
func (s *SqlLiteStore) CreateDog(*Dog) error {
	return nil
}

func (s *SqlLiteStore) CreateUser(user *User) error {
	query, err := s.db.Prepare(`insert into user(first_name, last_name, mail_id,user_name, encrypted_password ,is_active, created_at, lastmodified_at) values(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer query.Close()
	resp, err := query.Exec(user.FirstName, user.LastName, user.MailID, user.UserName, user.EncryptedPassword, user.IsActive, user.CreatedAt, user.LastModifiedAt)
	if err != nil {
		return err
	}
	userId, err := resp.LastInsertId()

	if err != nil {
		return err
	}
	id := int(userId)
	user.ID = id
	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *SqlLiteStore) DeleteUser(user_id int) error {
	_, err := s.GetUserById(user_id)
	if err != nil {
		return err
	}
	query, err := s.db.Prepare(`update user set is_active=0 where user_id = ?`)
	if err != nil {
		return err
	}
	defer query.Close()
	resp, err := query.Exec(user_id)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}
func (s *SqlLiteStore) UpdateDog(*Dog) error {
	return nil
}
func (s *SqlLiteStore) GetDogById(int) (*Dog, error) {
	return &Dog{}, nil
}
func (s *SqlLiteStore) GetAllUsers() ([]*User, error) {

	query, err := s.db.Prepare(`select * from  user where is_active=1`)

	if err != nil {
		return nil, err
	}
	defer query.Close()
	resp, err := query.Query()
	if err != nil {
		return nil, err
	}
	users := []*User{}
	for resp.Next() {
		user := new(User)
		var createdAt, lastModifiedAt string
		if err := resp.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.MailID, &user.UserName, &user.EncryptedPassword, &user.IsActive, &createdAt, &lastModifiedAt,
		); err != nil {
			return nil, err
		}

		user.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999", strings.Split(createdAt, "+")[0])
		if err != nil {
			return nil, err
		}
		user.LastModifiedAt, err = time.Parse("2006-01-02 15:04:05.999999999", strings.Split(lastModifiedAt, "+")[0])
		if err != nil {
			return nil, err
		}
		users = append(users, user)

	}
	return users, nil

}

func (s *SqlLiteStore) GetUserById(userId int) (*User, error) {
	query := `SELECT * FROM user WHERE is_active = 1 AND user_id = ?`

	// Execute the query with the userId parameter
	row := s.db.QueryRow(query, userId)

	// Initialize a User struct to store the result
	user := &User{}

	// Scan the row into the User struct
	var createdAt, lastModifiedAt string
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.MailID, &user.UserName, &user.EncryptedPassword, &user.IsActive, &createdAt, &lastModifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999", strings.Split(createdAt, "+")[0])
	if err != nil {
		return nil, err
	}
	user.LastModifiedAt, err = time.Parse("2006-01-02 15:04:05.999999999", strings.Split(lastModifiedAt, "+")[0])
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *SqlLiteStore) GetUserByUsername(userName string) (*User, error) {
	query := `SELECT * FROM user WHERE is_active = 1 AND user_name = ?`

	// Execute the query with the userId parameter
	row := s.db.QueryRow(query, userName)

	// Initialize a User struct to store the result
	user := &User{}

	// Scan the row into the User struct
	var createdAt, lastModifiedAt string
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.MailID, &user.UserName, &user.EncryptedPassword, &user.IsActive, &createdAt, &lastModifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999", strings.Split(createdAt, "+")[0])
	if err != nil {
		return nil, err
	}
	user.LastModifiedAt, err = time.Parse("2006-01-02 15:04:05.999999999", strings.Split(lastModifiedAt, "+")[0])
	if err != nil {
		return nil, err
	}

	return user, nil
}
