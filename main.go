package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
)

type User struct {
	ID      uint64 `validate:"required"`
	Name    string `validate:"required"`
	Email   string `validate:"required,email"`
	Age     int    `validate:"gte=18"`
	Active  bool
	Created int64
}

var users []User

func handleUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		handleGetUsers(w, r)
	} else if r.Method == http.MethodPost {
		handlePostUsers(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(users)

	if err != nil {
		fmt.Fprint(w, "Error marshaling the users data : %v", err)
	}
	w.Write(jsonData)
}

func handlePostUsers(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	validate := validator.New()

	valErr := validate.Struct(newUser)
	if valErr != nil {
		http.Error(w, "Input is invalid", http.StatusBadRequest)
		return
	}

	users = append(users, newUser)
	fmt.Println("Received user data: ", newUser)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}

func main() {
	users = append(users, User{ID: 1, Name: "John Doe", Email: "john.doe@example.com", Age: 30, Active: true, Created: time.Now().Unix()})
	users = append(users, User{ID: 2, Name: "Jane Smith", Email: "jane.smith@example.com", Age: 25, Active: false, Created: time.Now().Unix()})
	http.HandleFunc("/users", handleUsers)

	// start the http server on port 3000
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
