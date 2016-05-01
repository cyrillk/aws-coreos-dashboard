FROM golang:alpine

RUN apk update && apk upgrade && \
apk add wget && \
rm -rfv /var/cache/apk/* /tmp/* /var/tmp/*

ENV FLEET_VERSION 0.11.5

# Install fleetctl static binary
RUN \
  wget -P /tmp https://github.com/coreos/fleet/releases/download/v${FLEET_VERSION}/fleet-v${FLEET_VERSION}-linux-amd64.tar.gz && \
  gunzip /tmp/fleet-v${FLEET_VERSION}-linux-amd64.tar.gz && \
  tar -xf /tmp/fleet-v${FLEET_VERSION}-linux-amd64.tar -C /tmp && \
  mv /tmp/fleet-v${FLEET_VERSION}-linux-amd64/fleetctl /bin/ && \
  rm -rf /tmp/fleet-v${FLEET_VERSION}-linux-amd64*

COPY . /go/src/github.com/cyrillk/aws-coreos-dashboard

RUN go install github.com/cyrillk/aws-coreos-dashboard

EXPOSE 8080

CMD [""]
ENTRYPOINT ["/go/bin/aws-coreos-dashboard"]
