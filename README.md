# IP Accounting Bot

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

## License
This project is distributed under the [Business Source License 1.1](LICENSE) — free for non-commercial use; commercial use requires a separate agreement.
