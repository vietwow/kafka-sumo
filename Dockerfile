FROM golang:alpine

#Install dev tools
RUN apk add --update --no-cache alpine-sdk bash ca-certificates \
      libressl \
      tar \
      git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base coreutils librdkafka-dev

# Copy our source code into the container.
WORKDIR /go/src/kafka-sumo
ADD . /go/src/kafka-sumo

RUN ./build.sh

FROM alpine:latest
RUN apk --no-cache add ca-certificates librdkafka

WORKDIR /root/

COPY --from=0 /go/bin/kafka-sumo .

CMD ["./kafka-sumo", "--sumologic.url=https://endpoint3.collection.us2.sumologic.com/receiver/v1/http/ZaVnC4dhaV3roU3-kIZNtE2iTbUBN9crYJBdOVmhT0KD7_tmOjxO0dRuSn0DmGYcSOznK5nRyCqsaPm69tsW2CVr51cK9ook29GUHCJg8bd1oa1GpZzY7A==", "--sumologic.source.category=staging/logging"]