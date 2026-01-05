# Guidger

📒 Guidger is a signed webhook service that handles balance updates via webhooks sent from an external system.

## Features

- ✅ **Signed Webhook Processing**: Validates HMAC SHA256 signatures
- ✅ **Replay Attack Prevention**: Uses nonce and timestamp validation
- ✅ **Decimal Precision**: Handles decimal amounts with exact precision using `mathutil.Big`
- ✅ **Thread-Safe**: Concurrent request handling with proper synchronization
- ✅ **Structured Logging**: Request IDs and structured JSON logs
- ✅ **RESTful API**: Clean endpoints for webhook processing and balance queries

## Requirements

- Go 1.21 or higher
- `WEBHOOK_SECRET` environment variable (required)
- `PORT` environment variable (optional, defaults to 8080)

## Installation

```bash
git clone https://github.com/guiferpa/guidger.git
cd guidger
go build -o bin/server ./cmd/server
```

## Configuration

The service uses environment variables for configuration:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `WEBHOOK_SECRET` | Yes | - | Secret key for HMAC SHA256 signature validation |
| `PORT` | No | `8080` | Port number for the HTTP server |

### Example

```bash
export WEBHOOK_SECRET=your-secret-key-here
export PORT=8080
./bin/server
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

## License

See [LICENSE](LICENSE) file for details.
