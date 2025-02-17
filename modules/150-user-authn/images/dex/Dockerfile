# Based on https://github.com/dexidp/dex/blob/v2.31.0/Dockerfile
ARG BASE_GOLANG_16_ALPINE
ARG BASE_ALPINE
FROM $BASE_GOLANG_16_ALPINE as artifact
RUN apk add --no-cache git ca-certificates gcc build-base sqlite patch make curl
WORKDIR /dex
COPY patches/client-groups.patch patches/static-user-groups.patch /
RUN wget https://github.com/dexidp/dex/archive/v2.31.0.tar.gz -O - | tar -xz --strip-components=1 \
  && git apply /client-groups.patch \
  && git apply /static-user-groups.patch
RUN go build ./cmd/dex

FROM ghcr.io/dexidp/dex@sha256:104e0ba05220915de6e5d122895658e671bd8a252dbb550d4d10b8245f691ad4 as dex

FROM $BASE_ALPINE
RUN apk add --no-cache --update ca-certificates openssl
RUN mkdir -p /var/dex
RUN chown -R 1001:1001 /var/dex
RUN mkdir -p /etc/dex
RUN chown -R 1001:1001 /etc/dex
# Copy module files for CVE scanning / dependency analysis.
COPY --from=dex /etc/dex/config.docker.yaml /etc/dex/
COPY --from=dex /usr/local/src/dex/go.mod /usr/local/src/dex/go.sum /usr/local/src/dex/
COPY --from=dex /usr/local/src/dex/api/v2/go.mod /usr/local/src/dex/api/v2/go.sum /usr/local/src/dex/api/v2/

COPY --from=artifact /dex/dex /usr/local/bin/
COPY --from=artifact /dex/web /web

USER 1001:1001

CMD ["dex", "serve", "/etc/dex/config.docker.yaml"]
