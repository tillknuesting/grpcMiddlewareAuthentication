generate:
	protoc -I/usr/local/include -I. \
      --go_out=plugins=grpc:. \
     proto/helloworld.proto
