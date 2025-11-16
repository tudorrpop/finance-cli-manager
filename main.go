/*
Copyright © 2025 tuuudoor
*/
package main

import (
	"fmt"
	"log"

	"finance-cli-manager/cmd"
	"finance-cli-manager/internal/db"
)

func main() {
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer database.Close()

	fmt.Println("Connected to MySQL successfully!")

	cmd.Execute()
}
