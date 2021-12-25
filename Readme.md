## grpc 初期設定

```powershell
go get google.golang.org/grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go mod init <dir>
go mod tidy
```

## 使用方法

- サーバーの作成

```powershell
go run cmd/server/serer.go
```

- ユーザーの作成

```powershell
go run cmd/client/client.go <URL> <USER>
```
