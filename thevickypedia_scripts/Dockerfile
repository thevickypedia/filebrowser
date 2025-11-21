# -----------------------------
# 1. Frontend build stage
# -----------------------------
FROM node:24-alpine AS frontend-builder

# Install required build tools
RUN apk add --no-cache git curl jq mailcap ca-certificates

RUN git clone https://github.com/thevickypedia/filebrowser.git

WORKDIR /app

# Copy frontend source
COPY filebrowser/frontend ./frontend

# Install deps & build frontend
RUN npm install -g pnpm && \
    cd frontend && \
    pnpm install --frozen-lockfile && \
    pnpm run build


# -----------------------------
# 2. Backend build stage
# -----------------------------
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# Install required build tools
RUN apk add --no-cache git curl jq mailcap ca-certificates

RUN git clone https://github.com/thevickypedia/filebrowser.git

# Copy full project
COPY filebrowser .

# Copy built frontend from stage 1
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Determine version information
ARG VERSION
ARG VERSION_HASH

# Fallback if no git metadata exists
RUN if [ -z "$VERSION" ]; then \
        VERSION=$(git describe --tags --always --match=v* 2>/dev/null || cat .version 2>/dev/null || echo v0); \
    fi && \
    if [ -z "$VERSION_HASH" ]; then \
        VERSION_HASH=$(git rev-parse HEAD 2>/dev/null || echo dev); \
    fi && \
    echo "Version: $VERSION / Hash: $VERSION_HASH" && \
    export LDFLAGS="-X '$(go list -m)/version.Version=$VERSION' -X '$(go list -m)/version.CommitSHA=$VERSION_HASH'" && \
    go build -ldflags "$LDFLAGS" -o /filebrowser .

# -----------------------------
# 3. Runtime stage (minimal)
# -----------------------------
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates mailcap curl jq git && \
    mkdir -p -m 777 /config /data

# Copy binary
COPY --from=backend-builder /filebrowser /filebrowser

# Copy config
COPY docker_config.json /filebrowser.json

# Healthcheck script
COPY healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

EXPOSE 80

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD /healthcheck.sh || exit 1

ENTRYPOINT ["/filebrowser"]
CMD ["--root=/data", "--address=0.0.0.0", "--port=80", "--database=/config/filebrowser.db"]
