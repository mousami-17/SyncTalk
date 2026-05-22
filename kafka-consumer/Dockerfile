FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app

# Build the application.
RUN go build -o main .

# Run the executable.
CMD ["./main"]