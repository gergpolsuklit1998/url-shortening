# Stage 1: Build the Go application

FROM golang:1.26-alpine AS build

WORKDIR /app

# Cache Go module downloads
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and build
COPY . .

# Build a statically linked binary
# CGO_ENABLED=0 removes C library dependencies
# -ldflags strips debug info for a smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server .

# Stage 2: Alpine production image
FROM alpine:3.23 AS runtime

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=build /app/server /server

USER appuser
EXPOSE 8080
ENTRYPOINT ["/server"]