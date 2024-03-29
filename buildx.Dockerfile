# syntax=docker/dockerfile:1.4
FROM alpine:3.16

RUN apk --no-cache --no-progress add ca-certificates git \
    && rm -rf /var/cache/apk/*

COPY bibikoffi /

ENTRYPOINT ["/bibikoffi"]
EXPOSE 80
