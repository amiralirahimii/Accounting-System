package services

import (
	"accountingsystem/config"
	"accountingsystem/db"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := config.InitConfig("../../.env.test")
	if err != nil {
		log.Fatalf("Failed to load test configuration: %v", err)
	}

	err = db.Init()
	if err != nil {
		log.Fatalf("Failed to connect to the test database: %v", err)
	}

	os.Exit(m.Run())
}

func resetDatabase(t *testing.T) {
	err := db.DB.Exec("TRUNCATE TABLE dl RESTART IDENTITY CASCADE").Error
	if err != nil {
		t.Fatalf("Failed to reset database: %v", err)
	}
}

func seedDL(t *testing.T, dl *models.DL) {
	err := db.DB.Create(dl).Error
	if err != nil {
		t.Fatalf("Failed to seed DL: %v", err)
	}
}

func Test_CreateDL_WithValidRequest_Succeeds(t *testing.T) {
	// resetDatabase(t)

	req := &dl.InsertRequest{
		Code:  "DL00211",
		Title: "Test Ledger11",
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.Nil(t, err)
	assert.Equal(t, dl.Code, req.Code)
	assert.Equal(t, dl.Title, dl.Title)

	// Verify database state
	// var dbDL models.DL
	// err = db.DB.First(&dbDL, dl.ID).Error
	// if err != nil {
	// 	t.Fatalf("Failed to fetch DL from database: %v", err)
	// }

	// if dbDL.Code != req.Code || dbDL.Title != req.Title {
	// 	t.Errorf("Database DL mismatch: got %+v, expected %+v", dbDL, req)
	// }
}
