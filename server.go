package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

// All API responses
type Response struct {
	Data   interface{}
	Error  string
	Status int
}

// Request Body for Login/Signup
type AuthBody struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

// User Struct
type User struct {
	Id        string
	Email     string
	Password  string
	OtherData interface{}
}

// Token wrapper for Response struct Data field
type TokenResponse struct {
	Token string
}

// User wrapper for Response struct Data field
type UserResponse struct {
	User User
}

// Users wrapper for Response struct Data field
type UsersResponse struct {
	Users userDB
}

// Mock Database
// Actual User slice
type userDB []User

var (
	// Error
	err error
	// Initialize response
	result Response

	// Initialize Database
	users userDB

	// Map of authorized tokens
	tokenMap = map[string]string{}

	// Possible Characters in random token/id
	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// HTTP response codes
	success_status        = 200
	not_found_status      = 404
	conflict_status       = 409
	wrong_password_status = 400
	server_error_status   = 500
	unauthorized_status   = 401
)

// Serves Users
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(http.CanonicalHeaderKey("Authorization"))
	authArray := strings.Split(authHeader, " ")

	if len(authArray) < 2 {
		result = Response{
			Data:   nil,
			Error:  "Unauthorized",
			Status: unauthorized_status,
		}
	} else if _, ok := tokenMap[authArray[1]]; ok {
		// Get users from Database HERE

		// Create API response
		result = Response{
			Data:   UsersResponse{Users: users},
			Error:  "",
			Status: success_status,
		}
	} else {
		// Create API response
		result = Response{
			Data:   nil,
			Error:  "Unauthorized",
			Status: unauthorized_status,
		}
	}

	serveJSON(result, w, r)
}

// Serves User object
func UserHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(http.CanonicalHeaderKey("Authorization"))
	authArray := strings.Split(authHeader, " ")

	if len(authArray) < 2 {
		result = Response{
			Data:   nil,
			Error:  "Unauthorized",
			Status: unauthorized_status,
		}
	} else if _, ok := tokenMap[authArray[1]]; ok {
		// Grab the users id from the incoming url
		vars := mux.Vars(r)
		id := vars["id"]

		// Get user with id from database HERE
		user, err := users.getById(id)

		if err != nil {
			result = Response{Data: nil, Error: "user with that email does not exist", Status: not_found_status}
		} else {
			// Create API response
			result = Response{
				Data:   UserResponse{User: user},
				Error:  "",
				Status: success_status,
			}
		}
	} else {
		// Create API response
		result = Response{
			Data:   nil,
			Error:  "Unauthorized",
			Status: unauthorized_status,
		}
	}

	serveJSON(result, w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming json
	var body AuthBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	// Get user from database
	// & check password
	// Handle and errors
	user, err := users.getByEmail(body.Email)
	if err != nil {
		result = Response{Data: nil, Error: "user with that email does not exist", Status: not_found_status}
	} else {
		if user.Password != body.Password {
			result = Response{Data: nil, Error: "", Status: wrong_password_status}
		} else {
			// Make new token with super random func
			token := randSeq(10)
			tokenMap[token] = user.Id
			result = Response{Data: TokenResponse{Token: token}, Error: "", Status: success_status}
		}
	}

	// Log error if there is one
	if result.Error != "" {
		log.Printf("Error: %v", result.Error)
	}

	serveJSON(result, w, r)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming json
	var body AuthBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	_, err := users.getByEmail(body.Email)

	if err != nil {
		// Create new user
		user := User{
			randSeq(10),
			body.Email,
			body.Password,
			"User Data",
		}

		// Add to database here
		users = append(users, user)

		// Return new user
		result = Response{
			Data:   user,
			Error:  "",
			Status: success_status,
		}
	} else {
		// Return new user
		result = Response{
			Data:   nil,
			Error:  "user with that email already exists",
			Status: conflict_status,
		}
	}

	serveJSON(result, w, r)
}

func serveJSON(result Response, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(result)

	if err != nil {
		result = Response{
			Data:   nil,
			Error:  "json error",
			Status: server_error_status,
		}
	}

	j, err = json.Marshal(result)
	if err != nil {
		panic(err)
	}

	w.Write(j)
}

// Get User by id
// Function on Mock DB
func (db userDB) getById(id string) (User, error) {
	for _, user := range db {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, errors.New("no user in database with that id")
}

// Get User by email
// Function on Mock DB
func (db userDB) getByEmail(email string) (User, error) {
	for _, user := range db {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("no user in database with that email")
}

// Generates random sequence
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func main() {
	log.Println("Starting Server...")

	r := mux.NewRouter()

	// Get User
	r.HandleFunc("/api/users", UsersHandler).Methods("GET")
	r.HandleFunc("/api/users/{id}", UserHandler).Methods("GET")

	// Login Handler
	r.HandleFunc("/api/auth/login", LoginHandler).Methods("POST")

	// Signup Handler
	r.HandleFunc("/api/auth/signup", SignupHandler).Methods("POST")

	http.Handle("/api/", r)

	http.Handle("/", http.FileServer(http.Dir("./client/")))

	// Start Database session here
	users = make(userDB, 5)

	// Startup Web app
	log.Println("Listening on 1337")
	http.ListenAndServe(":1337", nil)
}
