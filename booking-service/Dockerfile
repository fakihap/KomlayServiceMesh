FROM golang:1.23.3-alpine

WORKDIR /app

COPY . .

RUN go build -o main main.go

CMD ["./main"]
