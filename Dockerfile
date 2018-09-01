FROM golang:1 AS builder

RUN apt-get update && apt-get install -y upx

ENV PACKAGE=github.com/xperimental/netatmo-exporter

RUN mkdir -p /go/src/${PACKAGE}
WORKDIR /go/src/${PACKAGE}

ENV LD_FLAGS="-w"
ENV CGO_ENABLED=0

COPY . /go/src/${PACKAGE}
RUN echo "-- TEST" \
 && go test ./... \
 && echo "-- BUILD" \
 && go install -tags netgo -ldflags "${LD_FLAGS}" . \
 && echo "-- PACK" \
 && upx -9 /go/bin/netatmo-exporter

FROM busybox
LABEL maintainer="Robert Jacob <xperimental@solidproject.de>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/netatmo-exporter /bin/netatmo-exporter

USER nobody
EXPOSE 9210

ENTRYPOINT ["/bin/netatmo-exporter"]
