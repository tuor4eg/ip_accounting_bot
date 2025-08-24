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
- Go 1.24+
- PostgreSQL
- Redis (cache and rate limiting)
- Telegram Bot API

## Installation & Run

### 1) Clone the repository
```bash
git clone https://github.com/tuor4eg/ip_accounting_bot.git
cd ip_accounting_bot
```

### 2) Install dependencies
Make sure you have **Go 1.24+** installed.
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

### 4) Run database migrations
```bash
go run ./cmd/migrate
```

### 5) Run the bot
```bash
go run ./cmd/bot
```

### 6) Build the binary
```bash
go build -o ip_bot ./cmd/bot
./ip_bot
```

## Repository Structure

```
ip_accounting_bot/
├── cmd/
│   ├── bot/
│   │   └── main.go                           # Bot application entry point
│   └── migrate/
│       └── main.go                           # Database migration entry point
├── internal/
│   ├── app/
│   │   ├── app.go                           # Main application logic and runner management
│   │   ├── handle_telegram_update.go        # Telegram update processing logic
│   │   ├── run_telegram_polling.go          # Telegram polling implementation
│   │   ├── runner.go                        # Runner interface and concurrent execution
│   │   ├── service.go                       # Service layer interface and implementation
│   │   ├── store.go                         # Storage layer interface and implementation
│   │   └── telegram_runner.go               # Telegram bot runner implementation
│   ├── bot/
│   │   ├── handlers_add.go                  # Add income command handler
│   │   ├── handlers_help.go                 # Help command handler
│   │   ├── handlers_start.go                # Start command handler
│   │   ├── handlers_total.go                # Total income command handler
│   │   ├── parse.go                         # Message parsing utilities
│   │   ├── router_dispatch.go               # Message routing and dispatch logic
│   │   └── text.go                          # Bot text messages and templates
│   ├── config/
│   │   └── config.go                        # Configuration loading from environment
│   ├── logging/
│   │   ├── logging.go                       # Logging configuration and setup
│   │   └── pkglogging.go                    # Package-level logging utilities
│   ├── migrations/
│   │   ├── embed.go                         # SQL migrations embedding
│   │   ├── lock.go                          # Database migration locking mechanism
│   │   ├── runner.go                        # Migration execution logic
│   │   └── sql/
│   │       └── 0001_init.up.sql             # Initial database schema
│   ├── money/
│   │   ├── format.go                        # Money formatting utilities
│   │   └── parse.go                         # Money parsing utilities
│   ├── period/
│   │   └── quarter.go                       # Quarter period calculations
│   ├── service/
│   │   └── income.go                        # Income business logic service
│   ├── storage/
│   │   └── postgres/
│   │       ├── base.go                      # Base database connection and operations
│   │       ├── identities.go                # User identity storage operations
│   │       └── incomes.go                   # Income data storage operations
│   └── telegram/
│       ├── client.go                        # Telegram Bot API HTTP client
│       ├── types.go                         # Telegram API data structures
│       └── updates.go                       # Telegram API methods (getUpdates, sendMessage)
├── go.mod                                   # Go module definition and dependencies
├── go.sum                                   # Go module checksums
├── .gitignore                               # Git ignore rules
├── CHANGELOG.md                             # Project changelog
├── LICENSE                                  # Business Source License 1.1
└── README.md                                # Project documentation
```

### File Descriptions

#### Project Files
- **`.gitignore`** - Git ignore rules for Go projects, excludes binaries, test files, coverage reports, and environment files

#### Command Line Applications
- **`cmd/bot/main.go`** - Bot application entry point, initializes configuration, creates application and starts Telegram bot
- **`cmd/migrate/main.go`** - Database migration application entry point

#### Application Core
- **`internal/app/app.go`** - Main application logic, manages registration and execution of various components (runners)
- **`internal/app/handle_telegram_update.go`** - Telegram update processing logic and message handling
- **`internal/app/run_telegram_polling.go`** - Telegram polling implementation for receiving updates
- **`internal/app/runner.go`** - Runner interface and function for concurrent execution of all registered components
- **`internal/app/service.go`** - Service layer interface and implementation for business logic
- **`internal/app/store.go`** - Storage layer interface and implementation for data persistence
- **`internal/app/telegram_runner.go`** - Telegram bot implementation, processes incoming messages and sends responses

#### Bot Handlers
- **`internal/bot/handlers_add.go`** - Add income command handler implementation
- **`internal/bot/handlers_help.go`** - Help command handler implementation
- **`internal/bot/handlers_start.go`** - Start command handler implementation
- **`internal/bot/handlers_total.go`** - Total income command handler implementation
- **`internal/bot/parse.go`** - Message parsing utilities for extracting commands and parameters
- **`internal/bot/router_dispatch.go`** - Message routing and dispatch logic to appropriate handlers
- **`internal/bot/text.go`** - Bot text messages, templates and localization

#### Configuration & Logging
- **`internal/config/config.go`** - Configuration loading from environment variables with .env file support
- **`internal/logging/logging.go`** - Logging configuration and setup
- **`internal/logging/pkglogging.go`** - Package-level logging utilities

#### Database & Migrations
- **`internal/migrations/embed.go`** - SQL migrations embedding for Go binary
- **`internal/migrations/lock.go`** - Database migration locking mechanism to prevent concurrent migrations
- **`internal/migrations/runner.go`** - Migration execution logic and version management
- **`internal/migrations/sql/0001_init.up.sql`** - Initial database schema creation

#### Business Logic
- **`internal/money/format.go`** - Money formatting utilities for displaying currency amounts
- **`internal/money/parse.go`** - Money parsing utilities for handling currency amounts
- **`internal/period/quarter.go`** - Quarter period calculations and date utilities
- **`internal/service/income.go`** - Income business logic service layer

#### Data Storage
- **`internal/storage/postgres/base.go`** - Base database connection and common operations
- **`internal/storage/postgres/identities.go`** - User identity storage operations
- **`internal/storage/postgres/incomes.go`** - Income data storage operations

#### Telegram Integration
- **`internal/telegram/client.go`** - HTTP client for Telegram Bot API, includes error handling and JSON parsing
- **`internal/telegram/types.go`** - Data structures for working with Telegram API (User, Chat, Message, Update)
- **`internal/telegram/updates.go`** - Methods for getting updates and sending messages via Telegram Bot API

## License
This project is distributed under the [Business Source License 1.1](LICENSE) — free for non-commercial use; commercial use requires a separate agreement.
