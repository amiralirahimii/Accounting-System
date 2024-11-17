package services

import (
	"accountingsystem/config"
	"accountingsystem/db"
	"accountingsystem/internal/constants"
	"accountingsystem/internal/models"
	"accountingsystem/internal/requests/dl"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_CreateDL_Succeeds_WithValidRequest(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.Nil(t, err)
	assert.Equal(t, dl.Code, req.Code)
	assert.Equal(t, dl.Title, dl.Title)
}

func Test_CreateDL_ReturnsErrCodeEmptyOrTooLong_WithEmptyCode(t *testing.T) {
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  "",
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrCodeEmptyOrTooLong_WithTooLongCode(t *testing.T) {
	randomCode := generateRandomString(65)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrTitleEmptyOrTooLong_WithEmptyTitle(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: "",
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrTitleEmptyOrTooLong_WithTooLongTitle(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := generateRandomString(65)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}

	service := DLService{}
	dl, err := service.CreateDL(req)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrCodeAlreadyExists_WithExistingCode(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	validReq := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	randomTitleNotExisting := "Test" + generateRandomString(20)
	reqWithExistingCode := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitleNotExisting,
	}

	service := DLService{}
	_, err := service.CreateDL(validReq)
	require.Nil(t, err)

	dl, err := service.CreateDL(reqWithExistingCode)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeAlreadyExists)
	assert.Nil(t, dl)
}

func Test_CreateDL_ReturnsErrTitleAlreadyExists_WithExistingTitle(t *testing.T) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	validReq := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	randomCodeNotExisting := "DL" + generateRandomString(20)
	reqWithExistingTitle := &dl.InsertRequest{
		Code:  randomCodeNotExisting,
		Title: randomTitle,
	}

	service := DLService{}
	_, err := service.CreateDL(validReq)
	require.Nil(t, err)

	dl, err := service.CreateDL(reqWithExistingTitle)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleAlreadyExists)
	assert.Nil(t, dl)
}

func createRandomDL(s *DLService) (*models.DL, error) {
	randomCode := "DL" + generateRandomString(20)
	randomTitle := "Test" + generateRandomString(20)
	req := &dl.InsertRequest{
		Code:  randomCode,
		Title: randomTitle,
	}
	createDL, err := s.CreateDL(req)
	return createDL, err
}

func Test_UpdateDL_Succeeds_WithValidRequest(t *testing.T) {
	service := DLService{}
	createdDL, err := createRandomDL(&service)
	require.Nil(t, err)

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:    createdDL.ID,
		Code:  newRandomCode,
		Title: newRandomTitle,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.Nil(t, err)
	assert.Equal(t, updatedDL.Code, updateReq.Code)
	assert.Equal(t, updatedDL.Title, updateReq.Title)
}

func Test_UpdateDL_ReturnsErrCodeEmptyOrTooLong_WithEmptyCode(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:    createdDL.ID,
		Code:  "",
		Title: newRandomTitle,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrCodeEmptyOrTooLong_WithTooLongCode(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := generateRandomString(65)
	newRandomTitle := "Test" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:    createdDL.ID,
		Code:  newRandomCode,
		Title: newRandomTitle,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrCodeEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleEmptyOrTooLong_WithEmptyTitle(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := "DL" + generateRandomString(20)
	updateReq := &dl.UpdateRequest{
		ID:    createdDL.ID,
		Code:  newRandomCode,
		Title: "",
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrTitleEmptyOrTooLong_WithTooLongTitle(t *testing.T) {
	service := DLService{}
	createdDL, _ := createRandomDL(&service)

	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := generateRandomString(65)
	updateReq := &dl.UpdateRequest{
		ID:    createdDL.ID,
		Code:  newRandomCode,
		Title: newRandomTitle,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrTitleEmptyOrTooLong)
	assert.Nil(t, updatedDL)
}

func Test_UpdateDL_ReturnsErrDLNotFound_WithNonExistingID(t *testing.T) {
	service := DLService{}
	newRandomCode := "DL" + generateRandomString(20)
	newRandomTitle := "Test" + generateRandomString(20)
	newID := generateRandomInt64()

	updateReq := &dl.UpdateRequest{
		ID:    newID,
		Code:  newRandomCode,
		Title: newRandomTitle,
	}

	updatedDL, err := service.UpdateDL(updateReq)

	require.NotNil(t, err)
	assert.ErrorIs(t, err, constants.ErrDLNotFound)
	assert.Nil(t, updatedDL)
}