# ERP Inventory Service 
ERP Inventory Service - Microservice EdenFarm And Integration With ERP GP

## Contents
- [Requirement](#requirement)
- [Migration](#migration)
- [Local Infrastructure](#local-infrastructure)
- [Tools](#tools)
- [Testing](#testing)

## Requirement
- Go 1.18+
- Docker
- Git
- VS Code / Goland
- Coffee & Cake

## Local Infrastructure
- Running your local infrastructure :
    - MySQL
    - MongoDB
    - Redis
    - Elasticsearch
    - Kafka
    - Jaeger
    - Prometheus
    - OtelCol

## Migration
- Open [Migration Project](https://git.edenfarm.id/project-version3/erp-databases)
- Follow the instruction

## Tools
- Install [Edenlabs CLI](https://git.edenfarm.id/edenlabs/cli)
- Start Debug & Run your instance
```
edenlabs run
# or
edenlabs grpc
# or
edenlabs consumer
```

## Testing
- Install mockery for automated create mocks
- Create file test
- Run Test
