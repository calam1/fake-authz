FROM golang:1.16-alpine as builder

# Adding the grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . /app

# Build the Go app
RUN go build -o main .

FROM alpine:3.13

COPY --from=builder /bin/grpc_health_probe ./grpc_health_probe

WORKDIR /app
# Copy our static executable.
COPY --from=builder /app .

#COPY --from=builder /bin/grpc_health_probe ./grpc_health_probe

# change user from root for security reasons
USER 9000

EXPOSE 50051

# Command to run the executable
ENTRYPOINT ["./main"]
