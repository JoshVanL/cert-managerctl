FROM alpine:3.9

RUN apk --no-cache --update add ca-certificates

COPY cert-managerctl /usr/bin/

CMD ["/usr/bin/cert-managerctl"]
