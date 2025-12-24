User CRUD API with Soft Delete (Go + Echo + GORM)
Overview

This project implements a User CRUD REST API using Golang, Echo, GORM, and PostgreSQL, following clean architecture and production-ready best practices.

The implementation incorporates all review feedback to ensure:

Proper separation of concerns

Secure handling of sensitive data

Correct soft delete behavior

Explicit error handling

Clean API contracts

Setup Instructions
1. Prerequisites

Ensure the following are installed on your system:

Go (v1.20 or higher)

PostgreSQL

Git

2. Clone the Repository
git clone https://github.com/Purviksiprtc/go-backend-learning.git
cd go-backend-learning

3. Configure Environment Variables

Create a .env file in the project root with the following content:

DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_backend
DB_PORT=5432
APP_PORT=8080


⚠️ Update the database credentials based on your local setup.

4. Install Dependencies
go mod tidy

5. Run Database Migrations

Database tables are automatically created using GORM when the application starts.

6. Start the Application
go run main.go


Server will start on:

http://localhost:8080

API Endpoints
Method	Endpoint	Description
POST	/users	Create a new user
GET	/users	Get all active users
GET	/users/:id	Get a user by ID
PUT	/users/:id	Update a user
DELETE	/users/:id	Soft delete a user
Key Improvements Based on PR Review
1. Separation of API Contracts and Database Models

Introduced dedicated request DTOs under /request

Introduced dedicated response DTOs under /response

API handlers do not bind or return database models directly

Explicit mapping performed:

Request → Model

Model → Response

2. Proper Soft Delete Implementation

Implemented soft delete using gorm.DeletedAt

All read operations filter records using:

deleted_at IS NULL


Delete operations mark records as deleted instead of removing them

Ensures deleted users never appear in API responses

3. Email Uniqueness Handled in Application Logic

Removed database-level unique constraint

Email uniqueness enforced manually:

Create User checks for existing active users

Update User checks excluding the current user

Allows reuse of emails after soft deletion

4. Secure Password Handling

Passwords are hashed using bcrypt

Bcrypt errors are explicitly handled

Password field removed from API responses

Password is never exposed in any response

5. Explicit Database Error Handling

All database operations handle .Error

API returns appropriate HTTP errors

Prevents silent failures and misleading success responses

6. Environment-Based Configuration

Server port is loaded from .env

Improves portability across environments

Folder Structure
user-crud-go/
├── config/        # Database configuration
├── handlers/      # HTTP handlers
├── models/        # Database models
├── request/       # Request DTOs
├── response/      # Response DTOs
├── routes/        # Route registration
├── main.go        # Application entry point
├── go.mod
├── go.sum 

