# Distributed-Architecture
Experimenting with Distributed Systems

Command to generate gRPC code of Protobuf message:
`make compile`

OR Alternative:
`protoc --proto_path=api/v1 --go_out=out --go_opt=paths=source_relative api/v1/log.proto`