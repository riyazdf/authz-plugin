FROM golang:1.7.1-alpine

RUN apk add --update make git gcc libc-dev && rm -rf /var/cache/apk/*

ENV PLUGINDIR /go/src/github.com/riyazdf/authz-plugin

COPY . ${PLUGINDIR}

RUN chmod -R a+rw /go

RUN cd ${PLUGINDIR} && make binary && mv bin/authz-plugin /go/bin/authz-plugin
RUN rm -rf ${PLUGINDIR}

ENTRYPOINT ["bin/authz-plugin"]