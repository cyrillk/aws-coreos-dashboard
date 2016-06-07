FROM golang:1.6-alpine

COPY . /go/src/github.com/cyrillk/aws-coreos-dashboard

RUN go build github.com/cyrillk/aws-coreos-dashboard
RUN go install github.com/cyrillk/aws-coreos-dashboard

EXPOSE 8080

CMD [""]
ENTRYPOINT ["aws-coreos-dashboard"]
