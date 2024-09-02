# Work In Progress
![Hexago](doc/hexago.png)

## Installation

```sh
go install github.com/ksckaan1/hexago@latest
```

## Get Started

If you didn’t hear about hexagonal architecture before, firstly, you could research about it.

Here it is nice blog posts about hexagonal architecture:

- https://medium.com/ssense-tech/hexagonal-architecture-there-are-always-two-sides-to-every-story-bc0780ed7d9c
- https://dev.to/bagashiz/building-restful-api-with-hexagonal-architecture-in-go-1mij
- https://medium.com/@janishar.ali/how-to-architecture-good-go-backend-rest-api-services-14cc4730c05b

## Example Folder Structure

```
.
├── .hexago/
│   ├── config.yaml
│   └── templates/
├── cmd/
│   └── api/
│       └── main.go
├── config/
├── doc/
├── go.mod
├── internal/
│   ├── domain/
│   │   └── core/
│   │       ├── application/
│   │       │   ├── restapi/
│   │       │   │   └── restapi.go
│   │       │   └── rpcapi/
│   │       │       └── rpcapi.go
│   │       ├── dto/
│   │       │   └── order.go
│   │       ├── model/
│   │       │   └── order.go
│   │       ├── port/
│   │       │   └── order.go
│   │       └── service/
│   │           ├── cancelorder/
│   │           │   └── cancelorder.go
│   │           ├── createorder/
│   │           │   └── createorder.go
│   │           ├── getorder/
│   │           │   └── getorder.go
│   │           ├── listorders/
│   │           │   └── listorders.go
│   │           └── updateorder/
│   │               └── updateorder.go
│   ├── infrastructure/
│   │   ├── cache/
│   │   │   └── cache.go
│   │   └── orderrepository/
│   │       └── orderrepository.go
│   └── pkg/
│       ├── authtoken/
│       │   └── authtoken.go
│       └── uniqueid/
│           └── uniqueid.go
├── pkg/
├── schemas/
└── scripts/
```
