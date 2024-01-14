package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rammyblog/go-product-subscriptions/config"
	"github.com/rammyblog/go-product-subscriptions/controller"
	"github.com/rammyblog/go-product-subscriptions/response"
	"github.com/rammyblog/go-product-subscriptions/tests"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	// Create a request body with user data
	userData := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}

	var response response.SuccessResponse

	testDb, close, mock := tests.NewMockGORM()

	defer close()
	// create an instance of our config with the db
	config.GlobalConfig = &config.AppConfig{
		DB: testDb,
	}

	reqBody, _ := json.Marshal(userData)

	rows := mock.NewRows([]string{"id", "email", "customer_code", "created_at", "updated_at"}).
		AddRow(0, "", "", time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"*").WillReturnRows(rows)

	// Create a new HTTP request with the request body
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"users\" (.+) VALUES (.+) RETURNING \"id\"").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Test User", "test@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	
	// Call the CreateUser handler function
	controller.CreateUser(rr, req)

	t.Log(rr.Body.String())

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Add assertions to validate the response body
	assert.Equal(t, "test@example.com", response.Data.(map[string]interface{})["email"])
	assert.Equal(t, "", response.Data.(map[string]interface{})["customer_code"])
	assert.Equal(t, 1, int(response.Data.(map[string]interface{})["id"].(float64)))

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Log("Error with expectations:", err)
	}

}
