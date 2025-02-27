# teya-home-assignment
 this is a solution for a home assignment Teya gave me as part of the interview process
 In this assignment I'm asked to implement a simple ledger system that supports adding transactions and getting the balance and transaction history.

 Important Decisions I made:
- Transactions are immutable.
- I made a separation between internal and external/ user-facing models as we don't want to expose the internal logic of the ledger to the user.
- The logic is in a separate component(pkg/ledger) to make it testable(in the context of this assignment, an integration test of the WS API routes is an overkill IMO).
- I decided to use primitive logs by printing to stdout. I decided on high verbosity logs.
- I would allow a negative balance(we Would like to charge interest for "lending" that money).
- I developed this solution while keeping in mind that the solution should be thread safe but that wasn't my focus. I do not guarantee there are no race conditions. 

Interesting discussion points for the interview:
- Test coverage and kinds of tests
- Minimalism - I did not add any properties that were not required by the assignment(e.g. date, merchant/the other as part of the transaction model)

## API Documentation

### Endpoints

#### Create Transaction
- **URL**: `/transaction`
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
- **URL**: `/transaction?offset=0&limit=10`
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






