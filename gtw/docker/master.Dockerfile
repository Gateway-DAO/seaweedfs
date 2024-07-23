FROM golang:1.22-alpine AS build
LABEL author="Gateway"

RUN apk add --no-cache file libc6-compat

WORKDIR /opt/weed

# Build binaries
# - Copy
COPY go.mod go.sum /opt/weed/
RUN go mod download -x

COPY ./weed /opt/weed/weed

# - Build
WORKDIR /opt/weed/weed
RUN go build -tags edv -o /usr/bin/weed

WORKDIR /opt/weed/weed/mq/client/cmd/weed_pub_kv
RUN go build -o /usr/bin/weed_pub_kv

WORKDIR /opt/weed/weed/mq/client/cmd/weed_pub_record
RUN go build -o /usr/bin/weed_pub_record

WORKDIR /opt/weed/weed/mq/client/cmd/weed_sub_kv
RUN go build -o /usr/bin/weed_sub_kv

WORKDIR /opt/weed/weed/mq/client/cmd/weed_sub_record
RUN go build -o /usr/bin/weed_sub_record

FROM alpine AS final
LABEL author="Gateway"

RUN apk add --no-cache file libc6-compat

# Copy swfs binaries
COPY --from=build /usr/bin/weed /usr/bin/
COPY --from=build /usr/bin/weed_pub* /usr/bin/
COPY --from=build /usr/bin/weed_sub* /usr/bin/

# Copy swfs config
RUN mkdir -p /etc/seaweedfs
COPY gtw/docker/config/filer.toml /etc/seaweedfs/filer.toml

COPY gtw/docker/config/entrypoint.sh /entrypoint.sh
RUN apk add --no-cache fuse # for weed mount
RUN apk add --no-cache curl # for health checks

# volume server grpc port
EXPOSE 18080
# volume server http port
EXPOSE 8080
# filer server grpc port
EXPOSE 18888
# filer server http port
EXPOSE 8888
# master server shared grpc port
EXPOSE 19333
# master server shared http port
EXPOSE 9333
# s3 server http port
EXPOSE 8333
# webdav server http port
EXPOSE 7333

RUN mkdir -p /data/filerldb2

VOLUME /data
WORKDIR /data

RUN chmod +x /usr/bin/weed /usr/bin/weed_pub* /usr/bin/weed_sub*
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/usr/bin/weed"]
