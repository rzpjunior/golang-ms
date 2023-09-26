# API DASHBOARD EDENFARM VERSION 2

[![build status](https://git.edenfarm.id/project-version2/api/badges/develop/build.svg)](https://git.edenfarm.id/project-version2/api/commits/develop)
[![coverage report](https://git.edenfarm.id/project-version2/api/badges/develop/coverage.svg)](https://git.edenfarm.id/project-version2/api/commits/develop)

API that will consumed by Dashboard.

## Installation

Requirement:
 1. Install Docker
 2. Install docker-compose

```bash
    1. go get git.edenfarm.id/project-version2/api
    2. go mod tidy
    3. go mod vendor
    4. migration up
```

## Database

How to use migration please see https://github.com/golang-migrate/migrate.

### Migrations

Database migrations are on migrations folder, see https://github.com/golang-migrate/migrate.

Example migrating:
```bash
migrate -url="mysql://username:password@tcp(127.0.0.1:3306)/eden(name of Database)" -path="./migrations" up
```

```bash
migrate -url="mysql://username:password@tcp(127.0.0.1:3306)/eden(name of Database)" -path="./migrations" down
```

### Running Test

Test can be executed by makefile
```bash
	make test -s
```

For fixing formating and lint on development mode
```bash
	make format -s
```

Completed command just read the help
```bash
	make
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a merge request :D
