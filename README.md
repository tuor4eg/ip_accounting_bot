# IP Accounting Bot

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://telegram.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)

**IP Accounting Bot** is a Telegram bot for sole proprietors that helps track income and reminds about tax payments.

## Features (MVP)
- Slash commands:
  - `/start`, `/help`
  - `/add <amount> [note]` — manual income input (kopecks, no floats)
  - `/total` — current quarter totals (sum + 6% tax)
  - `/undo` — void (soft-delete) the newest income in the current quarter
- Deterministic money math: `int64` in kopecks, no floats
- UTC dates (stored as `DATE`), quarter bounds are **inclusive**
- Soft-delete via `voided_at`, aggregates use only active rows

## Tech Stack
- Go (version per `go.mod`)
- PostgreSQL
- In-memory store for local tests
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
make deps
```

### 3) Set environment variables
```bash
make deps
```
Fill the file with your valid data

### 4) Run database migrations
```bash
make migrate
```

### 5) Build the binary
```bash
make build-bot
```

### 6) Run the bot
```bash
make run-bot
```
(Also build the binary for running)

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
│   │   ├── errors.go                        # Bot error handling and custom errors
│   │   ├── handlers_add.go                  # Add income command handler
│   │   ├── handlers_add_test.go             # Add income command handler tests
│   │   ├── handlers_help.go                 # Help command handler
│   │   ├── handlers_start.go                # Start command handler
│   │   ├── handlers_total.go                # Total income command handler
│   │   ├── handlers_undo.go                 # Undo last action command handler
│   │   ├── parse.go                         # Message parsing utilities
│   │   ├── router_dispatch.go               # Message routing and dispatch logic
│   │   └── text.go                          # Bot text messages and templates
│   ├── config/
│   │   └── config.go                        # Configuration loading from environment
│   ├── domain/
│   │   └── const.go                         # Domain constants and definitions
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
│   │   ├── format_test.go                   # Money formatting tests
│   │   ├── parse.go                         # Money parsing utilities
│   │   └── parse_test.go                    # Money parsing tests
│   ├── period/
│   │   ├── quarter.go                       # Quarter period calculations
│   │   └── quarter_test.go                  # Quarter period tests
│   ├── service/
│   │   └── income.go                        # Income business logic service
│   ├── storage/
│   │   ├── memstore/
│   │   │   ├── base.go                      # In-memory storage base implementation
│   │   │   ├── identities.go                # In-memory user identity storage
│   │   │   ├── incomes.go                   # In-memory income data storage
│   │   │   └── payments.go                  # In-memory payments data storage
│   │   └── postgres/
│   │       ├── base.go                      # Base database connection and operations
│   │       ├── identities.go                # User identity storage operations
│   │       ├── incomes.go                   # Income data storage operations
│   │       └── payments.go                  # PostgreSQL payments data storage
│   ├── telegram/
│   │   ├── client.go                        # Telegram Bot API HTTP client
│   │   ├── types.go                         # Telegram API data structures
│   │   └── updates.go                       # Telegram API methods (getUpdates, sendMessage)
│   └── validate/
│       ├── validate.go                      # Data validation utilities
│       └── wrap.go                          # Validation wrapper functions

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
- **`internal/bot/errors.go`** - Bot error handling, custom error types and error management
- **`internal/bot/handlers_add.go`** - Add income command handler implementation
- **`internal/bot/handlers_add_test.go`** - Tests for add income command handler
- **`internal/bot/handlers_help.go`** - Help command handler implementation
- **`internal/bot/handlers_start.go`** - Start command handler implementation
- **`internal/bot/handlers_total.go`** - Total income command handler implementation
- **`internal/bot/handlers_undo.go`** - Undo last action command handler implementation
- **`internal/bot/parse.go`** - Message parsing utilities for extracting commands and parameters
- **`internal/bot/router_dispatch.go`** - Message routing and dispatch logic to appropriate handlers
- **`internal/bot/text.go`** - Bot text messages, templates and localization

#### Configuration & Domain
- **`internal/config/config.go`** - Configuration loading from environment variables with .env file support
- **`internal/domain/const.go`** - Domain constants and business logic definitions

#### Logging
- **`internal/logging/logging.go`** - Logging configuration and setup
- **`internal/logging/pkglogging.go`** - Package-level logging utilities

#### Database & Migrations
- **`internal/migrations/embed.go`** - SQL migrations embedding for Go binary
- **`internal/migrations/lock.go`** - Database migration locking mechanism to prevent concurrent migrations
- **`internal/migrations/runner.go`** - Migration execution logic and version management
- **`internal/migrations/sql/0001_init.up.sql`** - Initial database schema creation

#### Business Logic
- **`internal/money/format.go`** - Money formatting utilities for displaying currency amounts
- **`internal/money/format_test.go`** - Tests for money formatting utilities
- **`internal/money/parse.go`** - Money parsing utilities for handling currency amounts
- **`internal/money/parse_test.go`** - Tests for money parsing utilities
- **`internal/period/quarter.go`** - Quarter period calculations and date utilities
- **`internal/period/quarter_test.go`** - Tests for quarter period calculations
- **`internal/service/income.go`** - Income business logic service layer

#### Data Storage
- **`internal/storage/memstore/base.go`** - In-memory storage base implementation for development/testing
- **`internal/storage/memstore/identities.go`** - In-memory user identity storage operations
- **`internal/storage/memstore/incomes.go`** - In-memory income data storage operations
- **`internal/storage/memstore/payments.go`** - In-memory payments data storage operations
- **`internal/storage/postgres/base.go`** - Base database connection and common operations
- **`internal/storage/postgres/identities.go`** - User identity storage operations
- **`internal/storage/postgres/incomes.go`** - Income data storage operations
- **`internal/storage/postgres/payments.go`** - PostgreSQL payments data storage operations

#### Telegram Integration
- **`internal/telegram/client.go`** - HTTP client for Telegram Bot API, includes error handling and JSON parsing
- **`internal/telegram/types.go`** - Data structures for working with Telegram API (User, Chat, Message, Update)
- **`internal/telegram/updates.go`** - Methods for getting updates and sending messages via Telegram Bot API

#### Validation
- **`internal/validate/validate.go`** - Data validation utilities and validation functions
- **`internal/validate/wrap.go`** - Validation wrapper functions for common validation patterns

## License
This project is distributed under the [Business Source License 1.1](LICENSE) — free for non-commercial use; commercial use requires a separate agreement.
