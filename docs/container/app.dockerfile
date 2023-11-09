# Image Base of Go
FROM golang:1.21

# Working Directory
WORKDIR /app

# Copy
# - dependencies
COPY go.mod .
RUN go mod download

# Copy
# - source code
COPY ./cmd/* .
COPY ./internal/* .
COPY ./pkg/* .

# Build
RUN go build -o ./bin/app ./cmd/main.go