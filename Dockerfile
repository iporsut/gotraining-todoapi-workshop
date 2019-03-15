FROM golang:1.12 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o todo-api

FROM ubuntu
COPY --from=builder /app/todo-api /
CMD ["/todo-api"]
