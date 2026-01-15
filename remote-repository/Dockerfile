FROM golang:1.25.5-alpine AS builder

RUN mkdir /app

WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8080 

CMD ["./main"]