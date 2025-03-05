FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o agent ./cmd/agent

FROM debian:latest
WORKDIR /root/
COPY --from=builder /app/agent .

EXPOSE 8081

CMD ["./agent"]
