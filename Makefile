build:
	go build -o bin/product-subsctiption

run: build
	./bin/product-subsctiption

seed: build
	 .bin/product-subsctiption --seed


test: 
	go test -v ./...