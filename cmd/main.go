package main

import (
	"accountingsystem/config"
	"accountingsystem/internal/services"
	"fmt"
)

func main() {
	// Initialize the database connection
	config.InitDB()

	// Instantiate the service
	dlService := &services.DLService{}

	// Example: Creating a new DL
	newDL, err := dlService.CreateDL(&services.CreateDLRequest{
		Code:  "001",
		Title: "Example DL",
	})
	if err != nil {
		fmt.Printf("Error creating DL: %v\n", err)
		return
	}

	fmt.Printf("DL created successfully: %+v\n", newDL)
}
