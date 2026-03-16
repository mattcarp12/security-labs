FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod main.go ./
RUN go build -o main .
CMD ["./main"]