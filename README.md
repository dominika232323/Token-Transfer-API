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


