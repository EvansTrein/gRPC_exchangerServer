FROM golang:1.23.3-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY migrations ./
RUN go mod download

# Need to work with SQLite
RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1

COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app

# Need to work with SQLite
RUN apk add --no-cache sqlite libc6-compat

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/internal/storages/exchanger.db ./internal/storages/exchanger.db

EXPOSE 44000
CMD ["./main", "-config", "./config.yaml"]