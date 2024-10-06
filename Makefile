build.cli:
	go build -o out/linux/bananas cmd/cli/main.go
	cp out/linux/bananas /home/darksied/.local/bin/

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
