############################
# STEP 1 build executable binary
############################
FROM golang:1.15 as builder

# kafka
ENV LIBRDKAFKA_VERSION 1.5.3

RUN apt-get -y update \
        && apt-get install -y --no-install-recommends upx-ucl zip libssl-dev \
        && apt-get clean \
        && rm -rf /var/lib/apt/lists/*


RUN curl -Lk -o /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz https://github.com/edenhill/librdkafka/archive/v${LIBRDKAFKA_VERSION}.tar.gz && \
      tar -xzf /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz -C /root && \
      cd /root/librdkafka-${LIBRDKAFKA_VERSION} && \
      ./configure --prefix /usr && make && make install && make clean && ./configure --clean

# source
COPY . /source/
WORKDIR /source

RUN go env -w GOPROXY=off

RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -tags dynamic -a -o /bin/svc ./cmd/kafka-producer/main.go

CMD ["/bin/svc"]
