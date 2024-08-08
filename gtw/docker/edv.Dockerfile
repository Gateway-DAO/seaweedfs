# ############## BUILD ################ #

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

# ############## FINAL ################ #

FROM alpine AS final
LABEL author="Gateway"

RUN apk add --no-cache file libc6-compat

# Copy swfs binaries
COPY --from=build /usr/bin/weed /usr/bin/

# Copy swfs config
RUN mkdir -p /etc/seaweedfs
COPY gtw/docker/config/filer.toml /etc/seaweedfs/filer.toml
COPY gtw/docker/config/kafka-edv.toml /etc/seaweedfs/kafka.toml

COPY gtw/docker/config/entrypoint.sh /entrypoint.sh
RUN apk add --no-cache fuse # for weed mount
RUN apk add --no-cache curl # for health checks

# volume server grpc port
EXPOSE 18080
# volume server http port
EXPOSE 8080

RUN mkdir -p /data/filerldb2

VOLUME /data
WORKDIR /data

RUN chmod +x /usr/bin/weed
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/usr/bin/weed"]
