

# User CRUD API with RabbitMQ (Event-Driven Backend)

## 📌 Overview

This project is a **backend-only Go application** that demonstrates:

* User CRUD operations (Create, Read, Update, Delete)
* **Event-driven architecture using RabbitMQ**
* Asynchronous communication between producer and consumer
* Reliable message handling using **Ack / Nack and Dead Letter Queue**
* Fully **Dockerized setup** (no local installations needed except Docker)

When a user is **created or updated**, an **event is published** to RabbitMQ.
A **consumer listens** to these events and simulates:

* Sending welcome emails
* Creating audit logs

---

## 🧠 What You Will Learn from This Project

* How RabbitMQ works in real backend systems
* How producers and consumers communicate
* Why events are better than direct service calls
* How reliable message processing is designed
* How to observe everything using RabbitMQ Management UI

---

## 🏗️ Architecture (Simple Explanation)

```
Client (Postman / Browser)
        |
        v
User CRUD API (Producer)
        |
        v
RabbitMQ (Exchange → Queue)
        |
        v
Consumer (Notification & Audit)
```

---

## 🧩 Features Implemented

### User CRUD

* Create user
* Get users
* Update user
* Soft delete user
* Password hashing with bcrypt

### RabbitMQ (Event Driven)

* USER_CREATED event
* USER_UPDATED event
* Topic exchange
* Separate queues per event
* Manual acknowledgements
* Dead Letter Queue (DLQ)

---

## 📂 Project Structure

```
user-crud-go/
├── config/          # Database config
├── handlers/        # User CRUD handlers
├── models/          # Database models
├── request/         # Request DTOs
├── response/        # Response DTOs
├── routes/          # API routes
├── rabbitmq/        # RabbitMQ connection, setup, publisher
├── consumer/        # RabbitMQ consumer (Notification & Audit)
├── main.go          # Application entry point
├── docker-compose.yml
├── Dockerfile
├── .env
├── example.env
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Prerequisites

You only need:

* **Docker**
* **Docker Compose**

👉 No need to install Go, PostgreSQL, or RabbitMQ locally.

---

## 🚀 How to Run This Project (Step by Step)

### 1️⃣ Clone the Repository

```bash
git clone <your-repository-url>
cd user-crud-go
```

---

### 2️⃣ Environment Configuration

Create a `.env` file in project root:

```env
APP_PORT=8080

DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_backend
DB_PORT=5432

RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

RABBITMQ_EXCHANGE=user.events
RABBITMQ_USER_CREATED_QUEUE=user.created.queue
RABBITMQ_USER_UPDATED_QUEUE=user.updated.queue
RABBITMQ_DLQ=user.dlq
```

(Use `example.env` as reference)

---

### 3️⃣ Start the Application

```bash
docker compose up -d --build
```

This starts:

* Go application
* PostgreSQL
* RabbitMQ (with management UI)

---

## 🌐 Access URLs

| Service        | URL                                              |
| -------------- | ------------------------------------------------ |
| API            | [http://localhost:8080](http://localhost:8080)   |
| RabbitMQ UI    | [http://localhost:15672](http://localhost:15672) |
| RabbitMQ Login | guest / guest                                    |

---

## 🧪 How to Test & Observe (VERY IMPORTANT)

### 🔹 Step 1: Observe RabbitMQ BEFORE API call

1. Open RabbitMQ UI
   👉 [http://localhost:15672](http://localhost:15672)
2. Go to **Queues**

You should see:

* `user.created.queue` → 0 messages
* `user.updated.queue` → 0 messages
* `user.dlq` → 0 messages

This is the **initial state**.

---

### 🔹 Step 2: Create User (USER_CREATED event)

Send request:

```http
POST http://localhost:8080/users
```

```json
{
  "name": "Rahul",
  "email": "rahul@test.com",
  "password": "123"
}
```

---

### 🔹 Step 3: Observe Logs (MOST IMPORTANT)

```bash
docker logs -f user-crud-app
```

You should see:

```
[USER_CREATED] Welcome email sent to rahul@test.com
[USER_CREATED] Audit log created for user 1
```

✔ This proves RabbitMQ is working.

---

### 🔹 Step 4: Observe RabbitMQ UI

* `user.created.queue` briefly receives message
* Consumer ACKs it
* Queue count returns to **0**

👉 **0 messages = SUCCESS**

---

### 🔹 Step 5: Update User (USER_UPDATED event)

```http
PUT http://localhost:8080/users/1
```

```json
{
  "name": "Rahul Sharma",
  "email": "rahul@test.com"
}
```

Logs:

```
[USER_UPDATED] User 1 profile updated
[USER_UPDATED] Audit log updated for rahul@test.com
```

---

### 🔹 Step 6: Dead Letter Queue (DLQ)

* Normally `user.dlq` = 0
* Messages go to DLQ only if consumer fails
* Used for debugging & reliability

---

## ❓ Why Queue Shows 0 Messages?

Because:

* Consumer is running
* Messages are processed immediately
* ACK deletes message

👉 This is **expected and correct behavior**

---

## 🧠 Important Concepts (Simple)

* **Producer**: User CRUD API (publishes events)
* **RabbitMQ**: Message broker
* **Consumer**: Notification & Audit service
* **Ack**: Message processed successfully
* **Nack**: Message failed → sent to DLQ

---

## 🎯 Interview-Ready Explanation

**Q: Why producer and consumer are in same service?**

> For learning purposes. In real production, they would be separate services, but the event-driven contract remains the same.

**Q: How do you verify RabbitMQ is working?**

> By observing consumer logs and message acknowledgements, not by queue message count.

---

## 🛑 How to Stop the Application

```bash
docker compose down
```

---

## ✅ Final Notes

* Backend-only project
* No frontend included
* Clean, production-style RabbitMQ usage
* Suitable for interviews and real-world learning

---

## 🎉 Conclusion

This project demonstrates a **real event-driven backend system** using:

* Go
* RabbitMQ
* PostgreSQL
* Docker

✔ Not beginner-level
✔ Industry-relevant
✔ Interview-ready

---

If you want, I can also provide:

* **Architecture diagram (text)**
* **Interview Q&A**
* **Troubleshooting guide**

Just tell me 👍
