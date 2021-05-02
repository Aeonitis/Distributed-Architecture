# Distributed-Architecture
Experimenting with Distributed Systems

Command to generate gRPC code of Protobuf message (w/libprotoc 3.15.8):
`protoc --proto_path=api/v1 --go_out=out --go_opt=paths=source_relative api/v1/log.proto`