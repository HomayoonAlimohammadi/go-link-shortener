.PHONY: help deps build start clean 

BINARY_NAME := linkshortener
COVERAGE_PROFILE := c.out

help:
	@echo "Link Shortener is a web app to shorten any given link to a unique phrase so you can use that any where."

build: deps 
	go build -o ${BINARY_NAME} ./cmd/linkshortener

deps: 
	go mod download

test: 
	go test ./... -coverprofile ${COVERAGE_PROFILE}

start:
	go run main.go serve

version:
	go run main.go version

clean:
	rm ${BINARY_NAME} ${COVERAGE_PROFILE}