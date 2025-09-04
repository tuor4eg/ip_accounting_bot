# IP Accounting Bot

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://telegram.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)

**IP Accounting Bot** is a Telegram bot for sole proprietors that helps track income and reminds about tax payments.

## Features (MVP)
- **Slash commands:**
  - `/start` — brief usage guide for the bot
  - `/help` — detailed help for all commands
  - `/add <amount> [note]` — add income (in kopecks, no floats)
  - `/add_contrib <amount> [note]` — add contribution
  - `/add_advance <amount> [note]` — add advance payment
  - `/total` — current quarter totals (income sum and 6% tax)
  - `/undo` — undo last income for the quarter
  - `/undo_contrib` — undo last contribution
  - `/undo_advance` — undo last advance payment
- **Amount format:** supports spaces/dots/commas as thousand separators, also "10р 50к" format
- **Deterministic math:** `int64` in kopecks, no floats
- **UTC dates** (stored as `DATE`), quarter bounds are **inclusive**
- **Soft delete** via `voided_at`, aggregates use only active rows

### Command usage examples:
```
/add 1000                    # Add income of 1000 rubles
/add 1 234,56 order #42      # Add income with note
/add 10р 50к advance         # Add income in "rubles kopecks" format
/add_contrib 5000            # Add contribution of 5000 rubles
/add_advance 3000            # Add advance payment of 3000 rubles
/total                       # Show current quarter totals
/undo                        # Undo last income
/undo_contrib                # Undo last contribution
/undo_advance                # Undo last advance payment
```

## Tech Stack
- Go (version per `go.mod`)
- PostgreSQL
- In-memory store for local tests
- Telegram Bot API

## Code Style

### Package Organization
Each package should follow a consistent structure with separate files for different concerns:

- **`types.go`** - Contains all type definitions and structs for the package
- **`interfaces.go`** - Contains all interface definitions for the package
- **`errors.go`** - Contains all error definitions and custom error types
- **Main files** - Contain only business logic, functions, and methods

### Code Separation Rules
- **Never** define types or errors in executable files (files with business logic)
- **Always** move types to `types.go` and errors to `errors.go`
- Keep main files focused on functionality, not data structure definitions
- Use clear, descriptive names for types and error variables
- Add comprehensive comments for complex types and error scenarios

### Benefits
- **Maintainability**: Easy to locate and modify types/errors
- **Readability**: Clear separation of concerns
- **Consistency**: Uniform structure across all packages
- **Documentation**: Self-documenting code organization

### Example Package Structure
```
internal/service/
├── types.go          # All type definitions
├── interfaces.go     # All interface definitions
├── errors.go         # All error definitions  
├── income.go         # Business logic only
├── payment.go        # Business logic only
└── total.go          # Business logic only
```

**Good** - types.go:
```go
// IncomeService handles income-related business logic
type IncomeService struct {
    store IncomeStore
}
```

**Bad** - income.go:
```go
// ❌ Don't define types in business logic files
type IncomeService struct {
    store IncomeStore
}
```

