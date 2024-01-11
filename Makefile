build:
	go build -o bin/product-subcription

run: build
	./bin/product-subcription

seed: build
	 ./bin/product-subcription --seed


test: 
	go test -v ./...