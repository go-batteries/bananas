# bananas

or maybe I should have called it `go-bananas`. This is a cli tool to setup an
opinionated golang base project.

The opinion being, write request, response using annotations in `proto` and `service`, and generate
- docs
- request, response models
- controllers

Rest of it is common setup needed for web dev.

### batteries

- `server` with [echo](https://github.com/labstack/echo) and graceful termination
- `cli`
- `zerolog` for logger
- `config` setup
- `database` setup
    - defaults to sqlx and sqlite3
- `redis` connection setup
- `bcrypt` and ulid for hashing and id generation.
- `protos` this is where one writes their api, and the tool generates:
    - `request`, `response` models
    - `dummy controller`
    - automatic `openapiv2` doc generation

### usage

Make sure the `bananas` executable is in $PATH.

```shell
$> mkdir -p testproj; cd testproj
$> go mod init testproj
$> bananas init -n testproj
$> go mod tidy

$> cp -r examples/protos/web/hellow protos/web/
$> bananas gen:docs --path=./protos/web # generates swagger.yaml

$> # generate the pb.go files, not fully supported yet in server
$> bananas gen:structs --path=./protos/web 
```

### building

For linux: `make build.cli`
