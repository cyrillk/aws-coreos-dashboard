FROM golang

RUN apk update && apk upgrade && \
apk add --update bash && \
rm -rfv /var/cache/apk/* /tmp/* /var/tmp/*

COPY aws-coreos-dashboard /opt/aws-coreos-dashboard
COPY docker-entrypoint.sh /opt/docker-entrypoint.sh

EXPOSE 8080

WORKDIR /opt

ENTRYPOINT ["/opt/docker-entrypoint.sh"]
CMD [""]
