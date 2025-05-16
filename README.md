# Token Transfer API

Token Transfer API is a GraphQL service written in Go that enables secure BTP token transfers between wallets.
It ensures consistent and safe operations using PostgreSQL transactions and row-level locking.

The system is designed to be simple, robust, and testable, featuring:
- transferring tokens between wallet addresses using a GraphQL mutation,
- GraphQL Playground for interactive testing,
- fully containerized development via Docker Compose,
- built-in tests to validate core functionality and race condition safety.

## Technologies Used

- **Go** (Golang)
- **GraphQL** (gqlgen)
- **GORM** (ORM for PostgreSQL)
- **PostgreSQL**
- **Testify**
- **Docker & Docker Compose**

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/dominika232323/Token-Transfer-API.git
cd Token-Transfer-API
```

### Environment Configuration

Before running the application, set up your environment variables.

#### Create an .env file in the root directory

```bash
touch .env
```

#### Fill in the following variables

```
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database
POSTGRES_PORT=your_port
POSTGRES_HOST=db
```

`POSTGRES_HOST=db` is used when running the app or tests via Docker Compose.
If you're running locally without Docker Compose, you can set `POSTGRES_HOST=localhost` instead.

### Running the Application

```bash
docker compose up --build app
```

Visit http://localhost:8080/ to use GraphQL Playground.

### Running tests

```bash
docker compose up --build test
```

This includes:
- Tests for successful and unsuccessful transfers
- Concurrency tests to simulate race conditions

## Example GraphQL Mutations

Use the `transfer` mutation to send tokens between wallets. Below are examples of both successful and failed transfer scenarios.

You can try them in the GraphQL Playground at http://localhost:8080

### Successful transaction

```
mutation {
  transfer(
    from_address: "0x0000000000000000000000000000000000000000",
    to_address: "0x0000000000000000000000000000000000000001",
    amount: 200
  )
}
```

Returns

```
{
  "data": {
    "transfer": 999800
  }
}
```

### Insufficient balance

```
mutation {
  transfer(
    from_address: "0x0000000000000000000000000000000000000000",
    to_address: "0x0000000000000000000000000000000000000001",
    amount: 2000000
  )
}
```

Returns 

```
{
  "errors": [
    {
      "message": "Insufficient balance",
      "path": [
        "transfer"
      ]
    }
  ],
  "data": null
}
```

**Note:** If a wallet with the address `0x0000000000000000000000000000000000000000` balance at leat 2000000 exists in the database, this transaction will succeed.

### Sender not found

```
mutation {
  transfer(
    from_address: "0x0000000000000000000000000000000000000002",
    to_address: "0x0000000000000000000000000000000000000001",
    amount: 200
  )
}
```

Returns 

```
{
  "errors": [
    {
      "message": "sender not found",
      "path": [
        "transfer"
      ]
    }
  ],
  "data": null
}
```

**Note:** If a wallet with the address `0x0000000000000000000000000000000000000002` and balance at leat 200 exists in the database, this transaction will succeed.

### Recipient not found

```
mutation {
  transfer(
    from_address: "0x0000000000000000000000000000000000000000",
    to_address: "0x0000000000000000000000000000000000000003",
    amount: 200
  )
}
```

Returns 

```
{
  "errors": [
    {
      "message": "recipient not found",
      "path": [
        "transfer"
      ]
    }
  ],
  "data": null
}
```

**Note:** If a wallet with the address `0x0000000000000000000000000000000000000003` exists in the database, this transaction will succeed.

### Amount cannot be negative

```
mutation {
  transfer(
    from_address: "0x0000000000000000000000000000000000000000",
    to_address: "0x0000000000000000000000000000000000000001",
    amount: -200
  )
}
```

Returns

```
{
  "errors": [
    {
      "message": "amount cannot be negative",
      "path": [
        "transfer"
      ]
    }
  ],
  "data": null
}
```