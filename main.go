package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserData struct {
	data map[string]string
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var revokedTokens = struct {
	data map[string]string
}{data: make(map[string]string)}

var secretKey = []byte("secret-key")

var userDB = UserData{data: make(map[string]string)}

func main() {
	http.HandleFunc("/signup", signUpUserHelper)

	http.HandleFunc("/signin", signInUserHelper)

	http.HandleFunc("/protected", ProtectedHandler)

	http.HandleFunc("/revoke", revokeJWTToken)

	http.HandleFunc("/renewToken", renewToken)

	http.ListenAndServe(":8080", nil)

}

// protected handler
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	tokenString = tokenString[len("Bearer "):]
	_, revoked := revokedTokens.data[tokenString]

	if revoked {
		http.Error(w, "Token has been revoked", http.StatusUnauthorized)
	}
	err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	fmt.Fprint(w, "Welcome to the the protected area")
}

// verify the token
func verifyToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invlid error")
	}
	return nil
}

// sign user helper
func signInUserHelper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	// If the email and password is not empty, now check if the user is already exist in the database
	password, isUserExist := userDB.data[user.Email]

	if !isUserExist {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if password != user.Password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// generate the JWT token

	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		http.Error(w, "Failes to generate token", http.StatusInternalServerError)
		return
	}
	// send token as response

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})

}

// sign up user handler
func signUpUserHelper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	userDB.data[user.Email] = string(user.Password)

	w.WriteHeader(http.StatusCreated)
	fmt.Println("User created succsssfully : ", user.Email)
}

// revoke token
// idea is to store the revoked tokens in a map and then check if token is revoked or not

func revokeJWTToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
	}

	tokenString = tokenString[len("Bearer "):]

	// now add the token string to blacklist this
	revokedTokens.data[tokenString] = ""

	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "Token has been revoked successfully")

}

// renew token handler
func renewToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
	}

	tokenString = tokenString[len("Bearer "):]

	_, revoke := revokedTokens.data[tokenString]

	if revoke {
		http.Error(w, "Revoked tokens can not be renew", http.StatusUnauthorized)
	}

	// now validate the token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	oldClaims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok || !token.Valid {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"user": oldClaims.Subject,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString(secretKey)

	if err != nil {
		http.Error(w, "Failes to generate token", http.StatusInternalServerError)
		return
	}
	// send token as response

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})

}
