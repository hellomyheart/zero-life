# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata musl

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .

RUN mkdir -p /app/data
VOLUME ["/app/data"]

EXPOSE 8080

CMD ["./server"]
