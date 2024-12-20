FROM golang:1.23.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o grpcserver ./cmd/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/grpcserver .
EXPOSE 8080

CMD ["./grpcserver"]