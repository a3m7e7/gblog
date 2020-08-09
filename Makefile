.PHONY: internal pkg proto

proto:
	protoc -Iproto/ proto/*.proto --go_out=plugins=grpc:pkg
