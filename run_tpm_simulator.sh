docker stop tpm-simulator
docker rm tpm-simulator
docker build -f Dockerfile.simulator -t tpm-simulator .
docker run -it --name tpm-simulator --network tpm-network tpm-simulator