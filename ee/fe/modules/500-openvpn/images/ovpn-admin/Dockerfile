ARG BASE_ALPINE
ARG BASE_NODE_16_ALPINE
ARG BASE_GOLANG_17_ALPINE

FROM $BASE_ALPINE as src
WORKDIR /app
COPY deckhouse.patch /
RUN apk update && apk add git patch && git clone -b fixes https://github.com/vitaliy-sn/ovpn-admin.git . && git checkout 05fcf99bdfe10f83b400fbdf41b17420c6651f07 && patch -p1 < /deckhouse.patch

FROM $BASE_NODE_16_ALPINE AS frontend
WORKDIR /app
COPY --from=src /app /app
RUN  cd frontend && npm install && npm run build

FROM $BASE_GOLANG_17_ALPINE AS backend
WORKDIR /app
COPY --from=src /app /app
RUN go build .

FROM $BASE_ALPINE
WORKDIR /app
RUN apk update && apk add bash openssl
COPY --from=backend /app/ovpn-admin /app/
COPY --from=frontend /app/frontend/static /app/frontend/static
COPY client.conf.tpl /app/templates/
COPY rootfs /
