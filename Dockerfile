# Use an official Golang image
FROM golang:1.23-alpine

# Install dependencies: OpenSSL, git, and any needed build tools
RUN apk add --no-cache openssl-dev gcc musl-dev curl unzip

# Set environment variables for CGO and OpenSSL
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/usr/include"
ENV CGO_LDFLAGS="-L/usr/lib"

# Set the working directory
WORKDIR /app

# Copy the Go project into the container
COPY . .

# Download and verify dependencies
RUN go mod tidy

# Install Protocol Buffer compiler
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip && \
    unzip protoc-21.12-linux-x86_64.zip -d /usr/local && \
    rm protoc-21.12-linux-x86_64.zip

# Expose the TPM simulator port (optional if you're running the simulator inside the container)
EXPOSE 2321

# Run the Go program by default
CMD ["go", "run", "main.go"]