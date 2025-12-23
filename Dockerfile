# use go based image to build the application
FROM golang:1.25.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# copy go mod and sum files
COPY go.mod go.sum ./

# download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# copy the source from the current directory to the Working Directory inside the container
COPY . .

# build the Go application
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/app/main.go

# start a new stage from scratch
FROM alpine:3.14

# create a non-root user to run the application
RUN adduser -D appuser

# copy application binary from builder stage and set ownership to non-root user
COPY --from=builder --chown=appuser:appuser /app/main /app/main
#COPY --from=builder --chown=appuser:appuser /app/.env /app/.env

# switch to non-root user
USER appuser

# exposed application port
EXPOSE 8080

# run the application binary
CMD ["/app/main"]