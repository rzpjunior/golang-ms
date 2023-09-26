# Client GRPC
This repository using for client GRPC. 
- All the client must create in this repository.

## Add To Your Project

- Clone the client grpc in repository
```
git clone git@git.edenfarm.id:project-version3/erp-pkg/erp-client-grpc.git
```

- Change module in your project `go.mod`
```
replace git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc v0.0.0 => ../erp-client-grpc

require (
	git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc v0.0.0
)
```
- Update your vendor project
```
go mod vendor
```


## How To Add new client

- Create your own connection at 
```bash
client->name_func.go
```
- Create function services at 
```bash
service->create_name_folder->create_name.go
```

## How To use in your Services

- Add new line to function `start()` in `api.go`
```bash
err = client.ConnectClientNameYourFunction()
```
- Call function service
```bash
 sNameService := name_service.ServiceNameYourFunctionGrpc()
 sNameService.Create(params)
```
