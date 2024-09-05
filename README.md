# Work In Progress
![Hexago](doc/hexago.png)

## Installation

```sh
go install github.com/ksckaan1/hexago@latest
```

## Dependencies
- [go](https://go.dev)
- [impl](https://github.com/josharian/impl)

> Make sure that the directory `$HOME/go/bin` is appended to the `$PATH` ortan variable

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

## Usage

### `doctor`

The `doctor` command displays the status of dependencies that are required for hexago to run properly.

**Example:**

https://github.com/user-attachments/assets/eebc095e-b806-41b7-bad7-0fb04cd379c7



### `init`
The `init` command initialize a Hexago project. This command creates a domain named `core` by default. Promts go module name. If leaves blank, uses project folder name as lowercase defaultly.

```sh
hexago init <project-path>
```
**Example:**

https://github.com/user-attachments/assets/b4ffd97a-a77a-4545-ae3b-41b168f32221


### `domain`
This is the parent command for all domain-related operations.

If the project does not contain any domain, a new `service` and `app` cannot be created. For this, a domain must be created first.

- #### `new`

  This command creates a new domain under the `internal/domain` directory.

  ```sh
  hexago domain new
  ```

  https://github.com/user-attachments/assets/09c5775b-39e6-47eb-bd46-090d2a07843a

- #### `ls`
  
  This command lists all domains under the `internal/domain` directory.

  ```sh
  hexago domain ls
  ```
  **Flags:**
  - `-l`: lists domains line-by-line
  
  https://github.com/user-attachments/assets/7014cbb4-5730-4278-affa-1a2dbc6d6ba5

### `port`
This is the parent command for all port-related operations.

Ports created in domains can be implemented when creating service, app, infrastructure and package. If there is no port in the project, it is not asked which port to implement in the creation screen.

You can create a port manually like bellow.

```go
// internal/domain/<domainname>/port/<portfilename>.go

package port

type ExamplePort interface {
  Create(ctx context.Context) error
  GetAll(ctx context.Context) ([]string, error)
}
```

You can use this port when creating a new service, app, infrastructure or package.

- #### `ls`:

  This command lists all ports under the `internal/domain/<domainname>/port`

  **Flags:**
  - `-l`: lists ports line-by-line

  https://github.com/user-attachments/assets/6898c794-ab31-4ed8-9eee-1702a566f655



### `service`
This is the parent command for all service-related (domain-service) operations.

- #### `new`

  This command creates a new service under the `internal/domain/<domainname>/service/<servicename>` directory.

  Domain is required to create a service. Steps applied when creating a service:

  - Insert service name (PascalCase)
  - Insert folder name (lowercase)
  - Select a domain
  - Select port which will be implemented (skips this step if there is no port)
  - Assert port if selected

  ```sh
  hexago service new
  ```

  https://github.com/user-attachments/assets/b3f94d5f-aa05-4b61-842c-37153b901328

- #### `ls`

  This command lists all services under the `internal/domain/<domainname>/service` directory.

  ```sh
  hexago service ls
  ```
  **Flags:**
  - `-l`: lists services line-by-line
  
  https://github.com/user-attachments/assets/39607acd-aed9-47cf-baed-3c7bd2ab5bce

### `app`
This is the parent command for all application-related (application-service) operations.

Application services are the places where endpoints such as controllers or cli applications are hosted.

- #### `new`

  This command creates a new application under the `internal/domain/<domainname>/app/<appname>` directory.

  Domain is required to create an application. Steps applied when creating an application:

  - Insert application name (PascalCase)
  - Insert folder name (lowercase)
  - Select a domain
  - Select port which will be implemented (skips this step if there is no port)
  - Assert port if selected

  ```sh
  hexago app new
  ```

  https://github.com/user-attachments/assets/a390cb4b-91f5-45dd-a4df-a076765452c9

- #### `ls`

  This command lists all applications under the `internal/domain/<domainname>/app` directory.

  ```sh
  hexago app ls
  ```
  **Flags:**
  - `-l`: lists applications line-by-line
  
  https://github.com/user-attachments/assets/8cd60430-4c83-4f42-a252-462974fa635e

### `infra`
This is the parent command for all infrastructure-related operations.

Infrastructures host databases (repositories), cache adapters or APIs that we depend on while writing applications

- #### `new`

  This command creates a new infrastructure under the `internal/infrastructure/<infraname>` directory.

  Steps applied when creating an infrastructure:

  - Insert infrastructure name (PascalCase)
  - Insert folder name (lowercase)
  - Select port which will be implemented (skips this step if there is no port)
  - Assert port if selected

  ```sh
  hexago infra new
  ```

  https://github.com/user-attachments/assets/8c000cfa-1f6b-42dd-8459-adf25182d972

- #### `ls`

  This command lists all infrastructures under the `internal/infrastructure` directory.

  ```sh
  hexago infra ls
  ```
  **Flags:**
  - `-l`: lists infrastructures line-by-line
  
  https://github.com/user-attachments/assets/47a1379d-3661-4e8c-88f3-13cea15bcf12

### `pkg`
This is the parent command for all package-related operations.

Packages are the location where we host features such as utils. There are two types of packages in a hexago project. 
- The first one is located under `/internal/pkg` and is not imported by other go developers. Only you use these packages in the project.
- The second is located under `/pkg`. The packages here can be used both by your project and by other go developers.

- #### `new`

  This command creates a new package under the `internal/pkg/<pkgname>` or `/pkg/<pkgname>` directory.

  Steps applied when creating a package:

  - Insert package name (PascalCase)
  - Insert folder name (lowercase)
  - Select port which will be implemented (skips this step if there is no port)
  - Assert port if selected
  - Select package scope (global or internal)

  ```sh
  hexago pkg new
  ```

  https://github.com/user-attachments/assets/1d518553-49ce-4c18-868a-f7ff87829a36

- #### `ls`

  This command lists all packages under the `internal/pkg` or `/pkg` directory.

  ```sh
  hexago pkg ls # list internal packages
  ```
  **Flags:**
  - `-g`: lists global packages
  - `-a`: list both global and internal packages.
  - `-l`: lists packages line-by-line
  
  https://github.com/user-attachments/assets/6d055ce5-deff-4096-8f8c-00964238cc59

### `cmd`
This is the parent command for all entry point-related (cmd) operations.

Entry points are the places where a go application will start running. entry points are located under the `cmd` directory.

- #### `new`

  This command creates a new entry point under the `cmd/<entry-point-name>` directory.

There is only one step creating an entry point.

  - Insert entry point folder name (kebab-case)

  Creates a go file like bellow.

  ```go
  package main

  func main(){

  }
  ```

  ```sh
  hexago cmd new
  ```

  https://github.com/user-attachments/assets/63503d4a-7691-4855-afa4-672602904d96

- #### `ls`

  This command lists all entry points under the `cmd` directory.

  ```sh
  hexago cmd ls
  ```
  **Flags:**
  - `-l`: lists entry points line-by-line
  
  https://github.com/user-attachments/assets/2f543ce7-fdaf-4970-9324-59027dcca179

  ### `run`
This command can be used for two different purposes. the `run` command create a log file under the `logs` directory defaultly.

- Firstly, if there is an entry point in your project, it can be used to run this entry point.

  ```sh
  hexago run <entry-point-name>
  ```

  https://github.com/user-attachments/assets/3f80bcd9-b737-4f79-8635-ad5a220c7680

  **Flags:**
  - `-e`: run entry point with environment variable. You can use multiple environment variable
    ```sh
    hexago run <entry-point-name> -e <ENV_KEY1>=<ENV_VALUE1> -e <ENV_KEY2>=<ENV_VALUE2>
    ```

  https://github.com/user-attachments/assets/ef73dfbe-79b3-482b-92d0-dbe0a29738f7
  
  You can customize this run command with given entry point in `.hexago/config.yaml` file.

  You can specify all envs from `config.yaml` file like bellow.

  ```yaml
  templates: # std | do | <custom>
  service: std
  application: std
  infrastructure: std
  package: std

  runners:
    api: # it runs "go run ./cmd/api", if exists
      env:
        - ENV_KEY1=ENV_VAL1
        - ENV_KEY2=ENV_VAL2
      log:
        disable: false # write logs to files
        seperate_files: true # create log files seperately as api.stderr.log and api.stdout.log
        overwrite: true # create new log file when runner called
  ```

  When the `hexago run api` command is executed as above, it starts the `api` entry point according to the settings in the `config.yaml` file.

- As a second method, you can use the `run` sub-command as an alternative to makefile. You can create a new entry in the `runners` section of the `.hexago/config.yaml` file to call it with the `run` command.

  The special commands created do not need to have an entry point equivalent. We can add a special command using the `cmd` key.

  ```yaml
  runners:
  custom-command:
    cmd: "go version" # overwrite default "go run ./cmd/mycommand/" command
    log:
      disabled: true # do not print log file
  ```

  When you run `hexago run custom-command` command, you will get the following result.

  ```text
  go version go1.23.0 darwin/arm64
  ```
