# Go To-Do REST API

---

## Project Structure

```
todo-app/
├── go.mod
├── main.go
├── models/
│   └── todo.go
├── store/
│   ├── store.go
│   └── store_test.go
└── handlers/
    ├── todo.go
    └── todo_test.go
```

---

## Prerequisites

| Tool    | Minimum version | Download                       |
| ------- | --------------- | ------------------------------ |
| Go      | 1.21            | https://go.dev/dl/             |
| VS Code | any             | https://code.visualstudio.com/ |

> **Tip (VS Code):** Install the official **Go** extension (`golang.go`) for IntelliSense, formatting, and test integration.

---

## Running the Server

```bash
# 1. Clone / copy the project
cd todo-app

# 2. Start the server (Windows PowerShell or Command Prompt)
go run main.go
```

You should see:

```
To-Do API server running on http://localhost:8080
Press Ctrl+C to stop.
```

---

## API Reference

All endpoints accept and return **JSON**. The `due_date` field always uses the format `YYYY-MM-DD`.

### Data Model

```json
{
  "id": "20240501120000.123456789",
  "task": "Buy groceries",
  "due_date": "2024-06-01T00:00:00Z",
  "status": "pending",
  "created_at": "2024-05-01T12:00:00Z",
  "updated_at": "2024-05-01T12:00:00Z"
}
```

| Field        | Type     | Values                               |
| ------------ | -------- | ------------------------------------ |
| `id`         | string   | Auto-generated                       |
| `task`       | string   | Free text                            |
| `due_date`   | datetime | ISO 8601 (set from YYYY-MM-DD input) |
| `status`     | string   | `pending` \| `completed`             |
| `created_at` | datetime | Set on creation                      |
| `updated_at` | datetime | Updated on every change              |

---

### 1. Create a To-Do Item

**POST** `/todos`

**Request body:**

```json
{
  "task": "Buy groceries",
  "due_date": "2024-06-01"
}
```

**Response `201 Created`:**

```json
{
  "id": "20240501120000.123456789",
  "task": "Buy groceries",
  "due_date": "2024-06-01T00:00:00Z",
  "status": "pending",
  "created_at": "2024-05-01T12:00:00Z",
  "updated_at": "2024-05-01T12:00:00Z"
}
```

**cURL example:**

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d "{\"task\": \"Buy groceries\", \"due_date\": \"2024-06-01\"}"
```

---

### 2. Get a To-Do Item

**GET** `/todos/{id}`

**Response `200 OK`:**

```json
{
  "id":       "20240501120000.123456789",
  "task":     "Buy groceries",
  "due_date": "2024-06-01T00:00:00Z",
  "status":   "pending",
  ...
}
```

**cURL example:**

```bash
curl http://localhost:8080/todos/20240501120000.123456789
```

---

### 3. List To-Do Items

**GET** `/todos`

Returns items sorted by `due_date` **ascending** (earliest first). Completed items are **excluded** by default.

| Query parameter     | Type    | Default | Description                            |
| ------------------- | ------- | ------- | -------------------------------------- |
| `include_completed` | boolean | `false` | Pass `true` to include completed items |

**Response `200 OK`:**

```json
[
  { "id": "...", "task": "Earlier task", "due_date": "2024-05-10T00:00:00Z", "status": "pending", ... },
  { "id": "...", "task": "Later task",   "due_date": "2024-06-01T00:00:00Z", "status": "pending", ... }
]
```

**cURL examples:**

```bash
# Pending tasks only (default)
curl http://localhost:8080/todos

# All tasks including completed
curl "http://localhost:8080/todos?include_completed=true"
```

---

### 4. Update a To-Do Item

**PUT** `/todos/{id}`

All fields in the request body are **optional** — only provided fields are updated.

**Request body (any combination):**

```json
{
  "task": "Updated task text",
  "due_date": "2024-07-15",
  "status": "completed"
}
```

**Response `200 OK`:** Returns the updated item.

**cURL example:**

```bash
# Mark an item as completed
curl -X PUT http://localhost:8080/todos/20240501120000.123456789 \
  -H "Content-Type: application/json" \
  -d "{\"status\": \"completed\"}"
```

---

### 5. Delete a To-Do Item

**DELETE** `/todos/{id}`

**Response `200 OK`:**

```json
{ "message": "todo deleted successfully" }
```

**cURL example:**

```bash
curl -X DELETE http://localhost:8080/todos/20240501120000.123456789
```

---

## To Run on Windows

```
# 1. Create
$item = Invoke-RestMethod -Uri "http://localhost:8080/todos" -Method POST -ContentType "application/json" -Body '{"task":"Buy milk","due_date":"2026-05-01"}'
$item

# 2. Get the ID automatically
$id = $item.id
Write-Host "Created item with ID: $id"

# 3. Get it
Invoke-RestMethod -Uri "http://localhost:8080/todos/$id" -Method GET

# 4. Update it
Invoke-RestMethod -Uri "http://localhost:8080/todos/$id" -Method PUT -ContentType "application/json" -Body '{"status":"completed"}'

# 5. List all (completed item won't show here)
Invoke-RestMethod -Uri "http://localhost:8080/todos" -Method GET

# 6. List including completed
Invoke-RestMethod -Uri "http://localhost:8080/todos?include_completed=true" -Method GET

# 7. Delete it
Invoke-RestMethod -Uri "http://localhost:8080/todos/$id" -Method DELETE

# Confirm it's deleted
Invoke-RestMethod -Uri "http://localhost:8080/todos/$id" -Method GET
```

---

---

## Error Responses

All errors return a JSON object with an `error` key:

```json
{ "error": "todo not found" }
```

| HTTP Status | Meaning                                  |
| ----------- | ---------------------------------------- |
| `400`       | Bad request — invalid body or parameters |
| `404`       | Item not found                           |
| `405`       | Method not allowed on that endpoint      |
| `500`       | Internal server error                    |

---
