package shared

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/anandMohanan/WoofAdopt_API/models"
	"github.com/anandMohanan/WoofAdopt_API/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func ValidateJwt(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	fmt.Println(secret)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
}
func CreateJwt(user *models.User) (string, error) {
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
func WithJWTAuth(handlerFunc http.HandlerFunc, api storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("with jwt middleware")
		tokenString := r.Header.Get("x-jwt-token")
		token, err := ValidateJwt(tokenString)
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
