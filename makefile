run-fetcher: 
	@go run market-data-fetcher/cmd/main.go

run-ingestor: 
	@go run market-data-ingestor/cmd/main.go

run-api:
	@go run market-data-service/cmd/main.go

compose-up:
	@docker compose up

.PHONY: run-fetcher run-ingestor compose-up run-api run-all