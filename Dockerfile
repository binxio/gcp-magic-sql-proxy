FROM golang:1.14
WORKDIR /go

ADD . /go/src/github.com/binxio/gcp-magic-sql-proxy

RUN go get github.com/binxio/gcp-magic-sql-proxy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' github.com/binxio/gcp-magic-sql-proxy

FROM gcr.io/cloudsql-docker/gce-proxy:1.16
COPY --from=0		/go/gcp-magic-sql-proxy /

ENTRYPOINT [ "/gcp-magic-sql-proxy" ]
