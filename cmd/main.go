package main

import (
	"accountingsystem/config"
	"accountingsystem/db"
	"accountingsystem/internal/requests/dl"
	"accountingsystem/internal/services"
	"fmt"
)

func main() {
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Error initing config: %v\n", err)
		return
	}
	if err := db.Init(); err != nil {
		fmt.Printf("Error initing database: %v\n", err)
		return
	}

	dlService := &services.DLService{}

	newDL, err := dlService.CreateDL(&dl.InsertRequest{
		Code:  "001",
		Title: "Example DL",
	})
	if err != nil {
		fmt.Printf("Error creating DL: %v\n", err)
		return
	}

	fmt.Printf("DL created successfully: %+v\n", newDL)
}
