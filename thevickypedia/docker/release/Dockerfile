# -----------------------------
# 1. Fetcher stage
# -----------------------------
FROM alpine AS fetcher

ARG VERSION
ENV VERSION=${VERSION:-latest}
ENV INSTALL_PATH="/"

RUN apk add --no-cache git
RUN git clone https://github.com/thevickypedia/filebrowser.git /src

COPY /src/thevickypedia/docker/release/release.sh /release.sh
RUN chmod +x /release.sh
RUN /release.sh

# -----------------------------
# 2. Runtime stage (minimal)
# -----------------------------
FROM alpine:latest

RUN apk add --no-cache ca-certificates mailcap curl jq git && \
    mkdir -p -m 777 /config /data

# Copy binary
COPY --from=backend-builder /filebrowser /filebrowser

# Healthcheck script
COPY --from=fetcher /src/docker/common/healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

EXPOSE 80

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD /healthcheck.sh || exit 1

COPY --from=fetcher /src/thevickypedia/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
