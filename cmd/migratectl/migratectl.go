package migratectl

import (
	"fmt" // import the package containing the MigrateCerts function
	"os"
)

func Execute() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gcipher migratectl [command]")
		fmt.Println("Available commands:")
		fmt.Println("  migrate-certs [path-to-certs-directory] [username] - Migrate certificates to the database")
		return
	}

	subcommand := os.Args[2]
	switch subcommand {
	case "migrate-certs":
		if len(os.Args) < 5 {
			fmt.Println("Usage: gcipher migratectl migrate-certs [path-to-certs-directory] [username]")
			return
		}

		certsPath := os.Args[3]
		username := os.Args[4]

		if err := MigrateCerts(certsPath, username); err != nil {
			fmt.Println("Migration failed:", err)
		} else {
			fmt.Println("Migration succeeded.")
		}

	default:
		fmt.Println("Unknown subcommand:", subcommand)
	}
}
