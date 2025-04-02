.DEFAULT_GOAL := generate
.PHONY : generate server  client 
server : 
	cd server && go run cmd/*.go localhost:8081
client : 
	cd client && go run *.go localhost:8081
generate : 
	@echo  "Genrating protobuf and grpc files"
	protoc --go_out=. --go_opt=module=github.com/paschalolo/grpc --go-grpc_out=. --go-grpc_opt=module=github.com/paschalolo/grpc proto/dummy/v1/dummy.proto 