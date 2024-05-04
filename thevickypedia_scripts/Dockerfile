FROM alpine:latest

# Set working directory
WORKDIR /opt/filebrowser

# Copy all files and folders in current directory to /opt/filebrowser
COPY . .

RUN apk --update add ca-certificates \
                     mailcap \
                     curl \
                     jq \
                     git

ENV GO_VERSION=1.22.2

# Install Node.js, npm, and Yarn
RUN apk add --no-cache nodejs npm yarn

# Download and install Go
RUN wget -O go.tar.gz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV GOBIN=/go/bin

# Set base path
ENV BASE_PATH $WORKDIR

# Determine version
RUN VERSION=$(git describe --tags --always --match=v* 2> /dev/null || cat "$BASE_PATH/.version" 2> /dev/null || echo v0) && \
    export VERSION

# Determine version hash
RUN VERSION_HASH=$(git rev-parse HEAD) && \
    export VERSION_HASH

# Disable GOGC
ENV GOGC=off

# Determine module
RUN MODULE=$(env GO111MODULE=on go list -m) && \
    export MODULE

# Set tools directory
ENV TOOLS_DIR $BASE_PATH/tools

# Set tools binary directory
ENV TOOLS_BIN $TOOLS_DIR/bin

# Create tools binary directory
RUN mkdir -p "$TOOLS_BIN"

# Add tools binary directory to PATH
ENV PATH="$TOOLS_BIN:$PATH"

# Set LDFLAGS
ENV LDFLAGS="-X \"$MODULE/version.Version=$VERSION\" -X \"$MODULE/version.CommitSHA=$VERSION_HASH\""

# IMPORTANT: frontend should be built first, since it is required when creating the binary
# Build frontend
RUN cd frontend && \
    npm install && \
    npm ci && \
    npm run build && \
    cd ..

# Build backend
RUN go build -ldflags "$LDFLAGS" -o .

# Expose port 80 to host on http within the container
EXPOSE 80

# COPY runs at start, so use RUN instead
RUN cp docker_config.json /filebrowser.json
RUN cp filebrowser /filebrowser

RUN rm -rf /opt

WORKDIR /

COPY healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh  # Make the script executable

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD /healthcheck.sh || exit 1

RUN mkdir -p -m 777 /config \
    && mkdir -p -m 777 /data

ENTRYPOINT ["/filebrowser"]
CMD ["--root=/data", "--address=0.0.0.0", "--port=80", "--database=/config/filebrowser.db"]
