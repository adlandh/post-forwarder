# Post-Forwarder

Post-Forwarder is a webhook forwarding service that receives webhook requests and forwards them to various notification services like Telegram, Slack, and Pushover.

## Features

- Receive webhook requests via HTTP GET and POST methods
- Forward messages to multiple notification services:
  - Telegram
  - Slack
  - Pushover
- Handle long messages by storing them in Redis and generating a URL to view them
- Secure webhook endpoints with token authentication
- Sentry integration for error tracking and monitoring
- Kubernetes deployment ready

## Installation
```bash
go install github.com/adlandh/post-forwarder@latest
```
or
```bash
docker pull ghcr.io/adlandh/post-forwarder/post-forwarder
```

## Configuration

The application is configured through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `AUTH_TOKEN` | Authentication token for webhook endpoints | Required |
| `NOTIFIERS` | Comma-separated list of enabled notifiers (TELEGRAM, SLACK, PUSHOVER) | `TELEGRAM` |
| `REDIS_URL` | Redis server URL | Required |
| `REDIS_PREFIX` | Prefix for Redis keys | `post-forwarder` |
| `TELEGRAM_TOKEN` | Telegram bot token | Required if TELEGRAM notifier is enabled |
| `TELEGRAM_CHAT_IDS` | Comma-separated list of Telegram chat IDs | Required if TELEGRAM notifier is enabled |
| `SLACK_TOKEN` | Slack API token | Required if SLACK notifier is enabled |
| `SLACK_CHANNEL_IDS` | Comma-separated list of Slack channel IDs | Required if SLACK notifier is enabled |
| `PUSHOVER_TOKEN` | Pushover application API token | Required if PUSHOVER notifier is enabled |
| `PUSHOVER_USER` | Pushover user API token | Required if PUSHOVER notifier is enabled |
| `SENTRY_DSN` | Sentry DSN for error tracking | Optional |
| `SENTRY_ENVIRONMENT` | Sentry environment | Optional |
| `SENTRY_TRACES_SAMPLE_RATE` | Sentry traces sample rate | `1.0` |

## API Documentation

The API is defined using OpenAPI 3.0 specification in `api/post-forwarder.yaml`.

### Endpoints

- `GET /` - Health check endpoint
- `POST /api/{token}/{service}` - Webhook endpoint for POST requests
- `GET /api/{token}/{service}` - Webhook endpoint for GET requests
- `GET /api/message/{id}` - Retrieve a stored message by ID

### Usage Examples

#### Send a webhook via POST:
```bash
curl -X POST http://your-server/api/your-auth-token/service-name --data "Your message content"
```

#### Send a webhook via GET:
```bash
curl "http://your-server/api/your-auth-token/service-name?param1=value1&param2=value2"
```

#### Retrieve a stored message:
```bash
curl http://your-server/api/message/message-id
```

## License

This project is licensed under the MIT License.