package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var secret = []byte("secret")

func connectdb() (*gorm.DB, error) {
	var err error
	datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
	db,err=gorm.Open(mysql.Open(datasource),&gorm.Config{})
    if err!=nil{
		return nil,err
	}
	return db,nil
}
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !Auth(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Auth(username, password string) bool {
	return username == "admin" && password == "123"
}
func GenerateToken(w http.ResponseWriter, r *http.Request) {

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
		Issuer:    "admin",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenstring)
}
func UserMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		authHeader := r.Header.Get("Authorization")
		if username != "" || password != "" {

			if !ok || !Auth(username, password) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		} else {

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				http.Error(w, "Invalid token format", http.StatusBadRequest)
				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Token expired or invalid", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})

}



func GenerateData(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Age:  22,
		name: "sarvesh",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
func getall(w http.ResponseWriter,r *http.Request){
   var user []User
   db.Find(&user)
   writer(w,http.StatusAccepted,user)
}
func writer(w http.ResponseWriter,status int,message interface{}){
	w.Header().Set("Content-Type","application-json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}