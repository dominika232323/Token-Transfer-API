FROM golang:1.23.8

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o token-transfer ./src/main.go

EXPOSE 8080

CMD ["./token-transfer"]
