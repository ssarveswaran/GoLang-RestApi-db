package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    
    "gorm.io/gorm"
)
var db *gorm.DB

const (
    username = "root"
    password = "password"
   hostname = "host.docker.internal:3306"
    // hostname = "db:3306"
    dbname   = "user"
)
// type User struct {
//     // gorm.Model
//     Name  string
//     Email string
// }


func main() {
    var err error
    db, err = connectdb()
    if err != nil {
        fmt.Println("Failed to connect to the database:", err)
        return
    }
   

    fmt.Println("Connection Successfully")

    // Auto migrate the User model
    err = db.AutoMigrate(&Users{})
    if err != nil {
        fmt.Println("Migration failed:", err)
        return
    }

    router := mux.NewRouter()
    router.Use(Middleware)

    router.HandleFunc("/getall", getall).Methods("GET")
    router.HandleFunc("/get/{id}", getbyid).Methods("GET")
    router.HandleFunc("/add", adduser).Methods("POST")
    router.HandleFunc("/update/{id}", update).Methods("PUT")
    router.HandleFunc("/delete/{id}", delete).Methods("DELETE")
    // router.HandleFunc("/deleteall", deleteall).Methods("DELETE")
    // http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    //     fmt.Fprintf(w, "Hello, World!\n")
    // })

    fmt.Println("Server listening on port 8080")
    http.ListenAndServe("0.0.0.0:8080", router)
}