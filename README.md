# GoLang REST Project

A Go learning project demonstrating how to build REST APIs using only the Go standard library (`net/http`) — no frameworks — and how two Go servers communicate with each other over HTTP.

---

## Project Structure

```
GoLang_REST_Project1/
├── Server1User/          # User CRUD API (in-memory store, port 8080)
│   ├── main.go
│   ├── go.mod
│   ├── models/           # User struct, request/response types
│   ├── store/            # In-memory data layer (sync.RWMutex)
│   ├── handlers/         # HTTP handlers (controller layer)
│   ├── router/           # URL routing via net/http ServeMux
│   └── UserCRUD.postman_collection.json
│
└── Server2Admin/         # Admin API that consumes Server1User (port 8081)
    ├── main.go
    ├── go.mod
    ├── models/           # Mirrors Server1User response shapes
    ├── client/           # HTTP client to call Server1User
    ├── handlers/         # Admin handlers using the client
    ├── router/           # URL routing via net/http ServeMux
    └── AdminCRUD.postman_collection.json
```

---

## Architecture

```
Postman / curl
      │
      ▼ :8081
┌─────────────────┐         HTTP call         ┌─────────────────┐
│   Server2Admin  │ ───────────────────────▶  │   Server1User   │
│                 │                           │                 │
│  router/        │         JSON response     │  router/        │
│  handlers/      │ ◀───────────────────────  │  handlers/      │
│  client/        │                           │  store/         │
└─────────────────┘                           └─────────────────┘
      :8081                                         :8080
```

**Key concept:** Server2Admin uses Go's `net/http.Client` to make HTTP calls to Server1User. It acts as a proxy/consumer — it has no database of its own.

---

## Server1User

**Port:** `8080`  
**Purpose:** Manages users with full CRUD using an in-memory map.

### API Endpoints

| Method | Endpoint | Description | Success | Error |
|--------|----------|-------------|---------|-------|
| `POST` | `/users` | Create a user | `201` | `400` bad input, `409` duplicate email |
| `GET` | `/users` | List all users | `200` | — |
| `GET` | `/users/{id}` | Get user by ID | `200` | `404` not found |
| `PUT` | `/users/{id}` | Update user | `200` | `404` not found, `409` duplicate email |
| `DELETE` | `/users/{id}` | Delete user | `200` | `404` not found |

### Request Body (Create / Update)

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "age": 30
}
```

### Response (User)

```json
{
  "id": "1",
  "name": "Alice",
  "email": "alice@example.com",
  "age": 30,
  "created_at": "2026-07-01T12:00:00Z",
  "updated_at": "2026-07-01T12:00:00Z"
}
```

### Run

```bash
cd Server1User
go run main.go

# Custom port
PORT=9000 go run main.go
```

---

## Server2Admin

**Port:** `8081`  
**Purpose:** Admin server that forwards all requests to Server1User via HTTP client.

### API Endpoints

| Method | Endpoint | Forwards To | Description |
|--------|----------|-------------|-------------|
| `POST` | `/admin/users` | `POST /users` | Create a user |
| `GET` | `/admin/users` | `GET /users` | List all users |
| `GET` | `/admin/users/{id}` | `GET /users/{id}` | Get user by ID |
| `PUT` | `/admin/users/{id}` | `PUT /users/{id}` | Update user |
| `DELETE` | `/admin/users/{id}` | `DELETE /users/{id}` | Delete user |

### Error Propagation

| Scenario | HTTP Status |
|----------|-------------|
| Server1User returns 404 | Admin returns `404` |
| Server1User returns 409 | Admin returns `409` |
| Server1User is down / unreachable | Admin returns `502 Bad Gateway` |
| Invalid input (caught before calling Server1User) | Admin returns `400` |

### Run

```bash
cd Server2Admin
go run main.go

# Custom port or custom Server1User URL
PORT=9001 SERVER1_URL=http://localhost:9000 go run main.go
```

---

## Running Both Servers Together

Open two terminals:

```bash
# Terminal 1 — Server1User
cd Server1User
go run main.go
# Server starting on http://localhost:8080

# Terminal 2 — Server2Admin
cd Server2Admin
go run main.go
# Server2Admin starting on http://localhost:8081
# Forwarding requests to Server1User at http://localhost:8080
```

Then hit Server2Admin and watch it talk to Server1User:

```bash
# Create a user through the Admin server
curl -X POST http://localhost:8081/admin/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":30}'

# Verify directly on Server1User — same data!
curl http://localhost:8080/users
```

---

## Postman Collections

| File | Server | Base URL |
|------|--------|----------|
| `Server1User/UserCRUD.postman_collection.json` | Server1User | `http://localhost:8080` |
| `Server2Admin/AdminCRUD.postman_collection.json` | Server2Admin | `http://localhost:8081` |

**Import:** Open Postman → Import → drag & drop the JSON file.

---

## Key Concepts Demonstrated

### 1. In-Memory Store with Thread Safety (`Server1User/store/`)
```go
type UserStore struct {
    mu    sync.RWMutex       // allows many readers, one writer
    users map[string]*User   // the "database"
}
```

### 2. HTTP Client between Servers (`Server2Admin/client/`)
```go
// Simple GET
resp, err := c.httpClient.Get(baseURL + "/users")

// PUT / DELETE — use NewRequest for custom methods
req, _ := http.NewRequest(http.MethodPut, url, body)
req.Header.Set("Content-Type", "application/json")
resp, err := c.httpClient.Do(req)
```

### 3. Layered Architecture (both servers)
```
Request → router → handler → store/client → Response
```

---

## Tech Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.22+ | Language |
| `net/http` | stdlib | HTTP server & client |
| `encoding/json` | stdlib | JSON encode/decode |
| `sync` | stdlib | Thread-safe map access |

> **No external dependencies.** Both servers use only the Go standard library.

---

## Build for Production

```bash
# Build Server1User binary
cd Server1User
go build -o server1user main.go
./server1user

# Build Server2Admin binary
cd Server2Admin
go build -o server2admin main.go
./server2admin
```

The output binaries are self-contained — no Go installation needed on the target machine.
