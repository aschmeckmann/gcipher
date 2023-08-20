package userctl

import (
	"fmt"
	"os"
)

func Execute() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gcipher userctl [command]")
		fmt.Println("Available commands:")
		fmt.Println("  register - Register a new user")
		return
	}

	subcommand := os.Args[2]
	switch subcommand {
	case "register":
		RegisterUser()
	default:
		fmt.Println("Unknown subcommand:", subcommand)
	}
}
