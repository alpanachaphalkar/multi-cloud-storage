FROM golang:1.15.0-alpine

ARG storageAccountName=abc
ENV storageAccountName=$storageAccountName
ARG accessKey=xxxxx
ENV accessKey=$accessKey
ARG containerName=container
ENV containerName=$containerName
ARG bucket_name=bkt
ENV bucket_name=$bucket_name
ARG PrivateKeyData=xxxxxxxxxx
ENV PrivateKeyData=$PrivateKeyData

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/xPlorinRolyPoly/multi-cloud-storage

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["multi-cloud-storage"]
