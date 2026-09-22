# Portfolio Tracker API

Backend service for managing investment portfolios, assets, and transactions.

The service provides:
- user registration and authentication;
- asset management;
- buy and sell transactions;
- portfolio calculation;
- real-time market prices;
- JWT-based authentication;
- PostgreSQL persistence;
- concurrent market data fetching.



## Tech Stack

- Go
- PostgreSQL
- Docker / Docker Compose
- JWT
- bcrypt
- REST API
- Binance API
- pgx
- standard library `net/http`



## Features

### Authentication
- User registration
- User login
- JWT authentication
- Password hashing with bcrypt
- Protected endpoints

### Assets
- List available assets
- Get asset by ID

### Transactions
- Create BUY and SELL transactions
- Get user's transactions
- Get transaction by ID
- Transaction validation
- Prevent selling more assets than available

### Portfolio
- Calculate current positions
- Calculate average purchase price
- Calculate current position value
- Calculate unrealized profit/loss
- Calculate realized profit/loss
- Calculate profit/loss percentage
- Fetch current market prices concurrently

### Infrastructure
- PostgreSQL persistence
- Database migrations
- Docker / Docker Compose
- Graceful HTTP server shutdown
- Context cancellation
- Unit and HTTP tests



## Development

The project is currently under development.

For development and testing purposes, the application can be run using Docker Compose.

## Local Development

### Requirements

- Go 1.27+
- Docker
- Docker Compose

### Configuration

Create a `.env` file in the project root:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=portfolio_tracker
JWT_SECRET=dev-secret
```

### Start

```bash
docker compose up -d --build
```

The API will be available at:
http://localhost:8080

Health check:
GET /health