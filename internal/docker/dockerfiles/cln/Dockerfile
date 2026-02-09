# Core Lightning (CLN) Dockerfile for Aurora
# Builds lightningd from a specified git repository and branch/commit

FROM debian:bookworm-slim AS builder

# Build arguments for customizing the source
ARG checkout="master"
ARG git_url="https://github.com/ElementsProject/lightning"

# Install build dependencies
RUN apt-get update && apt-get install -y \
    autoconf \
    automake \
    build-essential \
    git \
    gettext \
    libgmp-dev \
    libsqlite3-dev \
    python3 \
    python3-pip \
    python3-venv \
    net-tools \
    zlib1g-dev \
    libsodium-dev \
    libtool \
    pkg-config \
    libpq-dev \
    cargo \
    rustc \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Create virtual environment and install Python dependencies
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN pip3 install --upgrade pip setuptools wheel \
    && pip3 install mako grpcio-tools

# Clone and build Core Lightning
WORKDIR /build
RUN git clone $git_url lightning \
    && cd lightning \
    && git checkout $checkout \
    && pip3 install -r requirements.txt \
    && ./configure --prefix=/opt/cln \
    && make -j$(nproc) \
    && make install

# Final image
FROM debian:bookworm-slim AS final

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libgmp10 \
    libsqlite3-0 \
    libsodium23 \
    libpq5 \
    python3 \
    python3-pip \
    python3-venv \
    && rm -rf /var/lib/apt/lists/*

# Copy virtual environment
COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Copy binaries from builder
COPY --from=builder /opt/cln/bin/lightningd /usr/local/bin/
COPY --from=builder /opt/cln/bin/lightning-cli /usr/local/bin/
COPY --from=builder /opt/cln/libexec /usr/local/libexec

# Create data directory
RUN mkdir -p /root/.lightning
VOLUME /root/.lightning

# Expose ports (p2p, rpc)
EXPOSE 9735 9835

# Default command
CMD ["lightningd"]
