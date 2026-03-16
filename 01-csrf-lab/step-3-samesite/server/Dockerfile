# server/Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Run the executable
CMD ["./main"]