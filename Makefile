BINARY = weed

SOURCE_DIR = .
debug ?= 0
gco = 1

all: install

install:
	cd weed; go install

warp_install:
	go install github.com/minio/warp@v0.7.6

full_install:
	cd weed; go install -tags "elastic gocdk sqlite ydb tikv rclone"

server: install
	weed -v 0 server -s3 -filer -filer.maxMB=64 -volume.max=0 -master.volumeSizeLimitMB=1024 -volume.preStopSeconds=1 -s3.port=8000 -s3.allowEmptyFolder=false -s3.allowDeleteBucketNotEmpty=true -s3.config=./docker/compose/s3.json -metricsPort=9324

benchmark: install warp_install
	pkill weed || true
	pkill warp || true
	weed server -debug=$(debug) -s3 -filer -volume.max=0 -master.volumeSizeLimitMB=1024 -volume.preStopSeconds=1 -s3.port=8000 -s3.allowEmptyFolder=false -s3.allowDeleteBucketNotEmpty=false -s3.config=./docker/compose/s3.json &
	warp client &
	while ! nc -z localhost 8000 ; do sleep 1 ; done
	warp mixed --host=127.0.0.1:8000 --access-key=some_access_key1 --secret-key=some_secret_key1 --autoterm
	pkill warp
	pkill weed

# curl -o profile "http://127.0.0.1:6060/debug/pprof/profile?debug=1"
benchmark_with_pprof: debug = 1
benchmark_with_pprof: benchmark

.PHONY: images
images:
	make -C gtw/docker images

COMPOSE_FILES := -f ./gtw/docker/docker-compose.local.yml
COMPOSE_CMD := VOLUME_LOG_LEVEL=4 docker compose $(COMPOSE_FILES)

.PHONY: network stop logs
network: images
	VOLUME_LOG_LEVEL=4 $(COMPOSE_CMD) up -d --build
restart:
	$(COMPOSE_CMD) restart
stop:
	$(COMPOSE_CMD) down -v
logs:
	$(COMPOSE_CMD) logs -f

# Log events on locally mounted event store
.PHONY: events
events:
	leveldbutil dump ./data/volume1/events/*.ldb; leveldbutil dump ./data/volume1/events/*.log

.PHONY: ec2-binaries
ec2-binaries:
	@test -d ./bin || mkdir ./bin
	cd ./weed && CGO_ENABLED=$(cgo) GOOS=linux GOARCH=amd64 go build $(options) && mv weed ../bin/
	cd ./weed/mq/client/cmd/weed_pub_kv && CGO_ENABLED=$(cgo) GOOS=linux GOARCH=amd64 go build && mv weed_pub_kv ../../../../../bin/
	cd ./weed/mq/client/cmd/weed_pub_record && CGO_ENABLED=$(cgo) GOOS=linux GOARCH=amd64 go build && mv weed_pub_record ../../../../../bin/
	cd ./weed/mq/client/cmd/weed_sub_kv && CGO_ENABLED=$(cgo) GOOS=linux GOARCH=amd64 go build && mv weed_sub_kv ../../../../../bin/
	cd ./weed/mq/client/cmd/weed_sub_record && CGO_ENABLED=$(cgo) GOOS=linux GOARCH=amd64 go build && mv weed_sub_record ../../../../../bin/

.PHONY: dev
dev: ec2-binaries network

test:
	cd weed; go test -tags "elastic gocdk sqlite ydb tikv rclone" -v ./...
