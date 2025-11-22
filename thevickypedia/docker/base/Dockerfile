# -----------------------------
# 0. Fetcher stage
# -----------------------------
FROM alpine AS fetcher

RUN apk add --no-cache git
RUN git clone https://github.com/thevickypedia/filebrowser.git /src

# -----------------------------
# 1. Frontend build stage
# -----------------------------
FROM node:24-alpine AS frontend-builder

RUN apk add --no-cache curl jq mailcap ca-certificates

# Enable build scripts (pnpm 9+)
ENV PNPM_ENABLE_PREPOST_INSTALL_SCRIPTS=true

# Bring frontend source from fetcher stage
COPY --from=fetcher /src/frontend /app/frontend

WORKDIR /app/frontend

# Install deps & build
RUN npm install -g pnpm && \
    pnpm install --frozen-lockfile && \
    pnpm run build

# -----------------------------
# 2. Backend build stage
# -----------------------------
FROM golang:1.25-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev git curl jq mailcap ca-certificates

WORKDIR /app

# Bring entire repo
COPY --from=fetcher /src /app

# Bring built frontend
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Determine version information
ARG VERSION
ARG VERSION_HASH

RUN if [ -z "$VERSION" ]; then \
        VERSION=$(git describe --tags --always --match=v* 2>/dev/null \
                  || cat .version 2>/dev/null \
                  || echo v0); \
    fi && \
    if [ -z "$VERSION_HASH" ]; then \
        VERSION_HASH=$(git rev-parse HEAD 2>/dev/null || echo dev); \
    fi && \
    echo "Version: $VERSION / Hash: $VERSION_HASH" && \
    export LDFLAGS="-X '$(go list -m)/version.Version=$VERSION' \
                    -X '$(go list -m)/version.CommitSHA=$VERSION_HASH'" && \
    CGO_ENABLED=1 go build -ldflags "$LDFLAGS" -o /filebrowser .

# -----------------------------
# 3. Runtime stage (minimal)
# -----------------------------
FROM alpine:latest

ARG VERSION
ENV TAG=$VERSION

RUN apk add --no-cache ca-certificates mailcap curl jq git bash tar gzip && \
    mkdir -p -m 777 /config /data

# Copy binary
COPY --from=backend-builder /filebrowser /filebrowser

# Healthcheck script
COPY --from=fetcher /src/docker/common/healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

EXPOSE 80

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD /healthcheck.sh || exit 1

COPY --from=fetcher /src/thevickypedia/docker/base/upload.sh /upload.sh
RUN chmod +x /upload.sh
# Read secrets from BuildKit into environment variables DURING BUILD ONLY
RUN --mount=type=secret,id=git_token \
    --mount=type=secret,id=repo_name \
    export GIT_TOKEN="$(cat /run/secrets/git_token)" && \
    export GITHUB_REPOSITORY="$(cat /run/secrets/repo_name)" && \
    bash /upload.sh _alpine

COPY --from=fetcher /src/thevickypedia/docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
