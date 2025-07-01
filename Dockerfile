# Use a Go version consistent with go.mod
FROM golang:1.17-alpine AS builder

# Install protoc and other build dependencies
RUN apk update && apk add --no-cache \
    protobuf-dev \
    protoc \
    mysql-client \
    build-base

WORKDIR /app

# Install Go gRPC plugins
# Ensure GOPATH/bin is in PATH for subsequent RUN commands if needed,
# or call them via their full path.
# For go install in Go 1.17+, binaries are placed in GOBIN or GOPATH/bin.
# Setting GOBIN to /usr/local/bin for simplicity.
ENV GOBIN /usr/local/bin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Generate gRPC code
# This assumes protoc is in PATH (installed by apk add protoc)
# and protoc-gen-go* are in GOBIN (which should be in PATH or use full path)
# The -I. sets the import path to the current directory (WORKDIR /app)
RUN protoc -I. \
    --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    internal/adapters/framework/left/grpc/proto/messages.proto \
    internal/adapters/framework/left/grpc/proto/arithmetic_service.proto

COPY ./grpc_entrypoint.sh /usr/local/bin/grpc_entrypoint.sh
RUN /bin/chmod +x /usr/local/bin/grpc_entrypoint.sh

RUN go build -o /app/main cmd/main.go

# --- Final Stage ---
# Use a smaller image for the final application
FROM alpine:latest

RUN apk --no-cache add mysql-client

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/grpc_entrypoint.sh /usr/local/bin/grpc_entrypoint.sh
RUN /bin/chmod +x /usr/local/bin/grpc_entrypoint.sh

EXPOSE 9000
CMD ["./main"]
RUN mv main /usr/local/bin/

CMD ["main"]
ENTRYPOINT ["grpc_entrypoint.sh"]
