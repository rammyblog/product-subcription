package main

import (
	"os"
	"testing"
)

func TestPaystack(t *testing.T) {
	// Set up environment variable
	os.Setenv("PAYSTACK_SECRET_KEY", "your_secret_key")

	// Call the Paystack function
	client := Paystack()

	// Assert that the client is not nil
	if client == nil {
		t.Errorf("Expected client to be initialized, but got nil")
		return
	}

	// Assert that the client's secret key is set correctly
	if client.APIKey != "your_secret_key" {
		t.Errorf("Expected client's secret key to be 'your_secret_key', but got '%s'", client.APIKey)
	}
}
