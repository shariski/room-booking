# Room Booking API

A production-ready hotel room booking API built with Go, covering the full lifecycle from room discovery to reservation management. Built as a take-home assignment for the Senior Backend Engineer position at Bobobox.

## Tech Stack

- **Go 1.25** (stdlib `net/http`)
- **PostgreSQL 16** with exclusion constraints
- **Redis 7** for caching
- **golang-migrate** for database migrations
- **Docker Compose** for orchestration

## Quick Start

```bash
# Start all services (PostgreSQL, Redis, App) and run migrations
make start

# Run tests
make test

# Build Docker image
make build
```

The API will be available at `http://localhost:8080`.  
Swagger docs at `http://localhost:8080/swagger/index.html`.

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /users | No | Register a new user |
| POST | /login | No | Login, returns JWT token |
| GET | /rooms | Yes | List rooms, optional `?type=` filter |
| GET | /rooms/{id} | Yes | Get room detail (cached) |
| POST | /bookings | Yes | Create a booking |
| DELETE | /bookings/{id} | Yes | Cancel a booking (soft delete) |

## Architectural Decisions

### Why stdlib `net/http` over frameworks (Echo, Chi, Gin)

The API has 6 endpoints with no complex routing requirements. Since Go 1.22, the standard `ServeMux` supports method-based routing and path parameters natively. Adding a framework for this scope would introduce unnecessary dependencies without meaningful benefit. Middleware chaining is handled through simple function composition, which is idiomatic Go.

### Why raw SQL over GORM

This application is database-heavy — the core business rules (double-booking prevention, date validation, unique constraints) all live as PostgreSQL constraints. The application layer is essentially a thin interface on top of the database. Using an ORM would obscure the very thing the assignment asks to demonstrate: explicit control over database behavior.

With GORM, handling exclusion constraint violations, custom error codes (`23P01`, `23505`, `23514`), and precise query control would require frequent fallback to `db.Raw()` / `db.Exec()`, defeating the purpose of the ORM.

### Why Repository pattern over Hexagonal architecture

The assignment asks for clean, modular design patterns. Repository pattern was chosen over Hexagonal because:

1. The application has minimal business logic — almost every usecase method is a pass-through to the repository. Hexagonal's extra abstraction layers would add boilerplate without value.
2. The core rules live in the database (exclusion constraints, check constraints), not in domain logic. Dependency inversion toward a domain layer that barely exists is over-engineering.
3. The app could function entirely as raw SQL without a Go application layer. This makes a pragmatic, linear dependency flow (handler → usecase → repository → domain) more appropriate than ports-and-adapters.

### Why integration tests over unit tests with mocks

The assignment requests "unit tests for booking validation and concurrency handling." However, since all booking validation is enforced by PostgreSQL constraints (not application code), unit tests with mocked repositories would only prove the code can call a function — not that the constraints actually work.

Integration tests that hit a real database are the only meaningful way to verify:
- Exclusion constraints reject overlapping bookings
- Adjacent dates are correctly accepted (`[)` range semantics)
- Concurrent booking attempts result in exactly one success
- Soft-deleted bookings free up their date range for rebooking

This is a deliberate choice based on where the actual logic lives, not an omission.

## Race Condition Handling

### The Problem

Two users simultaneously attempt to book the same room for overlapping dates. Without protection, both `INSERT` statements could succeed, resulting in a double booking.

### The Solution: PostgreSQL Exclusion Constraint

```sql
CONSTRAINT no_overlapping_bookings
    EXCLUDE USING gist (
        room_id WITH =,
        daterange(start_date, end_date, '[)') WITH &&
    ) WHERE (deleted_at IS NULL)
```

This constraint operates at the database level, meaning it works correctly regardless of how many application instances are running. When two conflicting inserts arrive simultaneously, PostgreSQL serializes them internally — the first succeeds, the second receives error code `23P01` (exclusion violation), which the application maps to HTTP 409 Conflict.

The `[)` (start-inclusive, end-exclusive) range means a checkout on April 15 allows a new check-in on April 15, which follows standard hotel convention.

### Why not `SELECT FOR UPDATE`?

