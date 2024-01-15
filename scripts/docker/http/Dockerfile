FROM golang:1.21.6-alpine as builder

ENV ENVIRONMENT=dev
ENV SERVER_PORT=8080

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod verify

RUN go build -o /app/server ./cmd/http/main.go

EXPOSE $SERVER_PORT

CMD ["/app/server"]