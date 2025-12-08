# backend/Dockerfile
FROM golang:1.24-alpine AS build
WORKDIR /app

# Cache and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the cmd/server binary (refactor: main moved to cmd/server)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# --- Final minimal runtime image ---
FROM alpine:3.18
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /server /app/server
COPY --from=build /app/migrations /app/migrations
USER app
ENV MIGRATIONS_PATH=/app/migrations
EXPOSE 8080
CMD ["/app/server"]
