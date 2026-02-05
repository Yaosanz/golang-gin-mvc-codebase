# Root Dockerfile to build backend when Railway ignores rootDirectory
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# copy go mod and sum files (backend context)
COPY backend/go.mod backend/go.sum ./

# download dependencies
RUN go mod download

# copy backend source
COPY backend/ .

# build application binaries
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
