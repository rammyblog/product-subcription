package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rammyblog/go-product-subscriptions/config"
	"github.com/rammyblog/go-product-subscriptions/controller"
	"github.com/rammyblog/go-product-subscriptions/response"
	"github.com/rammyblog/go-product-subscriptions/tests"
	"gopkg.in/go-playground/assert.v1"
)

func TestGetAllProducts(t *testing.T) {

	testDb, close, mock := tests.NewMockGORM()
	var response response.SuccessResponse

	defer close()

	rows := mock.NewRows([]string{"id", "name", "description", "price", "plan_code", "created_at", "updated_at"}).
		AddRow(1, "Product 1", "Product 1 description", 1000, "plan_code_1", time.Now(), time.Now()).
		AddRow(2, "Product 2", "Product 2 description", 2000, "plan_code_2", time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM \"products\"*").WillReturnRows(rows)

	// create an instance of our config with the db
	config.GlobalConfig = &config.AppConfig{
		DB: testDb,
	}

	// Create a new request to /products
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	// now we execute our method
	controller.GetAllProducts(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
		return
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(response.Data.([]interface{})), 2)

}
