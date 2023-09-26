# erp-protobuf
EdenFarm protobuf test

## Installation Protoc Compiler

- Linux, using apt or apt-get, for example:
```
apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+
```

- MacOS, using Homebrew:
```
brew install protobuf
protoc --version  # Ensure compiler version is 3+
```

## How To Use

- Install protoc gen for golang in your project workspace
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
- Dont forget to export your go bin
```
export PATH="$PATH:$(go env GOPATH)/bin"
```

- Generate automaticaly files proto
```
make clean && make gen
# Or
make all
```

- All proto will be generated in folder `gen/proto/`

## Add To Your Project

- Clone the grpc protobuf in repository
```
git clone https://git.edenfarm.id/project-version3/erp-services/erp-protobuf
```
- Open your service project`
- Change module in your project `go.mod`
```
replace git.edenfarm.id/project-version3/erp-services/erp-protobuf v0.0.0 => ../erp-protobuf

require (
	git.edenfarm.id/project-version3/erp-services/erp-protobuf v0.0.0
)
```
- Update your vendor project
```
go mod vendor
```

- Add to your project, sample for `erp-boilerplate-service`
```
pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/boilerplate_service"
```

## Documentation

This is full [documentation](https://developers.google.com/protocol-buffers/docs/proto3) for protocol buffers. READ NOW !!!
