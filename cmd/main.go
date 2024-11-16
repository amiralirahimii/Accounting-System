package main

import (
	"accountingsystem/config"
	"accountingsystem/db"
	"accountingsystem/internal/requests/dl"
	"accountingsystem/internal/services"
	"fmt"
	"log"
)

func main() {
	if err := config.InitConfig(".env"); err != nil {
		log.Fatalf("Error initing config: %v\n", err)
		return
	}
	if err := db.Init(); err != nil {
		log.Fatalf("Error initing database: %v\n", err)
		return
	}

	dlService := &services.DLService{}

	// newDL, err := dlService.CreateDL(&dl.InsertRequest{
	// 	Code:  "001",
	// 	Title: "Example DL",
	// })
	// if err != nil {
	// 	fmt.Printf("Error creating DL: %v\n", err)
	// }

	newDL1, err1 := dlService.UpdateDL(&dl.UpdateRequest{
		ID:    1,
		Code:  "002",
		Title: "Example DL 1",
	})
	if err1 != nil {
		fmt.Printf("Error creating DL: %v\n", err1)
	}

	fmt.Printf("DL created successfully: %+v\n", newDL1)
}
