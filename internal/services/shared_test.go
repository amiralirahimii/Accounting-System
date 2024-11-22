package services

import (
	"accountingsystem/config"
	"accountingsystem/db"
	"accountingsystem/internal/dtos"
	"accountingsystem/internal/requests/dl"
	"accountingsystem/internal/requests/sl"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

var dlService DLService
var slService SLService
var voucherService VoucherService

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

func createRandomSL(hasDL bool) (*dtos.SLDto, error) {
	randomCode := "SL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &sl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
		HasDL: hasDL,
	}
	return slService.CreateSL(req)
}

func createRandomDL() (*dtos.DLDto, error) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	return dlService.CreateDL(req)
}
