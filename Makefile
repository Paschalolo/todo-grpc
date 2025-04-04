.DEFAULT_GOAL := generate
.PHONY : generate server  client proto
server : 
	cd server && go run cmd/*.go 8081
client : 
	cd client && go run cmd/*.go 8081
generate : 
	@echo  "Genrating protobuf and grpc files"
	cd proto && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative todo/v2/*.proto
proto:
	cd helpers && protoc --go_out=. --go_opt=paths=source_relative proto/*.proto
proto-gen-validate :
	 protoc -I. \
    --go_out=. \
    --go_opt=paths=source_relative \  
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \  
    --validate_out="lang=go,paths=source_relative:." \
    todo/v2/todo.proto