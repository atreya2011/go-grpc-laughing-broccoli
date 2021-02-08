# go-grpc-laughing-broccoli

Based on [grpc-gateway-boilerplate](https://github.com/johanbrandhorst/grpc-gateway-boilerplate.git)

## Running

Running `main.go` starts a web server on <https://0.0.0.0:11000/>. You can configure
the port used with the `$PORT` environment variable, and to serve on HTTP set
`$SERVE_HTTP=true`.

```shell
go run main.go
```

An OpenAPI UI is served on <https://0.0.0.0:11000/>.

### Running the standalone server

If you want to use a separate gRPC server, for example one written in Java or C++, you can run the
standalone web server instead:

```shell
go run ./cmd/standalone/ --server-address dns:///0.0.0.0:10000
```

## Getting started

After cloning the repo, there are a couple of initial steps;

1. Install the generate dependencies with `make install`.
   This will install `buf`, `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway`,
   `protoc-gen-openapiv2` and `statik` which are necessary for us to generate the Go, swagger and static files.

2. Finally, generate the files with `make generate`.

Now you can run the web server with `go run main.go`.
