What I Have Implemented 
🔹 1. RESTful User CRUD API

Create User

Get All Users (with pagination & filtering)

Get User by ID

Update User

Delete User (Soft Delete)

🔹 2. Clean Project Structure
user-crud-go/
│── config/        → Database & validator configuration
│── handlers/      → API request handling (business logic)
│── models/        → Database models (GORM)
│── request/       → Request DTOs
│── response/      → Response DTOs & pagination
│── routes/        → API route registration
│── messaging/     → RabbitMQ / Paota integration
│── main.go        → Application entry point

🔹 3. PostgreSQL Database Integration

Connected PostgreSQL using GORM

Auto-migration for User table

Soft delete support using gorm.DeletedAt

Email uniqueness enforcement

🔹 4. Secure Password Handling

Password hashing using bcrypt

Password never exposed in API responses

User responses return only:

id

name

email

🔹 5. Request Validation

Structured request validation using:

github.com/go-playground/validator/v10


Validation errors returned as proper HTTP responses

🔹 6. Pagination & Filtering

Pagination support:

page

limit

Filtering:

by name

by email

Pagination metadata returned in response

🔹 7. Soft Delete Implementation

Users are not permanently removed

Deleted users are excluded from:

Get All

Get by ID

Enables future recovery & audit safety

🔹 8. Event-Driven Architecture (RabbitMQ)

Events published after successful DB operations

Implemented two events:

USER_CREATED

USER_UPDATED

🔹 9. Paota-Based Messaging (Phase 2)

Replaced raw RabbitMQ implementation with Paota abstraction

Used Paota Publisher & WorkerPool

Messaging layer isolated (no business logic changes)

Durable exchange & queues

Persistent messages

Dead Letter Queue (DLQ) support

🔹 10. Environment-Based Configuration

All configs loaded from .env

No hardcoded credentials

Supports easy deployment & scaling

🔹 11. Production-Grade Code Quality

Clean error handling

Standard API responses

Modular design

Build-safe (go build ./... passes)

Ready for CI/CD pipelines

🚀 How to Run This Project
1️⃣ Prerequisites

Make sure you have installed:

Go 1.21+

PostgreSQL

RabbitMQ

2️⃣ Clone the Repository
git clone <repository-url>
cd user-crud-go

3️⃣ Configure Environment Variables

Create .env file:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=userdb

APP_PORT=8080

RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=user.events

4️⃣ Install Dependencies
go mod tidy

5️⃣ Build the Project (Recommended)
go build ./...


✔ Ensures entire project is error-free

6️⃣ Run the Application
go run main.go


Server starts on:

http://localhost:8080

📡 API Endpoints
➕ Create User
POST /users


Request Body

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}

📄 Get All Users
GET /users?page=1&limit=10

🔍 Get User by ID
GET /users/{id}

✏ Update User
PUT /users/{id}

❌ Delete User (Soft Delete)
DELETE /users/{id}

👀 What Output Should Be Observed
✅ API Responses

Proper JSON responses with status & message

Pagination metadata for list APIs

No password exposure

✅ Database

Users stored in PostgreSQL

Deleted users marked with deleted_at

No hard delete

✅ RabbitMQ

Events published on:

user.created

user.updated

Messages visible in queues

DLQ captures failed messages

