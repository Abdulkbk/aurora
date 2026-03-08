# btcd Dockerfile for Aurora
# Builds btcd from a specified git repository and branch/commit

FROM golang:1.24-alpine AS builder

# Build arguments for customizing the source
ARG checkout="master"
ARG git_url="https://github.com/btcsuite/btcd"

# Install build dependencies
RUN apk add --no-cache --update alpine-sdk git make gcc

# Clone and build btcd
WORKDIR /go/src/github.com/btcsuite/btcd
RUN git clone $git_url . \
    && git checkout $checkout \
    && go install -v . ./cmd/...

# Final image
FROM alpine:latest AS final

# Install runtime dependencies
RUN apk --no-cache add bash curl

# Copy binaries from builder
COPY --from=builder /go/bin/btcd /usr/local/bin/
COPY --from=builder /go/bin/btcctl /usr/local/bin/
COPY --from=builder /go/bin/gencerts /usr/local/bin/
COPY --from=builder /go/bin/addblock /usr/local/bin/
COPY --from=builder /go/bin/findcheckpoint /usr/local/bin/

# Create data directory
RUN mkdir -p /home/btcd/.btcd
VOLUME /home/btcd/.btcd

# Expose ports (rpc, simnet rpc)
EXPOSE 18334 18444

# Default command
CMD ["btcd"]
