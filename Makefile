# Agnes AI Studio — backend Makefile

BINARY ?= server
ADDR   ?= :8080
DATA   ?= ./data
ROUTE  ?= international

.PHONY: build run dev test clean tidy

build:
	cd backend && go build -o $(BINARY).exe ./cmd/server

run: build
	cd backend && ./$(BINARY).exe -addr $(ADDR) -data $(DATA) -route $(ROUTE)

dev:
	cd backend && go run ./cmd/server -addr $(ADDR) -data $(DATA) -route $(ROUTE)

test:
	cd backend && go test ./...

tidy:
	cd backend && go mod tidy

clean:
	rm -f backend/$(BINARY) backend/$(BINARY).exe