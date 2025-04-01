.DEFAULT_GOAL := generate
.PHONY : generate server 
server : 
	cd server 
	go run *.go localhost:8081
generate : 
	@echo  "Genrating protobuf and grpc files"
	cd proto
	protoc --go_out=. --go_opt=module=github.com/Paschalolo/grpc/proto --go-grpc_out=. --go-grpc_opt=module=github.com/Paschalolo/grpc/proto dummy/v1/dummy.proto 