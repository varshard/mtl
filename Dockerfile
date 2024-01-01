FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/myapp ./cmd/mtl

FROM alpine:latest

COPY --from=builder /app/myapp /app/myapp

EXPOSE 3000

ENTRYPOINT ["/app/myapp"]
