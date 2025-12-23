FROM golang:1.25-alpine AS builder

WORKDIR /app

# deps
COPY go.mod go.sum ./
RUN go mod download

# source
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o router ./cmd

# ---- runtime ----
FROM alpine:3.20

WORKDIR /app

# ca-certs для http
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/router /app/router

EXPOSE 8080

ENTRYPOINT ["/app/router"]