Pessimistic row-level locking was considered but rejected. It locks the room row during the transaction, which means even bookings for non-overlapping dates on the same room must wait for the lock to be released. This creates unnecessary contention and potential deadlocks without providing additional safety beyond what the exclusion constraint already guarantees.

### Why no explicit database transactions?

Every write operation in this application is a single SQL statement — booking creation is one `INSERT`, cancellation is one `UPDATE`. There are no multi-statement operations that require atomicity. Using `BEGIN/COMMIT` around a single statement adds overhead with no benefit.

## Caching Strategy

### Approach: Cache-Aside at the Usecase Level

`GET /rooms/{id}` implements a cache-aside pattern to minimize database load:

1. Check Redis for cached room data
2. On cache miss, query PostgreSQL
3. Store the result in Redis with a TTL (10 minutes)
4. On cache hit, return directly from Redis

### Why data-level caching, not HTTP response caching

The requirement says "minimize DB load," which points to reducing database queries, not reducing HTTP processing. Data-level caching in the usecase gives finer control over what gets cached and how it's invalidated.

### Redis failure handling

Redis is treated as optional infrastructure. If Redis is unavailable (connection error, timeout), the application falls back to querying PostgreSQL directly. Redis errors are logged but never propagated to the client. The database is always the source of truth.

### Cache invalidation

Room data is master data — it doesn't change when bookings are created or cancelled. Cache invalidation would only be relevant if room data itself were updated (e.g., name, type, description change), which is outside the scope of this API. TTL-based expiration is sufficient.

## Graceful Shutdown

The server listens for `SIGTERM` and `SIGINT` signals using `signal.NotifyContext`. When a signal is received:

1. The server stops accepting new connections
2. Active requests are given up to 10 seconds to complete
3. Database and Redis connections are closed cleanly
4. The process exits

This ensures no requests are dropped during deployment or restarts.

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
defer stop()

go func() {
    server.ListenAndServe()
}()

<-ctx.Done()
server.Shutdown(shutdownCtx)
db.Close()
rdb.Close()
```

## Soft Delete

Booking cancellation uses a soft-delete pattern — `DELETE /bookings/{id}` sets `deleted_at = NOW()` instead of removing the row. This preserves booking history for auditing purposes.

The exclusion constraint includes `WHERE (deleted_at IS NULL)`, so cancelled bookings don't block new bookings for the same dates. A `RETURNING` clause on the update returns the full booking data including the `deleted_at` timestamp.

## Database Schema

```sql
-- Users: minimal, supports JWT auth flow
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Rooms: master data, no soft delete (no delete use case)
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bookings: with exclusion constraint and soft delete
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID REFERENCES rooms(id),
    user_id UUID REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CHECK (end_date > start_date),
    CONSTRAINT no_overlapping_bookings
        EXCLUDE USING gist (
            room_id WITH =,
            daterange(start_date, end_date, '[)') WITH &&
        ) WHERE (deleted_at IS NULL)
);
```

## Project Structure

```
room-booking/
├── cmd/api/main.go           # Entrypoint, wiring, graceful shutdown
├── internal/
│   ├── domain/               # Entity structs (User, Room, Booking)
│   ├── model/                # Request/response structs, validation, errors
│   ├── handler/              # HTTP handlers, request parsing
│   ├── usecase/              # Business logic, caching, error mapping
│   ├── repository/           # PostgreSQL queries
│   ├── middleware/           # JWT auth, logging
│   └── config/               # Environment variable loading
├── migrations/               # SQL migration files (up/down)
├── docs/                     # Auto-generated Swagger docs
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Authentication

JWT-based authentication. Register via `POST /users`, login via `POST /login` to receive a token. Include the token in the `Authorization: Bearer <token>` header for protected endpoints.

The JWT middleware validates the token signature and expiry directly — no database lookup per request. The `user_id` from token claims is injected into the request context and used to associate bookings with the authenticated user.

## Running Tests

```bash
make test
```

Test scenarios cover:
- Successful booking creation
- Overlapping date rejection (exclusion constraint)
- Adjacent date acceptance (`[)` range)
- Start date after end date rejection (check constraint)
- Same dates on different rooms (constraint scoped per room)
- Soft delete and rebook on same dates
- Booking not found (non-existent ID)
- Concurrent booking attempts (10 goroutines, exactly 1 succeeds)
