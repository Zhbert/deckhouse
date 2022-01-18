ARG BASE_ALPINE
ARG BASE_NODE_16_ALPINE
ARG BASE_GOLANG_17_ALPINE

FROM $BASE_ALPINE as src
WORKDIR /app
RUN apk update && apk add git patch openvpn && git clone -b fixes https://github.com/vitaliy-sn/ovpn-admin.git . && git checkout 40cce34bd1c6e232c0b3b95263e36bdc5203b072

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
