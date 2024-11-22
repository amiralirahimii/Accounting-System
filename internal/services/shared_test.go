package services

import (
	"accountingsystem/config"
	"accountingsystem/db"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"accountingsystem/internal/requests/sl"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

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

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomInt64() int {
	return int(seededRand.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

func createRandomSL(s *SLService, hasDL bool) (*models.SL, error) {
	randomCode := "SL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &sl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
		HasDL: hasDL,
	}
	return s.CreateSL(req)
}

func createRandomDL(s *DLService) (*dtos.DLDto, error) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	createDL, err := s.CreateDL(req)
	return createDL, err
}
