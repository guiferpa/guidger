# Guidger

📒 Guidger is a signed webhook service that handles balance updates via webhooks sent from an external system.

## Features

- ✅ **Signed Webhook Processing**: Validates HMAC SHA256 signatures
- ✅ **Replay Attack Prevention**: Uses nonce and timestamp validation
- ✅ **Decimal Precision**: Handles decimal amounts with exact precision using `mathutil.Big`
- ✅ **Thread-Safe**: Concurrent request handling with proper synchronization
- ✅ **Structured Logging**: Request IDs and structured JSON logs
- ✅ **Telemetry**: Request metrics and error tracking via `/metrics` endpoint
- ✅ **RESTful API**: Clean endpoints for webhook processing and balance queries

## Requirements

- Go 1.21 or higher
- `WEBHOOK_SECRET` environment variable (required)
- `PORT` environment variable (optional, defaults to 8080)

## Installation


### Using Make

The project includes a Makefile with convenient commands for development and CI/CD:

**Available Make targets:**

| Target | Description | Output |
|--------|-------------|--------|
| `all` | Run tests, lint, and build (default workflow) | Runs `test`, `lint`, and `build-force` |
| `server` | Build the server binary | Creates `./target/bin/server` |
| `test` | Run all tests with coverage | Uses `tparse` for formatted output |
| `lint` | Run golangci-lint on the codebase | Auto-installs linter if needed |
| `bench` | Run benchmarks | Tests with 1, 2, 4, and 12 CPUs |
| `cover-html` | Generate and open HTML coverage report | Opens coverage in browser |
| `clean` | Remove build artifacts | Removes `./target/bin` directory |

**Build Configuration:**
- Binary output: `./target/bin/server`
- Build flags: `CGO_ENABLED=0` (static binary), `-race` (race detector enabled)
- Test flags: `-v -json -race -buildvcs -cover`

**Dependencies:**
The Makefile automatically installs required tools:
- `golangci-lint`: For code linting
- `tparse`: For formatted test output
- `act`: For GitHub Actions workflow testing (optional)

## Configuration

The service uses environment variables for configuration:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `WEBHOOK_SECRET` | Yes | - | Secret key for HMAC SHA256 signature validation |
| `PORT` | No | `8080` | Port number for the HTTP server |

### Example

```bash
# Build the server
make server

# Run the server
export WEBHOOK_SECRET=your-secret-key-here
export PORT=8080
./target/bin/server
```

## API Endpoints

### POST /webhook

Processes a signed webhook to update user balances.

**Headers:**
- `X-Timestamp`: UNIX timestamp of when the request was created
- `X-Signature`: HMAC SHA256 signature of the request body
- `X-Nonce`: Unique nonce for each request

**Request Body:**
```json
{
    "user": "user123",
    "asset": "BTC",
    "amount": "1.5"
}
```

**Signature Calculation:**
```
canonical_string = X-Timestamp + "\n" + X-Nonce + "\n" + <raw_request_body>
signature = HMAC_SHA256(WEBHOOK_SECRET, canonical_string)
```

**Response:**
- `200 OK`: Webhook processed successfully
- `400 Bad Request`: Missing headers or invalid request body
- `401 Unauthorized`: Invalid signature, expired timestamp, or duplicate nonce
- `500 Internal Server Error`: Server error processing the request

**Example:**
```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: $(date +%s)" \
  -H "X-Nonce: $(uuidgen)" \
  -H "X-Signature: <calculated-signature>" \
  -d '{"user":"user123","asset":"BTC","amount":"1.5"}'
```

### GET /balance/{user}

Retrieves the balance for a specific user.

**Response:**
```json
{
    "user": "user123",
    "balances": {
        "BTC": "1.5",
        "ETH": "2.25"
    }
}
```

**Example:**
```bash
curl http://localhost:8080/balance/user123
```

### GET /metrics

Returns telemetry metrics for the service.

**Response:**
```json
{
    "total_requests": 100,
    "success_requests": 95,
    "error_requests": 5,
    "webhook_requests": 80,
    "balance_requests": 20,
    "error_counts": {
        "400": 2,
        "401": 3
    }
}
```

**Example:**
```bash
curl http://localhost:8080/metrics
```

## Security

- **HMAC SHA256 Signature Validation**: All webhooks must be signed
- **Timestamp Validation**: Requests must be within 5 minutes of server time
- **Nonce Validation**: Prevents replay attacks by tracking used nonces
- **Thread-Safe Operations**: All concurrent operations are properly synchronized

## Architecture

The project follows a clean architecture pattern:

```
guidger/
├── cmd/server/          # Application entry point
├── domain/              # Business logic
│   ├── ledger/         # Ledger domain logic
│   └── webhook/        # Webhook validation logic
├── infra/              # Infrastructure implementations
│   ├── http/server/    # HTTP handlers and middleware
│   ├── signer/hmac/    # HMAC signature implementation
│   └── storage/memory/ # In-memory storage
├── mathutil/           # Decimal precision utilities
├── envutil/            # Environment variable utilities
└── httputil/           # HTTP utility functions
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test ./... -cover
```

Run integration tests:
```bash
go test ./infra/http/server/... -v
```

## Decimal Precision

The service uses `mathutil.Big` to preserve exact decimal precision:
- Values are stored as `big.Rat` for exact mathematical operations
- Original precision (number of decimal places) is preserved
- Operations maintain the maximum precision between operands
- Example: `"1.0" + "1.00" = "2.00"` (preserves greater precision)

## Logging

The service uses structured JSON logging with request IDs:
- Each request gets a unique `X-Request-ID` header
- Logs include request context, duration, and status codes
- Errors are logged with full context for debugging
- Logs are output in JSON format for easy parsing by log aggregation tools

Example log entry:
```json
{
    "time": "2024-01-01T12:00:00Z",
    "level": "INFO",
    "msg": "http_request",
    "request_id": "abc123...",
    "method": "POST",
    "path": "/webhook",
    "status": 200,
    "duration_ms": 15,
    "bytes_written": 45,
    "remote_addr": "127.0.0.1:12345"
}
```

## Telemetry

The service provides telemetry metrics via the `/metrics` endpoint:
- **Total Requests**: Total number of HTTP requests processed
- **Success Requests**: Number of successful requests (status 2xx, 3xx)
- **Error Requests**: Number of failed requests (status 4xx, 5xx)
- **Endpoint Metrics**: Separate counters for webhook and balance endpoints
- **Error Breakdown**: Count of errors by HTTP status code

All metrics are tracked in-memory and reset when the server restarts.

## Usage Examples

### Complete Webhook Flow

1. Generate a signature for your webhook request:
```bash
# Set your secret
export WEBHOOK_SECRET="your-secret-key"

# Prepare request data
TIMESTAMP=$(date +%s)
NONCE=$(uuidgen)
PAYLOAD='{"user":"user123","asset":"BTC","amount":"1.5"}'

# Generate canonical string
CANONICAL="${TIMESTAMP}\n${NONCE}\n${PAYLOAD}"

# Generate HMAC SHA256 signature (using openssl)
SIGNATURE=$(echo -n "$CANONICAL" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | cut -d' ' -f2)

# Send webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: $TIMESTAMP" \
  -H "X-Nonce: $NONCE" \
  -H "X-Signature: $SIGNATURE" \
  -d "$PAYLOAD"
```

2. Check user balance:
```bash
curl http://localhost:8080/balance/user123
```

3. Monitor service metrics:
```bash
curl http://localhost:8080/metrics | jq
```

## License

See [LICENSE](LICENSE) file for details.
