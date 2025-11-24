FROM ubuntu:22.04

ARG VERSION
ENV VERSION=${VERSION:-latest}
ENV INSTALL_PATH="/"

RUN apt-get update && apt-get install -y ca-certificates curl jq bash tar gzip && \
    mkdir -p -m 777 /config /data

COPY thevickypedia/docker/release/release.sh /release.sh
RUN chmod +x /release.sh
RUN /release.sh

COPY thevickypedia/docker/healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

COPY thevickypedia/docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN ls -l /
EXPOSE 80

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD /healthcheck.sh || exit 1

ENTRYPOINT ["/entrypoint.sh"]
