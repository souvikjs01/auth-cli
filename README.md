# CLI Authentication System

A secure command-line login system built with Go that supports user registration, authentication, optional TOTP-based 2FA, and session management. Runs in Docker containers with PostgreSQL for persistence.

## Features

- **User registration** with username and password
- **Login** with username and password (+ TOTP if enabled)
- **TOTP-based 2FA** (Google Authenticator compatible)
- **Secure password storage** using bcrypt hashing
- **Account lockout** after multiple failed login attempts
- **Session management** with configurable timeout
- **Interactive CLI** with command history

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go 1.25+](https://go.dev/dl/) (only for local development)

## Quick Start (Docker)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/souvikjs01/auth-cli.git
   cd auth-cli
   ```

2. **Create environment file:**
   ```bash
   cp .env.example .env
   ```

3. **Run the application:**
   ```bash
   docker compose run --rm --build app
   ```
   This starts PostgreSQL automatically and opens the interactive CLI.

4. **Use the CLI:**
   ```
   auth> register
   auth> login
   auth> help
   ```

## Local Development

1. **Start the database:**
   ```bash
   docker compose up postgres -d
   ```

2. **Create environment file:**
   ```bash
   cp .env.example .env
   ```

3. **Run the application:**
   ```bash
   go run .
   ```

## Configuration

All configuration is managed via the `.env` file:

| Variable             | Description                          | Default      |
|----------------------|--------------------------------------|--------------|
| `POSTGRES_USER`      | PostgreSQL username                  | `postgres`   |
| `POSTGRES_PASSWORD`  | PostgreSQL password                  | `postgres`   |
| `POSTGRES_HOST`      | PostgreSQL host                      | `localhost`   |
| `POSTGRES_DB`        | PostgreSQL database name             | `auth_cli`   |
| `POSTGRES_PORT`      | PostgreSQL port                      | `5432`       |
| `SESSION_TIMEOUT`    | Session expiration duration          | `30m`        |
| `MAX_LOGIN_ATTEMPTS` | Failed attempts before account lock  | `5`          |
| `LOCKOUT_DURATION`   | Account lockout duration             | `15m`        |
| `TOTP_ISSUER`        | Issuer name for TOTP codes           | `CLI_Auth`   |

## Available Commands

### Before Login
| Command    | Description              |
|------------|--------------------------|
| `register` | Create a new account     |
| `login`    | Login with credentials   |
| `help`     | Show available commands  |
| `exit`     | Quit the application     |

### After Login
| Command       | Description                |
|---------------|----------------------------|
| `whoami`      | Show current user details  |
| `enable-2fa`  | Enable TOTP-based MFA      |
| `disable-2fa` | Disable MFA                |
| `logout`      | End current session        |
| `help`        | Show available commands    |
| `exit`        | Quit the application       |

## Database Schema

The application uses GORM's AutoMigrate to manage the database schema automatically. Two tables are created:

### `users`
| Column           | Type         | Description                       |
|------------------|--------------|-----------------------------------|
| `id`             | `uint`       | Primary key (auto-increment)      |
| `username`       | `varchar(100)` | Unique username                 |
| `password_hash`  | `text`       | bcrypt hashed password            |
| `created_at`     | `timestamp`  | Registration date                 |
| `last_login_at`  | `timestamp`  | Last successful login (nullable)  |
| `failed_attempts`| `int`        | Consecutive failed login count    |
| `locked_until`   | `timestamp`  | Account lockout expiry (nullable) |
| `mfa_enabled`    | `boolean`    | Whether 2FA is active             |
| `mfa_secret`     | `text`       | TOTP secret key (nullable)        |

### `sessions`
| Column       | Type         | Description                        |
|--------------|--------------|------------------------------------|
| `id`         | `varchar(64)` | Primary key (random hex token)   |
| `user_id`    | `uint`       | Foreign key → `users.id`          |
| `created_at` | `timestamp`  | Session creation time              |
| `expires_at` | `timestamp`  | Session expiry time                |

## Project Structure

```
cli-auth/
├── cmd/                        # CLI commands and shell
│   ├── app.go                  # Application initialization
│   ├── auth.go                 # Auth commands (register, login, etc.)
│   ├── input.go                # Input helpers (prompts, passwords)
│   ├── root.go                 # Root cobra command
│   └── shell.go                # Interactive readline shell
├── internals/
│   ├── app/                    # Application orchestration layer
│   ├── config/                 # Configuration loading (env/viper)
│   ├── db/                     # Database connection and migrations
│   ├── models/                 # GORM data models
│   ├── repositories/           # Database access layer
│   ├── security/               # Password hashing (bcrypt)
│   └── service/                # Business logic (auth, sessions, TOTP)
├── Dockerfile                  # Multi-stage Go build
├── docker-compose.yml          # PostgreSQL + app services
├── .env.example                # Example environment configuration
└── go.mod                      # Go module definition
```

## Security

- Passwords are hashed with **bcrypt** (default cost factor)
- Sessions use **cryptographically random** 256-bit tokens
- Accounts are **locked** after 5 failed login attempts (configurable)
- Session tokens **expire** after 30 minutes (configurable)
- TOTP secrets are stored only after successful code verification
