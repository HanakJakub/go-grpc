# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13.3-buster

# Add Maintainer Info
LABEL maintainer="Jakub Hanak <jakub.hanak2@gmail.com>"

# Make a directory for app
RUN mkdir /app

# Install dependencies 
RUN apt-get update && apt-get install -y \
    golang-goprotobuf-dev \
    protobuf-compiler

# add folder to app directory
ADD . /app

# Set the Current Working Directory inside the container
WORKDIR /app

# Run make to build up server and client
RUN make

# Start server
CMD ["./server"]