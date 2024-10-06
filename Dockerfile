# Use the official Golang image to create a build artifact.
FROM golang:1.22.6 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Install the dependencies
RUN go mod download

# Copy the source code and .env file into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM gcr.io/distroless/base

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the config.yaml file from the builder stage to the final image
COPY --from=builder /app/config /config

# Command to run the executable
CMD ["/main"]