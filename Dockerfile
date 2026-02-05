# Build from repo root with backend context
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Copy backend go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migration cmd/migration/main.go

FROM alpine:3.14

RUN adduser -D appuser

WORKDIR /app
COPY --from=builder --chown=appuser:appuser /app/main /app/main
COPY --from=builder --chown=appuser:appuser /app/migration /app/migration
COPY --from=builder --chown=appuser:appuser /app/database/migrations /app/database/migrations
COPY --from=builder --chown=appuser:appuser /app/scripts/railway-start.sh /app/railway-start.sh

RUN chmod +x /app/railway-start.sh

USER appuser

EXPOSE 8080

CMD ["/app/railway-start.sh"]
