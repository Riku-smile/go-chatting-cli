grpc 初期設定

```
go get google.golang.org/grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go mod init <dir>
go mod tidy
```
