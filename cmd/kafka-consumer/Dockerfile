############################
# STEP 1 build executable binary
############################
FROM golang:1.14 as builder

# kafka
ENV LIBRDKAFKA_VERSION 1.0.0

RUN apt-get -y update \
        && apt-get install -y --no-install-recommends upx-ucl zip libssl-dev \
        && apt-get clean \
        && rm -rf /var/lib/apt/lists/*


RUN curl -Lk -o /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz https://github.com/edenhill/librdkafka/archive/v${LIBRDKAFKA_VERSION}.tar.gz && \
      tar -xzf /root/librdkafka-${LIBRDKAFKA_VERSION}.tar.gz -C /root && \
      cd /root/librdkafka-${LIBRDKAFKA_VERSION} && \
      ./configure --prefix /usr && make && make install && make clean && ./configure --clean

# source
COPY ./go.mod /source/go.mod
COPY ./src /source/src
COPY ./cmd/svc /source/
COPY ./api /source/api
COPY ./vendor /source/vendor
COPY ./version /source/
COPY ./migrations /source/migrations
WORKDIR /source

RUN go env -w GOPROXY=off

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o /bin/svc ./main.go


############################
# STEP 2 build a small image
############################
FROM confluentinc/cp-kafka-connect

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /source/version /
COPY --from=builder /bin/svc /svc

# Run the svc binary.
CMD ["./svc"]
