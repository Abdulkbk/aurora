# Bitcoin Core Dockerfile for Aurora
# Builds bitcoind from a specified git repository and branch/commit

FROM debian:bookworm-slim AS builder

# Build arguments for customizing the source
ARG checkout="master"
ARG git_url="https://github.com/bitcoin/bitcoin"

# Install build dependencies
RUN apt-get update && apt-get install -y \
    automake \
    autotools-dev \
    bsdmainutils \
    build-essential \
    ccache \
    clang \
    cmake \
    git \
    libtool \
    libboost-dev \
    libevent-dev \
    libsqlite3-dev \
    libzmq3-dev \
    pkg-config \
    python3 \
    && rm -rf /var/lib/apt/lists/*

# Clone and build Bitcoin Core
WORKDIR /build
RUN git clone $git_url bitcoin \
    && cd bitcoin \
    && git checkout $checkout \
    && cmake -B build \
        -DBUILD_TESTS=OFF \
        -DBUILD_BENCH=OFF \
        -DBUILD_WALLET_TOOL=ON \
        -DENABLE_WALLET=ON \
        -DENABLE_IPC=OFF \
        -DWITH_ZMQ=ON \
        -DWITH_SQLITE=ON \
        -DWITH_BDB=OFF \
    && cmake --build build -j$(nproc) \
    && cmake --install build --prefix /opt/bitcoin

# Final image
FROM debian:bookworm-slim AS final

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libevent-2.1-7 \
    libevent-extra-2.1-7 \
    libevent-pthreads-2.1-7 \
    libsqlite3-0 \
    libzmq5 \
    && rm -rf /var/lib/apt/lists/*

# Copy binaries from builder (core binaries only)
COPY --from=builder /opt/bitcoin/bin/bitcoind /usr/local/bin/
COPY --from=builder /opt/bitcoin/bin/bitcoin-cli /usr/local/bin/
COPY --from=builder /opt/bitcoin/bin/bitcoin-wallet /usr/local/bin/

# Create data directory
RUN mkdir -p /root/.bitcoin
VOLUME /root/.bitcoin

# Expose ports (mainnet p2p, mainnet rpc, testnet p2p, testnet rpc, regtest p2p, regtest rpc, zmq)
EXPOSE 8332 8333 18332 18333 18443 18444 28332 28333

# Default entrypoint
CMD ["bitcoind"]
