build.cli:
	go build -o out/linux/bananas cmd/cli/main.go
	cp out/linux/bananas /home/darksied/.local/bin/
