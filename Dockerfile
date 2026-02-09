# ---------- build stage ----------
FROM golang:1.24.7-alpine AS builder

WORKDIR /app

# кеш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -trimpath -ldflags="-s -w" -o titan-auth ./cmd/auth-service && chmod +x titan-auth

# ---------- runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/titan-auth /app/titan-auth

EXPOSE 1489

USER nonroot:nonroot

ENTRYPOINT ["/app/titan-auth"]
