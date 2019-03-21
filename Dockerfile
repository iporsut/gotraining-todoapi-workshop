FROM golang:1.12 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOARCH=amd64 GOOS=linux go build -o todo-api

FROM ubuntu
COPY --from=builder /app/todo-api /
CMD ["/todo-api"]
