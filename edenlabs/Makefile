test:
	go test -v -cover -covermode=atomic ./...
	
update:
	go get git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc@latest
	go mod tidy
	go mod vendor