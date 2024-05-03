FROM alpine:latest

# Set working directory
WORKDIR /opt/filebrowser

# Copy all files and folders in current directory to /opt/filebrowser
COPY . .
COPY thevickypedia_scripts/run.sh .

# todo: remove these steps and move them to volume mapping
# Copy config.json if present
RUN test -f thevickypedia_scripts/config.json && cp thevickypedia_scripts/config.json . || true

# Copy users.json if present
RUN test -f thevickypedia_scripts/users.json && cp thevickypedia_scripts/users.json . || true

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
ENV BASE_PATH /opt/filebrowser

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

# todo: change the run.sh, so that it can pick files from volume
VOLUME /srv
EXPOSE 80

# todo: make it so that the installation files get copied separately and installed separately
#RUN apk add --no-cache python3
#RUN python3 -c 'import os, shutil; files_to_keep = ["filebrowser", "config.json", "users.json", "run.sh"]; \
#                [os.remove(file) if os.path.isfile(file) else shutil.rmtree(file) for file in os.listdir() if file not in files_to_keep]'

RUN dos2unix run.sh
RUN chmod +x run.sh

RUN pwd
RUN ls -ltrha

ENTRYPOINT ["/opt/filebrowser/run.sh"]
