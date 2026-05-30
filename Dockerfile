# Development and CI image — not documented in README (host setup is the reference path).
FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o gin-mailexam .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/gin-mailexam .

ENV BIND_ADDR=0.0.0.0
ENV PORT=8080

EXPOSE 8080

CMD ["./gin-mailexam"]
