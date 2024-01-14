package controller

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rammyblog/go-product-subscriptions/config"
	"github.com/rammyblog/go-product-subscriptions/models"
)

type PaystackEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func TransactionWebhook(w http.ResponseWriter, r *http.Request) {
	// Retrieve the request's body
	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Secret key for HMAC
	secret := []byte(os.Getenv("PAYSTACK_SECRET_KEY"))

	h := hmac.New(sha512.New, []byte(secret))
	h.Write(reqBody)
	computedHash := hex.EncodeToString(h.Sum(nil))

	signature := r.Header.Get("X-Paystack-Signature")

	if computedHash == signature {
		var event PaystackEvent
		json.Unmarshal(reqBody, &event)
		fmt.Println(event.Event)
		if event.Event == "charge.success" {
			email := event.Data.(map[string]interface{})["customer"].(map[string]interface{})["email"].(string)
			customerCode := event.Data.(map[string]interface{})["customer"].(map[string]interface{})["customer_code"].(string)
			fmt.Println(email, customerCode)
			if err := config.GlobalConfig.DB.Model(&models.User{}).Where("email = ?", email).Update("CustomerCode", customerCode).Error; err != nil {
				w.WriteHeader(http.StatusBadGateway)
				log.Fatal(err)
				return
			}
			w.WriteHeader(http.StatusOK)

		}
	}
	w.WriteHeader(http.StatusOK)
}
