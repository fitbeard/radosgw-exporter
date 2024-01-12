FROM golang:1.21-alpine AS compile

WORKDIR /go/src

COPY . .

ENV CGO_ENABLED=0

RUN go build -o radosgw-exporter

FROM alpine:3.19

COPY --from=compile /go/src/radosgw-exporter /

ENTRYPOINT ["/radosgw-exporter"]
