# Fusion-framework fo Golang
### Notice: This information ralates only for fusion-framework

## Preface
GoLang is a great programming language for building Web Projects. I wanted to build a golang framework on httprouter after real benchmark comparison on my slow old netbook (gin, fasthttp and ozzo-routing were is slower than httprouter unexpectedly). First of all, it will be a framework for high-load web applications, and for resist network threats.

## Introduction
I make some changes in the https://github.com/qiangxue/go-rest-api template according web development needs. 
To get an idea of the features to be included in the project of this article provides a number of examples on how these features can be implemented ([Git Repo](https://github.com/brianwoo/server_side_dev_with_golang)). Structure of the project were be inspired from ozzo framework and https://github.com/qiangxue/go-rest-api former author of php yii framework and golang project layout [https://github.com/golang-standards/project-layout].

## Requirements

- Processor x32 or x64 intel compatible
- linux/bsd bash command compatible iterface
- rsync, make, git installed
- Go 1.16 up
- MySql 5.7 up
- libjpeg* installed for Image progressive jpeg optimisation
- root requirement provileges for remote server for deployment
- ssh key authentication connection established

Now features:
- [x] RESTful accepted format
- [x] CRUD operations for a database one table
- [x] JWT authentication
- [x] Environment dependent ozzo-config configuration management
- [x] ozzo-validation library
- [x] Structured logging with contextual information
- [x] Error handling with proper error response generation
- [x] ozzo-dbx database library
- [x] Database migration
- [x] Data validation
- [x] Test coverage
- [x] Live reloading during development
- [x] change db type to mysql (without dockerize this one)
- [x] Makefile for development
- [x] Golang standart structure

Todo change:
- [x] Healthchecks endpoints
- [ ] Default html template
- [ ] apply html templates/web-forms for manage records in db
- [ ] CSP Resource definitions
- [ ] Replace jwt library from drigvaila to nats.NKEYS
- [ ] JWT authentication in the cookie HttpOnly store
- [ ] Uploading files by secure pipelining
- [ ] Apply nats.io for async pipelines queues
- [ ] OAuth2 with Google
- [ ] OAuth2 with Facebook
- [ ] Integration frontend development pipeline
- [ ] Crud RESTAPI generator
- [ ] Migrations
- [ ] Docker implementation for development pipeline

The framework uses the following Go packages which can be replaced with your own favorite ones since their usages are mostly localized and abstracted. 

* Routing: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing)
* Database access: [ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
* Database migration: [golang-migrate](https://github.com/golang-migrate/migrate)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [zap](https://github.com/uber-go/zap)
* JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

# Building a Web Project with Fusion-framework

This fusion-framework is designed to get you up and running with a project structure optimized for HTML webapp developing and RESTful API services with Go. It promotes the best practices that follow the [SOLID principles](https://en.wikipedia.org/wiki/SOLID)
and [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). 
It encourages writing clean and idiomatic Go code. 

## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to
install Go on your computer. The kit requires **Go 1.15 or above**.

[Docker](https://www.docker.com/get-started) is also needed if you want to try the kit without setting up your
own database server. The kit requires **Docker 17.05 or higher** for the multi-stage build support.

After installing Go and Docker, rename config/_dev.yml to configs/dev.yml them fill appropriately. Run the following commands to start experiencing this fusion-framework:

```shell
# download the fusion-framework
git clone https://github.com/tvitcom/fusion-framework.git

cd fusion-framework

# start a PostgreSQL database server in a Docker container
make db-start

# seed the database with some test data
make testdata

# run the RESTful API server
make run

# Or develepment:
make dev

# or run the API server with live reloading, which is useful during development
# requires fswatch (https://github.com/emcrisostomo/fswatch)
make run-live
```

should return a list of album records in the JSON format

```
FAt this time, you have a RESTful API server running at `http://127.0.0.1:3000`. It provides the following endpoints:

* `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `POST /v1/login`: authenticates a user and generates a JWT
* `GET /v1/albums`: returns a paginated list of the albums
* `GET /v1/albums/:id`: returns the detailed information of an album
* `POST /v1/albums`: creates a new album
* `PUT /v1/albums/:id`: updates an existing album
* `DELETE /v1/albums/:id`: deletes an album

Try the URL `curl http://localhost:3000/healthcheck` in a browser, and you should see something like `"OK 0.19 0.29 0.43 2/885 27692"` displayed.

If you have `cURL` or some API client tools (e.g. [PostmanCanary](https://www.postman.com/downloads/canary/)), you may try the following 
more complex scenarios:

To use the fusion-framework as a starting point of a real project whose package name is `github.com/abc/xyz`, do a global 
replacement of the string `github.com/tvitcom/fusion-framework` in all of project files with the string `github.com/abc/xyz`.

## Project Layout

The fusion-framework uses the following project layout:
 
```shell
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── configs               configuration files for different environments
├── internal             private application and library code
│   ├── album            album-related features
│   ├── auth             authentication feature
│   ├── config           configuration library
│   ├── entity           entity definitions and domain logic
│   ├── errors           error types and handling
│   ├── healthcheck      healthcheck feature
│   └── test             helpers for testing purpose
├── migrations           database migrations
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── graceful         graceful shutdown of HTTP server
│   ├── log              structured and context-aware logger
│   └── pagination       paginated list
└── testdata             test data scripts
```

The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within `internal` and `pkg`, packages are structured by features in order to achieve the so-called
[screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). For example, 
the `album` directory contains the application logic related with the album feature. 

Within each feature package, code are organized in layers (API, service, repository), following the dependency guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).


## Common Development Tasks

This section describes some common development tasks using this fusion-framework.

### Implementing a New Feature

Implementing a new feature typically involves the following steps:

1. Develop the Agregator(or Service) that implements the business logic supporting the feature. Please refer to `internal/album/Agregator(or Service).go` as an example.
2. Develop the RESTful API exposing the Agregator(or Service) about the feature. Please refer to `internal/album/api.go` as an example.
3. Develop the repository that persists the data entities needed by the Agregator(or Service). Please refer to `internal/album/repository.go` as an example.
4. Wire up the above components together by injecting their dependencies in the main function. Please refer to 
   the `album.RegisterHandlers()` call in `cmd/server/main.go`.

### Managing Configurations

The application configuration is represented in `internal/config/config.go`. When the application starts,
it loads the configuration from a configuration file as well as environment variables. The path to the configuration 
file is specified via the `-config` command line argument which defaults to `./config/dev.yml`. Configurations
specified in environment variables should be named with the `APP_` prefix and in upper case. When a configuration
is specified in both a configuration file and an environment variable, the latter takes precedence. 

The `config` directory contains the configuration files named after different environments. For example,
`config/dev.yml` corresponds to the local development environment and is used when running the application 
via `make run`.

Do not keep secrets in the configuration files. Provide them via environment variables instead. For example,
you should provide `Config.DSN` using the `APP_DSN` environment variable. Secrets can be populated from a secret
storage (e.g. HashiCorp Vault) into environment variables in a bootstrap script (e.g. `cmd/server/entryscript.sh`). 

## Deployment

The application can be run as a docker container. You can use `make build-docker` to build the application 
into a docker image. The docker container starts with the `deployments/entryscript.sh` script which reads 
the `APP_ENV` environment variable to determine which configuration file to use. For example,
if `APP_ENV` is `prod`, the application will be started with the `config/prod.yml` configuration file.

You can also run `make build` to build an executable binary named `server`. Then start the API server using the following
command,

```shell
./server -config=./config/prod.yml
```
