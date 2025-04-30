FROM golang:1.24-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o subscriber-service

RUN mkdir -p /app/uploads
RUN mkdir -p /app/compressed

# Create directories for compressed images

# Run the application
CMD ["./subscriber-service"]