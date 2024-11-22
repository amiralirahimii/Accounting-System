package main

import (
	"accountingsystem/config"
	"accountingsystem/db"
	"accountingsystem/internal/services"
	"log"
)

func main() {
	if err := config.InitConfig(".env"); err != nil {
		log.Fatalf("Error initing config: %v\n", err)
		return
	}
	theDB, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return
	}

	dlService := &services.DLService{}
	slService := &services.SLService{}
	voucherService := &services.VoucherService{}

	dlService.InitService(theDB)
	slService.InitService(theDB)
	voucherService.InitService(theDB)

	log.Println("Successfully brought up the services")

	// Here we can start interacting with services as needed
}
