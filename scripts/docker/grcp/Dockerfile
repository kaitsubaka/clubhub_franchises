FROM golang:1.21.6-alpine as builder

ENV ENVIRONMENT=dev
ENV SERVER_PORT=5001

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod verify

RUN go build -o /app/server ./cmd/grpc/main.go

EXPOSE $SERVER_PORT

CMD ["/app/server"]