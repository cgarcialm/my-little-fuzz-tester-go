# Stage 1: Builder
FROM golang:1.23.2 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Install necessary packages
RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    pkg-config \
    libssl-dev \
    swtpm \
    tpm2-tools \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./

# Download dependencies (including go-tpm)
RUN go mod download

# Copy source code (for the fuzzer and any additional code)
COPY src/ ./src/

# Build the fuzzer application
RUN go build -o tpm_fuzz ./src/main.go

# Stage 2: Runtime
FROM golang:1.20

# Install necessary packages
RUN apt-get update && apt-get install -y \
    swtpm \
    tpm2-tools \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the built fuzzing binary from the builder
COPY --from=builder /app/tpm_fuzz ./

# Create directory for TPM state
RUN mkdir -p /tmp/mytpm1

# Start swtpm and tpm2-abrmd in the background
RUN swtpm socket --tpmstate dir=/tmp/mytpm1 \
    --ctrl type=tcp,port=2322 \
    --tpm2 \
    --server type=tcp,port=2321 & \
    tpm2-abrmd --tcti=tabrmd:bus_name=org.tpm2-abrmd &

# Wait for TPM services to initialize
CMD ["bash", "-c", "sleep 2 && ./tpm_fuzz"]
