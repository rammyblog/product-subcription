package main

import (
	"os"

	"github.com/rammyblog/go-paystack"
)

func Paystack() *paystack.Client {

	PAYSTACK_SECRET_KEY := os.Getenv("PAYSTACK_SECRET_KEY")

	if PAYSTACK_SECRET_KEY == "" {
		panic("PAYSTACK_SECRET_KEY is not set")
	}

	client := paystack.NewClient(PAYSTACK_SECRET_KEY)

	return client

}