### When to Create These Files
- **`types.go`** - Create when package has structs or custom types
- **`interfaces.go`** - Create when package has interface definitions
- **`errors.go`** - Create when package defines custom errors or error constants
- **Don't create empty files** - Only create if there are actual types/interfaces/errors to define
- **Follow Go conventions** - Use descriptive names and proper documentation

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
make env
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
├── bin/                                     # Compiled binaries
│   ├── ip_bot                               # Bot binary
│   └── migrate                              # Migration binary
├── cmd/
│   ├── bot/
│   │   └── main.go                           # Bot application entry point
│   └── migrate/
│       └── main.go                           # Database migration entry point
├── config/                                  # Application configuration
│   ├── config.go                            # Configuration loading from environment
│   ├── errors.go                            # Configuration error definitions
│   └── types.go                             # Configuration type definitions
├── migrations/                              # Database migrations
│   ├── embed.go                             # SQL migrations embedding
│   ├── errors.go                            # Migration error definitions
│   ├── lock.go                              # Database migration locking mechanism
│   ├── lock_errors.go                       # Lock mechanism error definitions
│   ├── runner.go                            # Migration execution logic
│   ├── types.go                             # Migration type definitions
│   └── sql/
│       └── 0001_init.up.sql                 # Initial database schema
├── internal/
│   ├── app/
│   │   ├── app.go                           # Main application logic and runner management
│   │   ├── errors.go                        # Application error definitions
│   │   ├── handle_telegram_update.go        # Telegram update processing logic
│   │   ├── interfaces.go                    # Application interface definitions
│   │   ├── run_telegram_polling.go          # Telegram polling implementation
│   │   ├── runner.go                        # Runner interface and concurrent execution
│   │   ├── service.go                       # Service layer interface and implementation
│   │   ├── store.go                         # Storage layer interface and implementation
│   │   ├── telegram_runner.go               # Telegram bot runner implementation
│   │   └── types.go                         # Application type definitions
│   ├── bot/
│   │   ├── deps.go                          # Bot dependencies and initialization
│   │   ├── errors.go                        # Bot error handling and custom errors
│   │   ├── handlers_add.go                  # Add income command handler
│   │   ├── handlers_add_advance.go          # Advanced add income handler
│   │   ├── handlers_add_contrib.go          # Contributory add income handler
│   │   ├── handlers_add_test.go             # Add income command handler tests
│   │   ├── handlers_help.go                 # Help command handler
│   │   ├── handlers_start.go                # Start command handler
│   │   ├── handlers_total.go                # Total income command handler
│   │   ├── handlers_undo.go                 # Undo last action command handler
│   │   ├── handlers_undo_advance.go         # Advanced undo handler
│   │   ├── handlers_undo_contrib.go         # Contributory undo handler
│   │   ├── parse.go                         # Message parsing utilities
│   │   ├── router_dispatch.go               # Message routing and dispatch logic
│   │   ├── text.go                          # Bot text messages and templates
│   │   ├── types.go                         # Bot type definitions and interfaces
│   │   ├── validate.go                      # Bot-specific validation
│   │   └── router_dispatch.go               # Message routing and dispatch logic
│   ├── crypto/
│   │   ├── crypto.go                        # Cryptographic utilities and functions
│   │   ├── errors.go                        # Cryptographic error definitions
│   │   └── types.go                         # Cryptographic type definitions
│   ├── cryptostore/
│   │   ├── base.go                          # Base cryptographic storage implementation
│   │   ├── errors.go                        # Cryptographic storage error definitions
│   │   ├── interface.go                     # Cryptographic storage interface
│   │   ├── README.md                        # Cryptographic storage documentation
│   │   └── types.go                         # Cryptographic storage type definitions
│   ├── domain/
│   │   ├── const.go                         # Domain constants and definitions
│   │   ├── interfaces.go                    # Domain interface definitions
│   │   ├── totals.go                        # Domain totals and aggregates logic
│   │   └── types.go                         # Domain type definitions
│   ├── logging/
│   │   ├── logging.go                       # Logging configuration and setup
│   │   └── pkglogging.go                    # Package-level logging utilities
│   ├── money/
│   │   ├── errors.go                        # Money error definitions
│   │   ├── format.go                        # Money formatting utilities
│   │   ├── format_test.go                   # Money formatting tests
│   │   ├── parse.go                         # Money parsing utilities
│   │   └── parse_test.go                    # Money parsing tests
│   ├── period/
│   │   ├── quarter.go                       # Quarter period calculations
│   │   ├── quarter_test.go                  # Quarter period tests
│   │   ├── year.go                          # Year period calculations
│   │   └── year_test.go                     # Year period tests
│   ├── service/
│   │   ├── income.go                        # Income business logic service
│   │   ├── interfaces.go                    # Service interface definitions
│   │   ├── payment.go                       # Payment business logic service
│   │   ├── total.go                         # Total calculation service
│   │   └── types.go                         # Service type definitions
│   ├── storage/
│   │   ├── memstore/
│   │   │   ├── base.go                      # In-memory storage base implementation
│   │   │   ├── identities.go                # In-memory user identity storage
│   │   │   ├── incomes.go                   # In-memory income data storage
│   │   │   ├── payments.go                  # In-memory payments data storage
│   │   │   └── types.go                     # In-memory storage type definitions
│   │   └── postgres/
│   │       ├── base.go                      # Base database connection and operations
│   │       ├── errors.go                    # PostgreSQL error definitions
│   │       ├── identities.go                # User identity storage operations
│   │       ├── incomes.go                   # Income data storage operations
│   │       ├── payments.go                  # PostgreSQL payments data storage
│   │       └── types.go                     # PostgreSQL storage type definitions
│   ├── tax/
│   │   ├── policy.go                        # Tax policy interface and implementation
│   │   ├── policy_test.go                   # Tax policy tests
│   │   ├── static_default.go                # Default static tax policy
│   │   ├── tax.go                           # Tax calculation logic
│   │   ├── tax_test.go                      # Tax calculation tests
│   │   └── types.go                         # Tax type definitions
│   ├── telegram/
│   │   ├── client.go                        # Telegram Bot API HTTP client
│   │   ├── errors.go                        # Telegram error definitions
│   │   ├── types.go                         # Telegram API data structures
│   │   └── updates.go                       # Telegram API methods (getUpdates, sendMessage)
│   └── validate/
│       ├── errors.go                        # Validation error definitions
│       ├── validate.go                      # Data validation utilities
│       └── wrap.go                          # Validation wrapper functions
├── go.mod                                   # Go module definition and dependencies
├── go.sum                                   # Go module checksums
├── .gitignore                               # Git ignore rules
├── CHANGELOG.md                             # Project changelog
├── LICENSE                                  # Business Source License 1.1
├── Makefile                                 # Build and development commands
└── README.md                                # Project documentation
```

### File Descriptions

#### Project Files
- **`bin/`** - Directory containing compiled binaries
  - **`ip_bot`** - Compiled bot application binary
  - **`migrate`** - Compiled database migration binary
- **`.gitignore`** - Git ignore rules for Go projects, excludes binaries, test files, coverage reports, and environment files
- **`Makefile`** - Build automation and development commands

#### Command Line Applications
- **`cmd/bot/main.go`** - Bot application entry point, initializes configuration, creates application and starts Telegram bot
- **`cmd/migrate/main.go`** - Database migration application entry point

#### Application Core
- **`internal/app/app.go`** - Main application logic, manages registration and execution of various components (runners)
- **`internal/app/errors.go`** - Application error definitions and error handling
- **`internal/app/handle_telegram_update.go`** - Telegram update processing logic and message handling
- **`internal/app/run_telegram_polling.go`** - Telegram polling implementation for receiving updates
- **`internal/app/runner.go`** - Runner interface and function for concurrent execution of all registered components
- **`internal/app/service.go`** - Service layer interface and implementation for business logic
- **`internal/app/store.go`** - Storage layer interface and implementation for data persistence
- **`internal/app/telegram_runner.go`** - Telegram bot implementation, processes incoming messages and sends responses
- **`internal/app/types.go`** - Application type definitions and structures

#### Bot Handlers
- **`internal/bot/deps.go`** - Bot dependencies and initialization logic
- **`internal/bot/errors.go`** - Bot error handling, custom error types and error management
- **`internal/bot/handlers_add.go`** - Add income command handler implementation
- **`internal/bot/handlers_add_advance.go`** - Advanced add income handler implementation
- **`internal/bot/handlers_add_contrib.go`** - Contributory add income handler implementation
- **`internal/bot/handlers_add_test.go`** - Tests for add income command handler
- **`internal/bot/handlers_help.go`** - Help command handler implementation
- **`internal/bot/handlers_start.go`** - Start command handler implementation
- **`internal/bot/handlers_total.go`** - Total income command handler implementation
- **`internal/bot/handlers_undo.go`** - Undo last action command handler implementation
- **`internal/bot/handlers_undo_advance.go`** - Advanced undo handler implementation
- **`internal/bot/handlers_undo_contrib.go`** - Contributory undo handler implementation
- **`internal/bot/parse.go`** - Message parsing utilities for extracting commands and parameters
- **`internal/bot/router_dispatch.go`** - Message routing and dispatch logic to appropriate handlers
- **`internal/bot/text.go`** - Bot text messages, templates and localization
- **`internal/bot/types.go`** - Bot type definitions, interfaces and dependency structures
- **`internal/bot/validate.go`** - Bot-specific validation functions

#### Configuration
- **`config/config.go`** - Configuration loading from environment variables with .env file support
- **`config/errors.go`** - Configuration error definitions and error handling
- **`config/types.go`** - Configuration type definitions and structures

#### Database Migrations
- **`migrations/embed.go`** - SQL migrations embedding for Go binary
- **`migrations/errors.go`** - Migration error definitions and error handling
- **`migrations/lock.go`** - Database migration locking mechanism to prevent concurrent migrations
- **`migrations/lock_errors.go`** - Lock mechanism error definitions and error handling
- **`migrations/runner.go`** - Migration execution logic and version management
- **`migrations/types.go`** - Migration type definitions and structures
- **`migrations/sql/0001_init.up.sql`** - Initial database schema creation

#### Domain
- **`internal/domain/const.go`** - Domain constants and business logic definitions
- **`internal/domain/interfaces.go`** - Domain interface definitions
- **`internal/domain/types.go`** - Domain type definitions and structures

#### Logging
- **`internal/logging/logging.go`** - Logging configuration and setup
- **`internal/logging/pkglogging.go`** - Package-level logging utilities

#### Business Logic
- **`internal/money/errors.go`** - Money error definitions and error handling
- **`internal/money/format.go`** - Money formatting utilities for displaying currency amounts
- **`internal/money/format_test.go`** - Tests for money formatting utilities
- **`internal/money/parse.go`** - Money parsing utilities for handling currency amounts
- **`internal/money/parse_test.go`** - Tests for money parsing utilities
- **`internal/period/quarter.go`** - Quarter period calculations and date utilities
- **`internal/period/quarter_test.go`** - Tests for quarter period calculations
- **`internal/period/year.go`** - Year period calculations and date utilities
- **`internal/period/year_test.go`** - Tests for year period calculations
- **`internal/service/income.go`** - Income business logic service layer
- **`internal/service/payment.go`** - Payment business logic service layer
- **`internal/service/total.go`** - Total calculation and aggregation service
- **`internal/service/types.go`** - Service type definitions and structures
- **`internal/tax/policy.go`** - Tax policy interface and implementation
- **`internal/tax/policy_test.go`** - Tests for tax policy implementation
- **`internal/tax/static_default.go`** - Default static tax policy implementation
- **`internal/tax/tax.go`** - Tax calculation logic and business rules
- **`internal/tax/tax_test.go`** - Tests for tax calculation logic
- **`internal/tax/types.go`** - Tax type definitions and structures

#### Cryptographic & Security
- **`internal/crypto/crypto.go`** - Cryptographic utilities and functions for security features
- **`internal/crypto/errors.go`** - Cryptographic error definitions and error handling
- **`internal/crypto/types.go`** - Cryptographic type definitions and structures
- **`internal/cryptostore/base.go`** - Base cryptographic storage implementation
- **`internal/cryptostore/errors.go`** - Cryptographic storage error definitions
- **`internal/cryptostore/interface.go`** - Cryptographic storage interface definition
- **`internal/cryptostore/README.md`** - Documentation for cryptographic storage system
- **`internal/cryptostore/types.go`** - Cryptographic storage type definitions

#### Data Storage
- **`internal/storage/memstore/base.go`** - In-memory storage base implementation for development/testing
- **`internal/storage/memstore/identities.go`** - In-memory user identity storage operations
- **`internal/storage/memstore/incomes.go`** - In-memory income data storage operations
- **`internal/storage/memstore/payments.go`** - In-memory payments data storage operations
- **`internal/storage/memstore/types.go`** - In-memory storage type definitions
- **`internal/storage/postgres/base.go`** - Base database connection and common operations
- **`internal/storage/postgres/errors.go`** - PostgreSQL error definitions and error handling
- **`internal/storage/postgres/identities.go`** - User identity storage operations
- **`internal/storage/postgres/incomes.go`** - Income data storage operations
- **`internal/storage/postgres/payments.go`** - PostgreSQL payments data storage operations
- **`internal/storage/postgres/types.go`** - PostgreSQL storage type definitions

#### Telegram Integration
- **`internal/telegram/client.go`** - HTTP client for Telegram Bot API, includes error handling and JSON parsing
- **`internal/telegram/errors.go`** - Telegram error definitions and error handling
- **`internal/telegram/types.go`** - Data structures for working with Telegram API (User, Chat, Message, Update)
- **`internal/telegram/updates.go`** - Methods for getting updates and sending messages via Telegram Bot API

#### Validation
- **`internal/validate/errors.go`** - Validation error definitions and error handling
- **`internal/validate/validate.go`** - Data validation utilities and validation functions
- **`internal/validate/wrap.go`** - Validation wrapper functions for common validation patterns

## License
This project is distributed under the [Business Source License 1.1](LICENSE) — free for non-commercial use; commercial use requires a separate agreement.
