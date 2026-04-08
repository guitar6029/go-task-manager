🔧 Task Manager API (Go + Gin)


📌 Overview

A production-style backend API built with Go, focused on authentication, task management, and scalable system design. This project demonstrates clean architecture, RESTful API design, and containerized deployment.

🚀 Features
User registration & authentication (JWT-based)
Protected routes with middleware
CRUD operations for tasks
Pagination (limit & offset)
PostgreSQL database integration
Dockerized setup for easy deployment
Input validation & error handling

🛠 Tech Stack
Language: Go
Framework: Gin
Database: PostgreSQL
Auth: JWT
Containerization: Docker
Docs: Swagger

📂 Project Structure

go-task-manager/
├── api/ # HTTP handlers (Gin routes)
├── db/ # Database queries & connection logic
├── docs/ # Swagger generated files
├── middleware/ # Middleware (JWT, rate limiting)
├── model/ # Data models (structs)
├── nginx/ # Nginx reverse proxy config
├── service/ # Business logic layer
├── .env # Environment variables
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go # Entry point (CLI + API bootstrap)
├── README.md
└── tasks.db # SQLite database (dev)


⚙️ Getting Started
1. Clone the repo
git clone https://github.com/yourusername/task-manager.git
cd task-manager
2. Set environment variables

Create a .env file:

DB_URL=postgres://user:password@localhost:5432/taskdb
JWT_SECRET=your_secret
3. Run with Docker
docker-compose up --build
4. Run locally (without Docker)
go mod tidy
go run cmd/main.go
🔐 API Endpoints
Auth
POST /register → Create user
POST /login → Get JWT
Tasks (Protected)
GET /tasks → List tasks
POST /tasks → Create task
PUT /tasks/:id → Update task
DELETE /tasks/:id → Delete task
🧪 Testing

You can test endpoints using:

Postman
curl
Swagger UI (/swagger/index.html)

🧠 What I Learned
Structuring Go projects for scalability
Implementing JWT auth & middleware
Handling DB connections and queries cleanly
Building containerized backend services
📈 Future Improvements
Add refresh tokens
Rate limiting
Unit & integration tests
Role-based access control (RBAC)

<img width="1672" height="867" alt="image" src="https://github.com/user-attachments/assets/8f8e33f6-b81a-45f0-a0f9-9923a16458c2" />

📘 API Documentation (Swagger)

Interactive API documentation is available at:

http://localhost/swagger/index.html

<img width="2242" height="1101" alt="image" src="https://github.com/user-attachments/assets/78a2c5b4-ca29-4555-b7b4-bad30d64c791" />

