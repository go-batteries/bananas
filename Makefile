OUT_DIR=out

build.cli:
	go build -ldflags "-s -w" -o $(OUT_DIR)/bananas cmd/cli/main.go
	cp $(OUT_DIR)/bananas $(HOME)/.local/bin/

build.all:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o tmp/linux/bananas cmd/cli/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o tmp/linux/arm/bananas cmd/cli/main.go

	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o tmp/darwin/bananas cmd/cli/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o tmp/darwin/amd/bananas cmd/cli/main.go

zip:
	mkdir -p releases
	zip releases/bananas.linux_amd64.zip tmp/linux/bananas
	zip releases/bananas.darwin_amd64.zip tmp/darwin/amd/bananas
	zip releases/bananas.darwin_arm64.zip tmp/darwin/bananas
	zip releases/bananas.linux_arm64.zip tmp/linux/arm/bananas

release.cli: build.all zip

gen.web.proto:
	( protoc -I protos/web \
		-I protos/includes/googleapis \
		-I protos/includes/grpc_ecosystem \
		--openapiv2_out ./openapiv2 --openapiv2_opt logtostderr=true \
		./protos/web/**/*.proto \
	)

gen.grpc.proto:
	( protoc -I protos/web \
		-I protos/includes/googleapis \
		-I protos/includes/grpc_ecosystem \
		--go_out=./app/core --go_opt=paths=import \
		--go-grpc_out=./app/core --go-grpc_opt=paths=import \
		--grpc-gateway_out=./app/core --grpc-gateway_opt=paths=import \
		./protos/web/**/*.proto \
	)

gen.apbodyi.proto:
	( protoc -I protos/web \
		-I protos/includes/googleapis \
		-I protos/includes/grpc_ecosystem \
		--go_out=./app/core --go_opt=paths=import \
		--go-grpc_out=./app/core --go-grpc_opt=paths=import \
		./protos/web/**/*.proto \
	)
