# Token Transfer API

A GraphQL API written in Go for transferring BTP tokens between wallets.

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

### Running tests

```bash
docker compose up --build test
```
