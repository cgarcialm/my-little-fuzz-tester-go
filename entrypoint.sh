#!/bin/bash
set -e

# Start swtpm and tpm2-abrmd in the background
swtpm socket --tpmstate dir=/tmp/mytpm1 \
    --ctrl type=tcp,port=2322 \
    --tpm2 \
    --server type=tcp,port=2321 &

tpm2-abrmd --tcti=tabrmd:bus_name=org.tpm2-abrmd &

# Wait for TPM services to initialize
sleep 2

# Execute the passed command
exec "$@"
