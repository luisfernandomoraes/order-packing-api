# Order Packing Calculator API

An intelligent order packing API built in Go that finds the optimal combination of package sizes to fulfill any order, minimizing shipping waste and number of packages.

## 🚀 Quick Start

```bash
git clone https://github.com/luisfernandomoraes/order-packing-api
cd order-packing-api
make run
```

Or using Docker:

```bash
make build-container
make run-container
```

Access:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html

## 🌐 Live Demo

Try the application online without installation:

**https://order-packing-api.onrender.com/**

- **Web Interface**: https://order-packing-api.onrender.com/
- **Swagger API Documentation**: https://order-packing-api.onrender.com/swagger/index.html

## 📋 Table of Contents

- [Quick Start](#🚀-quick-start-)
- [About the Project](#🎯-about-the-project)
- [Calculation Logic](#🧮-calculation-logic)
- [Project Structure](#📁-project-structure)
- [Middlewares](#🔧-middlewares)
- [API Documentation](#📚-api-documentation)
- [Endpoints](#🚀-endpoints)
- [How to Run](#🏃-how-to-run)
- [Testing](#🧪-testing)

## 🎯 About the Project

This project implements an optimized package calculation system to fulfill customer orders. The goal is to ship the minimum number of items (minimizing waste) using the fewest packages possible.

### Business Rules

1. **Only whole packs**: Packages cannot be broken
2. **Minimize items shipped**: Prioritize combinations that send fewer total items
3. **Minimize number of packages**: Within the previous constraint, use fewer packages

### Example

For an order of 501 items with package sizes [250, 500, 1000]:

- ❌ 3x250 = 750 items (3 packages)
- ❌ 1x1000 = 1000 items (1 package, but too much waste)
- ✅ 1x500 + 1x250 = 750 items (2 packages) ← Optimal solution

## 🧮 Calculation Logic

The algorithm uses **Dynamic Programming** to find the optimal solution.

### How it Works

```go
// Candidate solution structure
type solution struct {
    totalItems     int         // Total items in this combination
    packsBySize    map[int]int // Quantity of each package size
    totalPackCount int         // Total number of packages
}
```

## Algorithm: Dynamic Programming

1. **Dynamic Programming Table**: Creates a table where `dp[i]` represents the best solution for `i` items
2. **Bottom-Up Construction**: For each quantity from 1 to `order + largestPack`:
   - Try adding each available package size
   - Compare with the current best solution for that quantity
   - Keep only the best solution (fewer items, then fewer packages)
3. **Solution Search**: Find the first quantity >= order that has a valid solution

### Comparison Criteria

```go
func isBetterSolution(new, current *solution) bool {
    // Priority 1: Fewer total items
    if new.totalItems < current.totalItems {
        return true
    }

    // Priority 2: Fewer packages (with same total items)
    if new.totalItems == current.totalItems &&
       new.totalPackCount < current.totalPackCount {
        return true
    }

    return false
}
```

### Complexity

- **Time**: O(n × m), where n = order + largestPack and m = number of sizes
- **Space**: O(n) to store optimal solutions

## 📁 Project Structure

```sh
order-packing-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Application configuration
│   ├── domain/
│   │   ├── pack_calculator.go     # Core business logic
│   │   └── pack_calculator_test.go # Business logic tests
│   ├── handlers/
│   │   ├── health.go              # Health check handler
│   │   ├── health_test.go
│   │   ├── calculate.go           # Package calculation handler
│   │   ├── calculate_test.go
│   │   ├── pack_sizes.go          # Pack sizes management handler
│   │   └── pack_sizes_test.go
│   ├── middleware/
│   │   ├── chain.go               # Middleware chaining
│   │   ├── cors.go                # CORS headers
│   │   ├── logging.go             # Request logging
│   │   └── recovery.go            # Panic recovery
│   ├── response/
│   │   └── json.go                # JSON response utilities
│   └── server/
│       ├── routes.go              # Route definitions
│       ├── server.go              # HTTP server configuration
│       └── server_test.go         # Integration tests
├── static/
│   └── index.html                 # Web interface (UI)
├── Makefile                       # Build and test commands
├── go.mod                         # Project dependencies
└── README.md                      # This file
```

### Application Layers

- **cmd/**: Application entry points
- **internal/domain/**: Pure business logic (calculation algorithm)
- **internal/handlers/**: HTTP handlers (presentation layer)
- **internal/middleware/**: Reusable HTTP middlewares
- **internal/response/**: HTTP response utilities
- **internal/server/**: Server configuration and setup
- **static/**: Static files (UI)

## 🔧 Middlewares

The application uses a chain middleware architecture:

### 1. **CORS** (`middleware/cors.go`)
```go
// Adds CORS headers to allow cross-origin requests
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

**Responsibilities**:

- Allows requests from any origin
- Supports GET, POST, OPTIONS methods
- Handles preflight requests (OPTIONS)

### 2. **Logging** (`middleware/logging.go`)

```go
// Logs information about each request
log.Printf("%s %s %d %v %s", method, path, statusCode, duration, remoteAddr)
```

**Responsibilities**:

- Logs method, path, status code
- Measures response time
- Logs client remote address

### 3. **Recovery** (`middleware/recovery.go`)

```go
// Recovers from panics and returns 500 error
defer func() {
    if err := recover(); err != nil {
        log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
        http.Error(w, "Internal server error", 500)
    }
}()
```

**Responsibilities**:

- Captures unhandled panics
- Logs complete stack trace
- Returns appropriate error response
- Keeps server running after errors

### Middleware Chain

```go
handler := middleware.Chain(
    finalHandler,
    middleware.CORS,      // 1st: Adds CORS headers
    middleware.Logging,   // 2nd: Logs request
    middleware.Recovery,  // 3rd: Catches panics (innermost)
)
```

Order matters: Recovery must be innermost to catch errors from all others.

## 📚 API Documentation

The API is fully documented using **Swagger/OpenAPI 3.0**.

### Access Swagger UI

Once the server is running, access the interactive API documentation at:

**http://localhost:8080/swagger/index.html**

The Swagger UI provides:

- ✅ Complete API specification
- ✅ Interactive endpoint testing
- ✅ Request/response examples
- ✅ Schema definitions

### Generate Swagger Docs

To regenerate the Swagger documentation after making changes to the API:

```bash
make swagger
```

This will:

1. Install `swag` CLI if not already installed
2. Parse annotations from code
3. Generate `docs/swagger.json`, `docs/swagger.yaml`, and `docs/docs.go`

## 🚀 Endpoints

### Health Check

**GET** `/health`

Checks if the API is running.

**Response**:

```json
{
  "status": "healthy",
  "app": "Order Packing Calculator API"
}
```

---

### Calculate Packages

**POST** `/api/calculate`

Calculates the best package combination for an order.

**Request Body**:
```json
{
  "order": 501
}
```

**Response**:
```json
{
  "order": 501,
  "total_items": 750,
  "packs": {
    "250": 1,
    "500": 1
  },
  "pack_sizes": [250, 500, 1000, 2000, 5000],
  "surplus": 249,
  "total_packs": 2
}
```

**Validations**:
- ❌ `order < 0`: Returns 400 "Order must be positive"
- ❌ Invalid JSON: Returns 400 "Invalid request body"

---

### Get Package Sizes

**GET** `/api/pack-sizes`

Returns the configured package sizes.

**Response**:
```json
{
  "pack_sizes": [250, 500, 1000, 2000, 5000]
}
```

---

### Update Package Sizes

**POST** `/api/pack-sizes`

Updates the available package sizes.

**Request Body**:
```json
{
  "pack_sizes": [100, 250, 500, 1000]
}
```

**Response**:
```json
{
  "message": "Pack sizes updated successfully",
  "pack_sizes": [100, 250, 500, 1000]
}
```

**Validations**:
- ❌ Empty array: Returns 400 "Pack sizes cannot be empty"
- ❌ Negative or zero values: Returns 400 "All pack sizes must be positive"

---

### Web Interface

**GET** `/`

Serves the interactive web interface to use the API.

## 🏃 How to Run

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)

### Using Make (Recommended)

#### 1. Run the application
```bash
make run
```

The API will be available at `http://localhost:8080`

#### 2. Run in development mode (with auto-reload)
```bash
make dev
```

#### 3. Run tests
```bash
make test
```

#### 4. Run tests with coverage
```bash
make test-coverage
```

#### 5. Build the application
```bash
make build
```

The binary will be created at `./bin/api`

#### 6. Clean generated files
```bash
make clean
```

### Available Makefile Commands

```makefile
make build          # Compile the application
make run            # Run the application
make dev            # Run with auto-reload
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make clean          # Remove generated files
make help           # Show help
```

### Environment Variables

```bash
# Server port (default: 8080)
PORT=8080

# Default package sizes (default: 250,500,1000,2000,5000)
DEFAULT_PACK_SIZES=250,500,1000,2000,5000
```

## 🧪 Testing

The project has complete test coverage at three levels:

### 1. Unit Tests

**Domain Layer** (`internal/domain/pack_calculator_test.go`)
- ✅ Basic cases (exact order, order with surplus, zero order)
- ✅ Edge cases (prime numbers, coprimes, large packages)
- ✅ Business rules (item and package minimization)
- ✅ Branch coverage (empty pack sizes, impossible solutions)

**Handlers** (`internal/handlers/*_test.go`)
- ✅ Input validations (invalid JSON, negative values)
- ✅ HTTP methods (GET, POST, other methods)
- ✅ Responses and JSON formats
- ✅ Concurrency tests

### 2. Integration Tests

**Server** (`internal/server/server_test.go`)
- ✅ End-to-end endpoints
- ✅ CORS headers
- ✅ Middlewares (logging, recovery)
- ✅ Static file serving

### 3. Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Only one package
go test ./internal/domain -v

# With race detector
go test -race ./...

# Benchmarks
go test -bench=. ./internal/domain
```

### Test Coverage

- **Domain**: ~95% coverage
- **Handlers**: ~90% coverage
- **Integration**: ~85% coverage
- **Total**: ~75 tests passing

## 📊 Usage Examples

### Using cURL

```bash
# Health check
curl http://localhost:8080/health

# Calculate packages
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"order": 501}'

# Get sizes
curl http://localhost:8080/api/pack-sizes

# Update sizes
curl -X POST http://localhost:8080/api/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{"pack_sizes": [100, 250, 500, 1000]}'
```

### Using the Web Interface

1. Access `http://localhost:8080`
2. Configure the desired package sizes
3. Enter the order quantity
4. Click "Calculate"
5. See the result with the optimal package distribution

## 🛠️ Technologies Used

- **Go 1.25+**: Programming language
- **net/http**: Native HTTP server
- **testify**: Testing framework and assertions
- **Standard Library**: Only Go standard libraries (no external frameworks)

## 📈 Performance

- **Algorithm**: O(n × m) where n ≈ order size, m = number of sizes
- **Memory**: O(n) for dynamic programming table
- **Concurrency**: Thread-safe with `sync.RWMutex` for pack sizes read/write

### Benchmarks (Macbook Pro M1, Go 1.25)

```
BenchmarkCalculate_SmallOrder-8     50000    25.3 µs/op
BenchmarkCalculate_MediumOrder-8    10000   120.5 µs/op
BenchmarkCalculate_LargeOrder-8      2000   850.2 µs/op
```

## 👨‍💻 Author

Luís Fernando Moraes