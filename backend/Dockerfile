FROM golang:1.25.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kidtask ./cmd/server/main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/kidtask .
COPY --from=builder /app/migrations ./migrations
RUN chmod +x kidtask
CMD ["./kidtask"]