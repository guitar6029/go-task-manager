# Task Manager API

Task Manager is a production-style backend system built in Go, demonstrating authenticated task management, asynchronous job processing, and horizontally scalable API design.

It combines a REST API, background worker, and reverse proxy layer to simulate a more realistic distributed backend architecture than a basic CRUD app.

## Overview

The system is composed of three main components:

- `cmd/api`: Gin-based HTTP server handling authentication and task APIs
- `cmd/worker`: Background worker that processes asynchronous jobs from Redis
- `cmd/cli`: Lightweight local development shell

Task reads are synchronous. Task writes are processed asynchronously through a Redis-backed queue, which separates request handling from mutation processing.

## Features

- User registration and login with JWT-based authentication
- Protected task endpoints
- Task listing with pagination and completion filtering
- Asynchronous task creation, deletion, and completion updates via Redis jobs
- Worker retry handling with dead-letter queue support (`jobs:failed`)
- PostgreSQL-backed persistence
- Redis-backed queueing, plus cache hooks for task reads
- Basic API rate limiting middleware
- NGINX reverse proxy with:
  - load balancing across multiple API instances
  - HTTPS support via SSL/TLS termination
- Docker Compose multi-service setup
- Swagger UI for API exploration
- Health check endpoint for container orchestration

## System Design Highlights

- Asynchronous job processing using Redis queues and a dedicated worker
- Separation of read and write paths: synchronous reads, queued mutations
- Horizontal scaling through multiple API instances behind NGINX
- TLS termination at the proxy layer
- Retry and dead-letter queue handling for failed jobs

## Architecture

```text
Client (HTTPS)
  |
  v
NGINX (TLS termination + load balancing)
  |
  +--> API instance 1
  |
  +--> API instance 2
          |
          +--> PostgreSQL
          |
          +--> Redis
                 |
                 +--> jobs
                 +--> jobs:failed
                        |
                        v
                     Worker
```

## Request Flow

- `GET /tasks` uses the synchronous read path and is backed by PostgreSQL, with Redis cache lookup hooks in the service layer
- `POST /tasks`, `DELETE /tasks/:id`, and `PATCH /tasks/:id` enqueue jobs in Redis
- The worker processes queued jobs, retries failures up to 3 times, and moves exhausted jobs to `jobs:failed`
- `GET /health` checks database connectivity

## Tech Stack

- Language: Go
- Framework: Gin
- Database: PostgreSQL
- Queue/Cache: Redis
- Proxy: NGINX
- Containerization: Docker Compose
- Docs: Swagger

## Project Structure

```text
go-task-manager/
├── cmd/
│   ├── api/
│   ├── cli/
│   └── worker/
├── docs/
├── internal/
│   ├── api/
│   ├── cache/
│   ├── config/
│   ├── db/
│   ├── middleware/
│   ├── model/
│   ├── queue/
│   ├── redis/
│   └── service/
├── nginx/
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## Environment Variables

Create `.env` or `.env.local`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskdb
REDIS_ADDR=localhost:6379
APP_ENV=local
```

Notes:

- `.env.local` is loaded before `.env`
- `APP_ENV` is optional, but useful for local development logging
- Database tables are created automatically on startup if they do not already exist

## Running With Docker

```bash
docker compose up --build
```

This starts:

- 2 API instances
- Worker
- PostgreSQL
- Redis
- NGINX on ports `80` and `443`, with HTTP to HTTPS redirect

## Running Locally

Run the API:

```bash
go run ./cmd/api
```

Run the worker in a separate terminal:

```bash
go run ./cmd/worker
```

Run the CLI:

```bash
go run ./cmd/cli
```

## API Usage

### Public Endpoints

- `POST /register`
- `POST /login`
- `GET /health`
- `GET /swagger/index.html`

### Protected Endpoints

Require:

```text
Authorization: Bearer <token>
```

Endpoints:

- `GET /tasks`
- `POST /tasks`
- `DELETE /tasks/:id`
- `PATCH /tasks/:id`

### Important Behavior

Task mutation endpoints are asynchronous. They enqueue a job and currently return `202 Accepted` with a queued response instead of returning the final mutated task immediately.

### Example

```bash
curl https://localhost/tasks \
  -H "Authorization: Bearer <token>" \
  -k
```

`-k` is required for self-signed TLS in local development.

## Swagger

- `http://localhost:8080/swagger/index.html`
- `https://localhost/swagger/index.html`

## Windows Notes

If using Git Bash with OpenSSL:

```bash
MSYS_NO_PATHCONV=1 openssl ...
```

See [DEV_NOTES.MD](DEV_NOTES.MD) for details.

## Current Limitations

- JWT signing currently uses a hardcoded secret and should move to secure configuration
- CLI functionality is intentionally minimal
- Observability is still limited to basic logging
- Rate limiting is present, but still simple
- Cache helpers exist, but the read path is not yet a fully realized cache-write-through flow

## Production Considerations

In a real deployment, this system would likely evolve to:

- use managed services such as RDS and ElastiCache
- replace local NGINX with a cloud load balancer
- store secrets in a secure manager
- add structured logging and monitoring
- introduce CI/CD pipelines

## Future Improvements

- Move JWT secret to environment or secret management
- Add unit and integration tests
- Implement per-user task ownership
- Introduce DB migrations
- Improve observability with metrics and tracing
- Enhance retry strategy with backoff
