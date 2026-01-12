How to Run the Project
1. Start infrastructure
docker compose up -d


Starts PostgreSQL and RabbitMQ.

2. Run API (User CRUD + Event Producer)

Windows (PowerShell)

$env:APP_MODE="api"
go run main.go


macOS / Linux

export APP_MODE=api
go run main.go


API runs on:

http://localhost:8080

3. Run USER_CREATED Consumer

Open a new terminal.

Windows

$env:APP_MODE="consumer_created"
go run main.go


macOS / Linux

export APP_MODE=consumer_created
go run main.go


Handles USER_CREATED events.

4. Run USER_UPDATED Consumer

Open another terminal.

Windows

$env:APP_MODE="consumer_updated"
go run main.go


macOS / Linux

export APP_MODE=consumer_updated
go run main.go


Handles USER_UPDATED events.

5. Verify

Create user → USER_CREATED consumer logs

Update user → USER_UPDATED consumer logs

Queues visible in RabbitMQ UI (localhost:15672)

APP_MODE values

api – REST API and producer

consumer_created – USER_CREATED consumer

consumer_updated – USER_UPDATED consumer