# Gin + Mailexam

Minimal [Gin](https://gin-gonic.com/) example that sends test mail through [Mailexam](https://mailexam.ru/) SMTP via the Go standard library [`net/smtp`](https://pkg.go.dev/net/smtp).

Based on the [Mailexam Gin guide](https://wiki.mailexam.ru/en/examples/gin/).

## What you need

- A Mailexam account and a project with SMTP credentials.
- Go 1.22+.

From your Mailexam welcome email or dashboard:

| Variable | Description |
|----------|-------------|
| `MAILEXAM_LOGIN` | SMTP login (for example, `xxxxx`) |
| `MAILEXAM_PASSWORD` | SMTP password (paired with the login) |
| Host | `{MAILEXAM_LOGIN}.mailexam.ru` (built automatically in code) |

## Quick start (host)

1. Install dependencies:

```bash
go mod tidy
```

2. Copy the example environment file and fill in your credentials:

```bash
cp .env.example .env
```

3. Edit `.env`:

```env
MAILEXAM_LOGIN=YOUR_LOGIN
MAILEXAM_PASSWORD=YOUR_PASSWORD
MAILEXAM_PORT=587
MAIL_FROM=noreply@example.test
```

4. Run the server:

```bash
go run .
```

The server listens on `http://127.0.0.1:8080` by default.

5. Send a test message:

```bash
curl -X POST http://127.0.0.1:8080/mail/test \
  -H 'Content-Type: application/json' \
  -d '{"to":"user@example.test","subject":"Test","body":"Hello"}'
```

The message appears in the Mailexam dashboard → your project → inbox.

## Environment variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `MAILEXAM_LOGIN` | yes | — | SMTP login; also used to build the host name |
| `MAILEXAM_PASSWORD` | yes | — | SMTP password |
| `MAILEXAM_PORT` | no | `587` | SMTP port (`587`, `2525`, or `25`) |
| `MAIL_FROM` | no | `noreply@example.test` | Sender address (any test address is fine) |
| `BIND_ADDR` / `HTTP_HOST` | no | `127.0.0.1` | HTTP bind address |
| `PORT` / `HTTP_PORT` | no | `8080` | HTTP listen port |

For port **587** sending uses STARTTLS (`sendWithSTARTTLS`). For port **25** it uses plain `smtp.SendMail`.

## Project layout

```
.
├── go.mod
├── mail.go             # Mailexam SMTP config and sendTest()
├── main.go             # HTTP server and POST /mail/test
├── .env.example
├── Dockerfile          # for local debugging only
└── docker-compose.yml
```

## Docker (debugging)

Docker is provided for local debugging. For day-to-day development, run the app on the host with `go run .` (see above).

```bash
cp .env.example .env
# edit .env with your credentials

docker compose up --build
```

Then call the same endpoint on the mapped port:

```bash
curl -X POST http://127.0.0.1:8080/mail/test \
  -H 'Content-Type: application/json' \
  -d '{"to":"user@example.test","subject":"Test","body":"Hello"}'
```

Inside the container the server binds to `0.0.0.0:8080` so the port mapping works.

## CI

Set these secrets in your CI environment:

```yaml
variables:
  MAILEXAM_LOGIN: $MAILEXAM_LOGIN
  MAILEXAM_PASSWORD: $MAILEXAM_PASSWORD
  MAILEXAM_PORT: "587"
  MAIL_FROM: "noreply@example.test"
```

After sending a message in a test, verify delivery via the [Mailexam API](https://mailexam.ru/api).

## Troubleshooting

**TLS or authentication failed**

- Host must be `{login}.mailexam.ru`, where `{login}` matches `MAILEXAM_LOGIN`.
- Login and password must come from the same Mailexam project.

**Port 587**

- Requires STARTTLS (`sendWithSTARTTLS`), not direct `smtp.SendMail` without TLS.

**Message not in the dashboard**

- Open the inbox of the same Mailexam project.
- Run with `GIN_MODE=debug` and check the error text in the `500` response.

**Port already in use**

- Change `PORT` or `HTTP_PORT` in `.env`.

## See also

- [Mailexam Gin guide (wiki)](https://wiki.mailexam.ru/en/examples/gin/)
- [Axum](https://github.com/mailexam/Axum), [Actix Web](https://github.com/mailexam/Actix) — other high-performance frameworks
- [Gin documentation](https://gin-gonic.com/docs/)
- [Mailexam API documentation](https://mailexam.ru/api)
