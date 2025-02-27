# teya-home-assignment
# Overview
This is a solution for a home assignment from Teya, implementing a simple ledger system that supports:
- Adding transactions
- Getting account balance 
- Retrieving transaction history

## Key Design Decisions

### Architecture
- **Immutable Transactions**: Once created, transactions cannot be modified
- **Separation of Models**: Internal and external/user-facing models are separated to encapsulate ledger logic
- **Modular Design**: Core logic is in a separate component (`pkg/ledger`) to enable isolated testing
- **In-Memory Storage**: Transactions are stored in memory using a slice property of the ledger component (`pkg/ledger`)

- **Thread Safety**: Basic thread safety considerations, though not fully guaranteed

### Technical Choices
- **Logging**: Simple stdout logging with high verbosity for debugging
- **Negative Balances**: Allowed to support potential interest charging on overdrafts
- **Minimal Implementation**: Focused on core requirements without additional fields like dates or merchant info

## Intresting discussion Topics for the interview
- Test coverage and kinds of tests
- Minimalism - I did not add any properties that were not required by the assignment(e.g. date, merchant/the other as part of the transaction model)


## API Documentation

### Endpoints

#### Create Transaction
- **URL**: `/api/v1/transaction`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "amount": "10.50"
  }
  ```
- **Response**:
  - Status: 201 Created (Success)
  - Status: 400 Bad Request (Invalid request body)
  - Status: 500 Internal Server Error (Server error)

#### Get Transaction History
- **URL**: `/api/v1/transaction?offset=0&limit=10`
- **Method**: `GET`
- **Query Parameters**:
  - `offset` (required): Starting position for pagination (must be >= 0)
  - `limit` (optional): Number of transactions to return (default: 10, max: 100)
- **Response**:
  - Status: 200 OK
    ```json
    {
      "transactions": [
        {
          "id": "1",
          "amount": "10.50"
        }
      ],
      "pagination": {
        "offset": 0,
        "limit": 10
      }
    }
    ```
  - Status: 400 Bad Request (Invalid query parameters)
  - Status: 500 Internal Server Error (Server error)

#### Get Account Balance
- **URL**: `/account`
- **Method**: `GET`
- **Response**:
  - Status: 200 OK
    ```json
    {
      "balance": "42.75"
    }
    ```
  - Status: 500 Internal Server Error (Server error)




# Running the Project

You can run this project either using the provided Makefile commands or Docker Compose. Both methods are explained below.

## Using Makefile

The Makefile provides several commands to build, run, and test the application:

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Stop the webserver container
make stop
```

## Using Docker Compose

Docker Compose allows you to run the application in a containerized environment:

1. Make sure you have Docker and Docker Compose installed on your system.

2. Build and start the container:
   ```bash
   docker-compose up --build
   ```

3. To run in detached mode:
   ```bash
   docker-compose up -d
   ```

4. To stop the container:
   ```bash
   docker-compose down
   ```

## API Access

Once the application is running, you can access the API endpoints as described in the API Documentation section:

- Create Transaction: `POST /api/v1/transaction`
- Get Transaction History: `GET /api/v1/transaction?offset=0&limit=10`
- Get Account Balance: `GET /account`

The API will be available at `http://localhost:8000` by default.

# API Usage Examples with cURL

Here are some example cURL commands to interact with the ledger API:

## Create a Transaction

```bash
# Add a positive transaction (deposit)
curl -X POST http://localhost:8000/api/v1/transaction \
  -H "Content-Type: application/json" \
  -d '{"amount": "25.50"}'

# Add a negative transaction (withdrawal)
curl -X POST http://localhost:8000/api/v1/transaction \
  -H "Content-Type: application/json" \
  -d '{"amount": "-10.75"}'
```

## Get Account Balance

```bash
# Retrieve the current account balance
curl -X GET http://localhost:8000/api/v1/account
```

## Get Transaction History

```bash
# Get the first 10 transactions (default limit)
curl -X GET "http://localhost:8000/api/v1/transaction?offset=0"

# Get 5 transactions starting from position 10
curl -X GET "http://localhost:8000/api/v1/transaction?offset=10&limit=5"

# Get the most recent transactions (assuming transactions are ordered by recency)
curl -X GET "http://localhost:8000/api/v1/transaction?offset=0&limit=20"
```