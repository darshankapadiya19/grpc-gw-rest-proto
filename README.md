stub generation command

```aiignore
protoc ./proto/hello.proto --openapiv2_out=. --go_out=. --go-grpc_out=. --grpc-gateway_out=. 
```

dependencies for swagger docs
```aiignore
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

post request
```aiignore
curl -X POST -d '{"name": "world"}' http://localhost:8080/hello
```