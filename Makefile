all: gen build test vet

build: server

server:
	@cd cmd/server && go build -o ../../bin/server

gen:
	@buf generate --path proto/data.proto

vet:
	@go vet ./...

test:
	@go test ./...

clean:
	@rm -rf bin pkg/apipb

count:
	@echo ""
	@echo "WITH EVERYTHING"
	@gocloc  --include-lang="Go,Protocol Buffers" .
	@echo ""
	@echo "WITH ONLY THE FILES WE MAINTAIN"
	@gocloc --not-match-d pkg/apipb --include-lang="Go,Protocol Buffers" .
	
init:
	@go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/bufbuild/buf/cmd/buf
