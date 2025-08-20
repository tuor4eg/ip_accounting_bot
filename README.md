# IP Accounting Bot

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://telegram.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)

**IP Accounting Bot** is a Telegram bot for sole proprietors that helps track income and reminds about tax payments.

## Features (MVP)
- Income tracking for simplified taxation (6% rate)
- Automatic tax calculation
- Payment deadline reminders
- CSV data export

## Tech Stack
- Go
- PostgreSQL / SQLite (for initial setup)
- Redis (cache and rate limiting)
- Telegram Bot API

## Installation & Run

### 1) Clone the repository
```bash
git clone https://github.com/your-username/ip_accounting_bot.git
cd ip_accounting_bot
```

### 2) Install dependencies
Make sure you have **Go 1.21+** installed.
```bash
go mod tidy
```

### 3) Set environment variables
Create a `.env` file in the project root:
```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
DATABASE_URL=postgres://user:password@localhost:5432/ip_accounting
REDIS_URL=redis://localhost:6379
```

### 4) Run the bot
```bash
go run ./cmd/bot
```

### 5) Build the binary
```bash
go build -o ip_bot ./cmd/bot
./ip_bot
```

## Repository Structure

```
ip_accounting_bot/
├── cmd/
│   └── bot/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go               # Main application logic and runner management
│   │   ├── runner.go            # Runner interface and concurrent execution
│   │   └── telegram_runner.go   # Telegram bot runner implementation
│   ├── config/
│   │   └── config.go            # Configuration loading from environment
│   └── telegram/
│       ├── client.go            # Telegram Bot API HTTP client
│       ├── types.go             # Telegram API data structures
│       └── updates.go           # Telegram API methods (getUpdates, sendMessage)
├── go.mod                       # Go module definition and dependencies
├── go.sum                       # Go module checksums
├── LICENSE                      # Business Source License 1.1
└── README.md                    # Project documentation
```

### File Descriptions

- **`cmd/bot/main.go`** - Application entry point, initializes configuration, creates application and starts Telegram bot
- **`internal/app/app.go`** - Main application logic, manages registration and execution of various components (runners)
- **`internal/app/runner.go`** - Runner interface and function for concurrent execution of all registered components
- **`internal/app/telegram_runner.go`** - Telegram bot implementation, processes incoming messages and sends echo responses
- **`internal/config/config.go`** - Configuration loading from environment variables with .env file support
- **`internal/telegram/client.go`** - HTTP client for Telegram Bot API, includes error handling and JSON parsing
- **`internal/telegram/types.go`** - Data structures for working with Telegram API (User, Chat, Message, Update)
- **`internal/telegram/updates.go`** - Methods for getting updates and sending messages via Telegram Bot API
- **`go.mod`** - Go module definition and dependencies (godotenv for .env file loading)
- **`go.sum`** - Go module dependency checksums

## License
This project is distributed under the [Business Source License 1.1](LICENSE) — free for non-commercial use; commercial use requires a separate agreement.
