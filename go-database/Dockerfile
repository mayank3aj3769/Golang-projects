# Start from the official Go image to create a build artifact.
# This image contains all the necessary tools to compile your application.
FROM golang:1.21.5 AS builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies.
# Doing this before copying the code can help leverage Docker cache
# and speed up subsequent builds if dependencies haven't changed.
COPY go.mod go.sum ./

# Download all the dependencies specified in go.mod and go.sum.
RUN go mod download

# Copy the rest of the project's source code.
COPY . .

# Compile the application to a binary.
# CGO_ENABLED=0 is set to ensure a static build, important for a scratch or an alpine image.
# GOOS=linux and GOARCH=amd64 ensure the binary is compatible with the target system.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Start from a fresh, minimal image to create a lean and secure deployment image.
FROM alpine:latest

# Add CA certificates
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container.
WORKDIR /root/

# Copy the binary from the builder stage to the production image.
COPY --from=builder /app/main .

# Expose port 8080 to the outside world.
EXPOSE 8000

# Command to run the executable.
CMD ["./main"]
