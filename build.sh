protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/githubtasks.proto
mv proto/github.com/brotherlogic/githubtasks/proto/* ./proto
