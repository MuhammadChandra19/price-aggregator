run-fetcher: 
	@go run market-data-fetcher/cmd/main.go

run-ingestor: 
	@go run market-data-ingestor/cmd/main.go

compose-up:
	@docker compose up

.PHONY: run-fetcher run-ingestor compose-up