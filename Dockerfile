FROM golang:1.7-onbuild

RUN apt-get update && \
apt-get install -y bash ca-certificates && \
rm -rf /var/lib/apt/lists/*

ENV GIN_MODE=debug
ENV PORT=8080

EXPOSE 8080

