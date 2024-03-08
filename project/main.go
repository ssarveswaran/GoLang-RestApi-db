package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	username = "root"
	password = "password"
	hostname = "localhost:3306"
	dbname   = "user"
)

var db *sql.DB // Declare db variable to store the database connection globally

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func main() {
	var err error
	// Initialize the database connection
	db, err = connectDB()
	if err != nil {
		fmt.Println("Error occurred in connection db:", err)
		return
	}
	defer db.Close()

	// Create a new router instance
	router := mux.NewRouter()

	// Define the routes
router.HandleFunc("/getall", getAllHandler).Methods("GET")
router.HandleFunc("/insert", insertUserHandler).Methods("POST")
router.HandleFunc("/update/{id}", updateUserHandler).Methods("PUT")
router.HandleFunc("/delete/{id}", deleteUserHandler).Methods("DELETE")


	// Start the HTTP server
	fmt.Println("Server is running on port 9099")
	http.ListenAndServe("localhost:9991", router)
}

// connectDB connects to the MySQL database and returns a database connection
func connectDB() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
// updateUserHandler updates an existing user in the database
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from request URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Parse user ID into integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		writer(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Decode JSON request body into a map
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		writer(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	// Prepare SQL statement for updating user
	query := "UPDATE Users SET "
	var params []interface{}

	for key, value := range updates {
		// Append each field to the SQL query
		query += fmt.Sprintf("%s=?, ", key)
		params = append(params, value)
	}

	// Remove the trailing comma and space
	query = strings.TrimSuffix(query, ", ")

	// Append the WHERE clause to the SQL query
	query += " WHERE ID = ?"
	params = append(params, id)

	// Prepare and execute SQL statement
	stmt, err := db.Prepare(query)
	if err != nil {
		writer(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		writer(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	writer(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}
// deleteUserHandler deletes an existing user from the database
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from request URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Parse user ID into integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		writer(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Prepare SQL statement for deleting user
	stmt, err := db.Prepare("DELETE FROM Users WHERE ID = ?")
	if err != nil {
		writer(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(id)
	if err != nil {
		writer(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	writer(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// getAllHandler retrieves all users from the database and returns them as JSON
func getAllHandler(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers()
	if err != nil {
		writer(w, http.StatusInternalServerError, err.Error())
		return
	}
	writer(w, http.StatusOK, users)
}

// getUsers retrieves all users from the database
func getUsers() ([]User, error) {
	rows, err := db.Query("SELECT ID,Name,Age FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// insertUserHandler inserts a new user into the database
func insertUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		writer(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}else if newUser.Age ==0&&newUser.Name==""{
		writer(w, http.StatusBadRequest, map[string]string{"name": "Name field is required","Age": "Age field is required"})
		return
	}else if newUser.Name == "" {
		writer(w, http.StatusBadRequest, map[string]string{"Name": "Name field is required"})
		return
	}else if newUser.Age == 0 {
		writer(w, http.StatusBadRequest, map[string]string{"Age": "Age field is required"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO Users(Name, Age) VALUES(?, ?)")
	if err != nil {
		writer(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUser.Name, newUser.Age)
	if err != nil {
		writer(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Get the last inserted ID
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		writer(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Update the newUser's ID field with the last inserted ID
	newUser.ID = uint(lastInsertedID)

	writer(w, http.StatusCreated, newUser)
}





// writer writes the HTTP response with the specified status code and message
func writer(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}
