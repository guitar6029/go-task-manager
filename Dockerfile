FROM golang:1.26-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
#ENV CGO_ENABLED=1

# build api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# later cli
# RUN go build -o cli ./cmd/cli

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y curl

COPY --from=builder /app/api .

EXPOSE 8080
CMD ["./api"]