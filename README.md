# bananas

![logo](./logo.jpg)

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

# generates swagger.json
$> bananas gen:docs 

# generate the pb.go files
$> bananas gen:structs 

$> go run cmd/server/main.go
```

#### to enable grpc use

```shell
$> bananas init -n testproj --grpc

# from here on, when need to update grpc structs and docs use

$> bananas gen:structs --grpc
$> bananas gen:docs

$> go run cmd/server/main.go
```

_N.B.: The `--path` flag for overriding protos/ directory is off. Will add that later_

A default hello example has been provided.
After running the server, test it with curl:

```shell
curl -XGET 'http://localhost:9001/hellow?name=DudePerfect'
curl -XDELETE 'http://localhost:9001/hellow?name=DudePerfect'
```

### building

`make build.cli`
