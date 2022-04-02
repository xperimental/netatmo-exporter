FROM golang:1.17.8-alpine AS builder

RUN apk add --no-cache make git bash

WORKDIR /build

COPY go.mod go.sum /build/
RUN go mod download
RUN go mod verify

COPY . /build/
RUN make

FROM busybox
LABEL maintainer="Robert Jacob <xperimental@solidproject.de>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/netatmo-exporter /bin/netatmo-exporter

USER nobody
EXPOSE 9210

ENTRYPOINT ["/bin/netatmo-exporter"]
