package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rammyblog/go-product-subscriptions/config"
	"github.com/rammyblog/go-product-subscriptions/database"
	"github.com/rammyblog/go-product-subscriptions/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	seed := flag.Bool("seed", false, "seed the db")

	flag.Parse()
	paystackClient := Paystack()
	db, err := database.Init(seed)

	if err != nil {
		log.Fatal("Could not connect to db")
		panic(err)
	}
	fmt.Println("db conected", db.Name())

	config.GlobalConfig = &config.AppConfig{
		DB:             db,
		PaystackClient: paystackClient,
	}

	if seed != nil && *seed {
		database.Seed(db)
	}

	port := fmt.Sprintf(":%v", os.Getenv("PORT"))

	handler := router.Init()

	server := &http.Server{
		Addr:    port,
		Handler: handler,
	}
	log.Printf("[info] start http server listening %s", port)
	server.ListenAndServe()
}
