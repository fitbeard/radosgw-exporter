FROM golang:1.21-alpine AS compile

WORKDIR /go/src

COPY . .

ENV CGO_ENABLED=0

RUN apk add upx

RUN go build -ldflags "-s -w" -o radosgw-exporter

RUN upx radosgw-exporter

FROM gcr.io/distroless/static-debian12:latest

COPY --from=compile /go/src/radosgw-exporter /

ENTRYPOINT ["/radosgw-exporter"]
