all: gen build test vet

build: server client

server:
	@cd cmd/$@ && go build -o ../../bin/$@

client:
	@cd cmd/$@ && go build -o ../../bin/$@

gen:
	@buf generate

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
	
dep-install:
	@go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go get github.com/bufbuild/buf/cmd/buf@latest
	@go get github.com/mgechev/revive@latest
	@buf mod update

