FROM golang:1.6-alpine

RUN apk add --no-cache --upgrade bash ca-certificates

COPY . /go/src/github.com/cyrillk/aws-coreos-dashboard
WORKDIR /go/src/github.com/cyrillk/aws-coreos-dashboard
RUN mv docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN go build github.com/cyrillk/aws-coreos-dashboard
RUN go install github.com/cyrillk/aws-coreos-dashboard

EXPOSE 8080

ENTRYPOINT ["docker-entrypoint.sh"]
CMD [""]
