package userctl

import (
	"fmt"
	"gcipher/internal/db/models"
	"gcipher/internal/db/repositories"
	"gcipher/internal/util"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: gcipher userctl register [username] [password]")
		return
	}

	username := os.Args[3]
	password := os.Args[4]

	userRepo, err := repositories.NewUserRepository()
	if err != nil {
		fmt.Println("Failed to initialize user repository:", err)
		return
	}

	// Check if the username already exists
	existingUser, err := userRepo.FindByUsername(username)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println("Error checking username:", err)
		return
	}

	if existingUser != nil {
		fmt.Println("Username already exists.")
		return
	}

	// Hash the password
	hashedPassword, err := util.GenerateFromPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	// Create and insert the user record
	newUser := models.User{
		Username: username,
		Password: hashedPassword,
	}

	err = userRepo.Insert(newUser)
	if err != nil {
		fmt.Println("Error registering user:", err)
		return
	}

	fmt.Println("User registered successfully.")
}
