docker stop tpm-client
docker rm tpm-client
docker build -f Dockerfile.client -t tpm-client .
docker run -d --name tpm-client --network tpm-network tpm-client
docker logs tpm-client