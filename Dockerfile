FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./cmd/main.exe ./cmd/main.go

EXPOSE 8080

CMD ["./cmd/main.exe"]
