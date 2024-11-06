# Use an official Golang image
FROM golang:1.23-alpine

# Install dependencies: OpenSSL, git, build tools, and protoc
RUN apk add --no-cache openssl-dev gcc musl-dev curl unzip protobuf

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

# Add Go bin to PATH to make the protoc plugins accessible
ENV PATH="/go/bin:${PATH}"

# Install Go Protobuf plugins: protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code from Protobuf schema
RUN protoc --go_out=MTest-Suite/testSchema.proto=github.com/cgarcialm/my-little-fuzz-tester-go/Test-Suite:. --go_opt=paths=source_relative \
           --go-grpc_out=MTest-Suite/testSchema.proto=github.com/cgarcialm/my-little-fuzz-tester-go/Test-Suite:. --go-grpc_opt=paths=source_relative \
           Test-Suite/testSchema.proto

# Expose the TPM simulator port (optional if you're running the simulator inside the container)
EXPOSE 2321

# Run the Go program by default
CMD ["go", "run", "main.go"]
