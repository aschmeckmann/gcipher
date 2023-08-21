package main

import (
	"fmt"
	"gcipher/cmd/userctl"
	"gcipher/internal/server"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gcipher [command]")
		fmt.Println("Available commands:")
		fmt.Println("  server - Start the server")
		fmt.Println("  userctl - User management")
		return
	}

	command := os.Args[1]
	switch command {
	case "server":
		server.StartServer()
	case "userctl":
		userctl.Execute()
	default:
		fmt.Println("Unknown command:", command)
	}
}
