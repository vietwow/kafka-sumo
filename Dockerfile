FROM golang:alpine

#Install dev tools
RUN apk add --update --no-cache alpine-sdk bash ca-certificates \
      libressl \
      tar \
      git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base coreutils librdkafka-dev

# Copy our source code into the container.
WORKDIR /go/src/kafka-sumo
ADD . /go/src/kafka-sumo

RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates librdkafka

WORKDIR /root/

COPY --from=0 /go/bin/kafka-sumo .

CMD ["./kafka-sumo"]