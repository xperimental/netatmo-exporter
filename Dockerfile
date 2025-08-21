FROM --platform=$BUILDPLATFORM golang:1.24.6-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN apk add --no-cache make git bash

WORKDIR /build

COPY go.mod go.sum /build/
RUN go mod download
RUN go mod verify

COPY . /build/
RUN make build-binary

FROM busybox
LABEL maintainer="Robert Jacob <xperimental@solidproject.de>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/netatmo-exporter /bin/netatmo-exporter

RUN mkdir -p /var/lib/netatmo-exporter/ \
 && chown nobody /var/lib/netatmo-exporter/

USER nobody
EXPOSE 9210

ENV NETATMO_EXPORTER_TOKEN_FILE=/var/lib/netatmo-exporter/netatmo-token.json
VOLUME /var/lib/netatmo-exporter/

ENTRYPOINT ["/bin/netatmo-exporter"]
