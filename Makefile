include .env
export

run:
	API_BASE_URL=${API_BASE_URL} go run main.go

.PHONY: run
