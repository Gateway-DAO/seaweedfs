all: gen

.PHONY : gen

gen:
	protoc -I ./ master.proto --go_out=./master_pb --go-grpc_out=./master_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ volume_server.proto --go_out=./volume_server_pb --go-grpc_out=./volume_server_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ filer.proto --go_out=./filer_pb --go-grpc_out=./filer_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ event.proto --go_out=./event_pb --go-grpc_out=./event_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ remote.proto --go_out=./remote_pb --go-grpc_out=./remote_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ iam.proto --go_out=./iam_pb --go-grpc_out=./iam_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ mount.proto --go_out=./mount_pb --go-grpc_out=./mount_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ s3.proto --go_out=./s3_pb --go-grpc_out=./s3_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ mq.proto --go_out=./mq_pb --go-grpc_out=./mq_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
	protoc -I ./ schema.proto --go_out=./schema_pb --go-grpc_out=./schema_pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

fbs:
	flatc --go -o . --go-namespace message_fbs message.fbs
