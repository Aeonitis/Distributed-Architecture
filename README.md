# Distributed-Architecture
Experimenting with Distributed System concepts via Go

###### Log Package
* Record — data stored in our log.
* Store — file for store records.
* Index — file for store index entries.
* Segment — abstraction that ties a store & an index.
* Log — abstraction that ties all the segments together.

Command to generate gRPC code of Protobuf message:
`make compile`

OR Alternative:
`protoc --proto_path=api/v1 --go_out=out --go_opt=paths=source_relative api/v1/log.proto`