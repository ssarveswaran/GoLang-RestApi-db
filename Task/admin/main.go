package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)
const(
	username="root"
	password="password"
	hostname="localhost:3306"
	dbname="user"
)
var db *gorm.DB
func main() {
	var err error;
	db,err=connectdb();
	if err!=nil{
		fmt.Println("Failed to Connect to the db");
		return
	}
	err=db.AutoMigrate(&User{},&Account{})
	if err != nil {
        fmt.Println("Migration failed:", err)
        return
    }
	router := mux.NewRouter()
    router.HandleFunc("/generate-token", Middleware(GenerateToken)).Methods("GET")
	// router.HandleFunc("/generate-data",data).Methods("GET")
	router.HandleFunc("/data",UserMiddleWare(GenerateData)).Methods("GET")
	router.HandleFunc("/getall",UserMiddleWare(getall)).Methods("GET")
	http.ListenAndServe(":8090", router)
	
}

