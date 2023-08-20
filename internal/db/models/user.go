package models

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// Create a new user instance
func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}
