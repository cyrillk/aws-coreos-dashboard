FROM golang:alpine

RUN apk update && apk upgrade && \
rm -rfv /var/cache/apk/* /tmp/* /var/tmp/*

COPY . /go/src/github.com/cyrillk/aws-coreos-dashboard

RUN go install github.com/cyrillk/aws-coreos-dashboard

EXPOSE 8080

CMD [""]
ENTRYPOINT ["/go/bin/aws-coreos-dashboard"]
