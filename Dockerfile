# ---------- build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# кеш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o auth ./cmd/server

# ---------- runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/auth /app/auth

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/auth"]
 