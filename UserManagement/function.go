package main

import (
	"encoding/json"
	"fmt"
	"strconv"


	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



func connectdb() (*gorm.DB, error) {
	var err error
	datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
	db, err = gorm.Open(mysql.Open(datasource), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !checkCredentials(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized access"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkCredentials(username, password string) bool {

	var user Users
	fmt.Println("Checking credentials for:", username)

	if err := db.Where("Username = ? AND Password = ?", username, password).First(&user).Error; err != nil {
		fmt.Println("Error fetching user:", err)
		return false
	}

	fmt.Println("User found:", user.Username)
	return true
}

func getall(w http.ResponseWriter, r *http.Request) {
	var user []Users
	db.Find(&user)
	writer(w, http.StatusAccepted, user)

}
func getbyid(w http.ResponseWriter, r *http.Request) {
	var user Users
	para := mux.Vars(r)
	id, _ := strconv.Atoi(para["id"])
	if err := db.Find(&user, id).Error; err != nil {
		writer(w, http.StatusBadRequest, map[string]string{
			"message": err.Error()})
		return
	}
	writer(w, http.StatusAccepted, user)
}
func update(w http.ResponseWriter, r *http.Request) {
	
	para := mux.Vars(r)
	id, err := strconv.Atoi(para["id"])
	if err != nil {
		writer(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid ID",
		})
		return
	}


	var user Users
	if err := db.First(&user, id).Error; err != nil {
		writer(w, http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
		return
	}


	var updatedUser Users
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		writer(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
		return
	}


	db.Model(&user).Updates(updatedUser)
    writer(w, http.StatusOK, user)
}

func adduser(w http.ResponseWriter, r *http.Request) {
	var user Users
	json.NewDecoder(r.Body).Decode(&user)
	db.Save(&user)
	writer(w, http.StatusCreated, user)
}
func delete(w http.ResponseWriter,r *http.Request){
	var user Users;
	para:=mux.Vars(r)
	id,_:=strconv.Atoi(para["id"])
	if err:=db.First(&user,id).Error;err!=nil{
    writer(w,http.StatusBadRequest,map[string]string{
		"message":err.Error(),
	})
	return
	}
	db.Delete(&user)
	writer(w, http.StatusAccepted, user)

}

func writer(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)

}
