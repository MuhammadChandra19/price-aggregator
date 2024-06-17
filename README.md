# Microservices for Real-Time Market Data
This repository contains microservices for fetching, ingesting, and serving real-time market data using Go and Docker Compose.


## Prerequisites

- Go installed (version 1.18+ recommended)
- Docker and Docker Compose installed


## Running the Services
### 1. Start Docker Compose
To spin up necessary infrastructure components (like Kafka, Redis, etc.):

```bash
  make compose-up

```

### 2.  Run the Ingestor
Start the Kafka consumer to ingest real-time data:
```bash
  make run-ingestor

```

### 3.  Start the API
Launch the API service to serve the stored market data:
```bash
 make run-api

```

### 4.  Fetch Real-Time Data
Begin fetching real-time data from Binance and DeGate:
```bash
 make run-fetcher

```

## Testing the API
Once the services are up and running, you can test the API using Postman:

Open Postman: Launch the Postman application or use the Postman web interface.

Send a GET Request:

#### URL: http://localhost:8080/market/ETH
#### Method: GET


# simple-queue-channel
