package main
// Define User model
type User struct {
    ID       uint `gorm:"primaryKey"`
    Username string

}

// Define Account model
type Account struct {
    ID       uint `gorm:"primaryKey"`
    UserID   uint // Foreign key referencing User ID
    Balance  float64
    // Define other account fields
}

type Data struct {
	Age  int
	name string
}


