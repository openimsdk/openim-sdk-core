# Stage 1: Build environment
FROM golang:latest AS build-env

WORKDIR /app

# Copy the necessary files for building
COPY . .

# Define build arguments
ARG OS
ARG ARCH
ARG BIN_DIR=/_output/bin

# Build the binary for the specified platform
RUN GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o ${BIN_DIR}/openIM.wasm wasm/cmd/main.go

# Stage 2: Ubuntu-based final image
FROM ubuntu:latest

# Copy the built binaries from the previous stage
COPY --from=build-env /output /app

# Set the working directory
WORKDIR /app

# Set up any necessary dependencies or configurations here

# Define any runtime instructions
CMD ["./_output/bin/openIM.wasm"]  # Update this with the actual name of your executable
