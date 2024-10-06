build.cli:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o out/linux/bananas cmd/cli/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o out/linux/bananas_arm cmd/cli/main.go

	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o out/darwin/bananas cmd/cli/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o out/darwin/bananas_amd cmd/cli/main.go

	cp out/linux/bananas $(HOME)/.local/bin/

gen.web.proto:
	( protoc -I protos/web \
		-I protos/includes/googleapis \
		-I protos/includes/grpc_ecosystem \
		--openapiv2_out ./openapiv2 --openapiv2_opt logtostderr=true \
		./protos/web/**/*.proto \
	)

gen.api.proto:
	( protoc -I protos/web \
		-I protos/includes/googleapis \
		-I protos/includes/grpc_ecosystem \
		--go_out=./app/core --go_opt=paths=source_relative \
		--go-grpc_out=./app/core --go-grpc_opt=paths=source_relative \
		./protos/web/**/*.proto \
	)
