gen-proto:
	protoc \
	--proto_path=protos \
	--go_out=internal/gen \
	--go_opt=paths=source_relative \
	--go-grpc_out=internal/gen \
	--go-grpc_opt=paths=source_relative \
	protos/card_validator_svc.proto

run:
	go run ./cmd

update-proto:
	git submodule update --recursive --remote

build:
	CGO_ENABLED=0 GOOS=linux GOFLAGS=-mod=vendor go build -o ./build/app -ldflags="-s -w" ./cmd

compress:
	upx --best --lzma ./build/app

test:
	go test $$(go list ./... | grep -v gen) -v -cover

.PHONY: gen-proto run update-proto build compress test