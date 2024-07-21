# Table Of Contents

> Note: this project is used for gateway which exposing the API, for grpc-services you may refer to another repository which mainly for services development.

- [Table Of Contents](#table-of-contents)
- [Installation](#installation)
- [Project Setup](#project-setup)
- [Documentation](#documentation)
  - [Config \& Environment Variables](#config--environment-variables)
  - [API Version Control](#api-version-control)
  - [Request Handler](#request-handler)
  - [Response Handler](#response-handler)
  - [Model](#model)
  - [Repository](#repository)
  - [Extra: Background](#extra-background)

# Installation

If your primary go version was 1.20, you may skip this installation steps.

1. Install SDK

```
$ go install golang.org/dl/go1.20@latest
$ go1.20 download
$ go1.20 mod tidy
```

2. Setup Development Environment

If you're using vscode you can open "command + P" + select "> Go: Choose Go Enviroment" and select go1.20.

> Note: using `go version` check it's using 1.20 or your system primary version, if not the one you select please remove the GOPATH & PATH settings for Go SDK in your `.zshrc` / `.bashrc` / `.bashprofile`, and restart should working now.

# Project Setup 

Before you start the development please proceed steps as per below to setup your project specific settings & configuration. There's might some differences for CI setup if you're using bitbucket + circleci instead gitlab. This tutorial is mainly for gitlab repositories.

1. Go to `.vscode/launch.json`, and replace `{project-name}` with your project names.
2. Go to `main.go`, and replace `{port}` to your favorite port 
3. Go to `deployment/sandbox/kubernetes.yaml`, and add in your deployment information.
4. Go to `deployment/production/kubernetes.yaml`, and add in your deployment information.
5. Go to `deployment/Dockerfile`, and based on your needs setup `git config` for the go modules fetching and also `GOPRIVATE` configuration.
6. Go to `.gitlab-ci.yml`
   - Replace `{docker-image-path}` to your docker image path example - `registry-intl-vpc.ap-southeast-3.aliyuncs.com/dinar/wallet-engine` 
   - Replace `{project-name}` with your project name 
7. Go to `Makefile`
   - Replace `{project-name}` with your project name
   - Setup `GOPRIVATE` configuration with your needs
8. Go to `sonar-project.properties`, and replace `{project-name}` with your project name.
9. Go to `gitlab.revenuemonster.my` find your projects and update environment variables
   - set `KUBE_CREDENTIALS` to kubernetes base64 encoded kubectl configuration yaml
   - set `ALIYUN_DOCKER_USERNAME` to docker username
   - set `ALIYUN_DOCKER_PASSWORD` to docker password
   - set `SONAR_HOST_URL` to `https://sqube.superapp.my`
   - set `SONAR_TOKEN` to sonarqube token
10. Replace `{project-name}` inside `app/bootstrap/opentracing.go`.

# Documentation

Before you're proceed the development please make sure you know what are IoC ( Inversion of Control ) and Di ( Dependency Injection ), you can refer to [Martin Fowler's article](https://martinfowler.com/articles/injection.html#InversionOfControl). This repository come with some built-in features like versioning control, response, opentracing, goroutine with retry control so you can remove if some features you're not using but most of the time you will need all of them since it's quite a utility for all the projects.

## Config & Environment Variables

All environment variables will be defined in the config instead, you can refer to `app/config/config.go` for the example. For local development using vscode you may add your environment variables under `.vscode/launch.json` and there's a property key `configurations.*.env`.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "{project-name}",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "cwd": "${workspaceFolder}",
            "env": {
                "ENV": "local",
                "REDIS_HOST": "127.0.0.1:6379",
                "REDIS_PASSWORD": "",
                "MYSQL_HOST": "127.0.0.1",
                "MYSQL_PORT": "3306",
                "MYSQL_USER": "root",
                "MYSQL_PASSWORD": "",
                "MYSQL_DATABASE": "example"
                // define your own environment variables
            }
        }
    ]
}
```

## API Version Control

Versioning is expected to be done by seperate the folders `app/api/v1` mean version 1. Registering routes and the dependency container you may refer to `app/app.go` and code like `authorizedApp.Party("/v1", middleware.Opentracing("example-v1")).ConfigureContainer(v1.RegisterRoutes)`.

Ok so for the `*iris.Container` is the IoC Container for injecting depdendency so all your handler defined using this container will able to perform depdencny inject from the handler parameters. You can also inject the dependency based on the versioning instead of the globally by calling `RegisterDepedency` via `*iris.APIContainer` instead of `*iris.Application`. Here introduction and example from [Iris Documentation](https://docs.iris-go.com/iris/dependency-injection/context-register-dependency).

Example router definition as per below :-

```go
package v1

import (
	"gateway/app/api/v1/handler"
	"github.com/kataras/iris/v12"
)

func RegisterRoutes(app *iris.APIContainer) {
	user := app.Party("/example")
	{
		user.Post("", handler.ExampleHandler)
	}
}
```

## Request Handler

As per previous topic mention, you will able inject all the dependencies you're injected when calling the `RegisterDependency`, so the function signature will be dynamic and based on your needs. There's also come with some default injectable resources as per documented at [Iris Documentation](https://docs.iris-go.com/iris/dependency-injection/inputs). Example function signature as per below :-

```go
package handler

import (
	"gateway/app/repository"
	"gateway/app/response"
	"gateway/app/package/redsync"
	"github.com/kataras/iris/v12"
)

func ExampleHandler(ctx iris.Context, repository *repository.Repository) (int, response.ApiResponse) {
	return iris.StatusOK, response.Item(ctx, "example handler invoked")
}

func ExampleHandlerWithRedsync(ctx iris.Context, repository *repository.Repository, redsync *redsync.Redsync) (int, response.ApiResponse) {
	return iris.StatusOK, response.Item(ctx, "example handler invoked")
}
```

## Response Handler

We will suggest follow the function signature we're define since it's have better code readability. You might not following but make sure you're understand what you're doing, there's come with some default response parameters combination as per documented at [Iris Documentation](https://docs.iris-go.com/iris/dependency-injection/outputs).

And there's a package under `app/response` with opentracing, and validation translation supported so you can have your own response if you wish not to use the built-in one. So for the built-in usage will be quite straight forward :-

1. Response one result `response.Item(ctx, json)`
1. Response array result `response.Items(ctx, arrayJson, cursor)`
1. Response error message `response.Error(ctx, errCode, err)`

```go
package handler

import (
	"gateway/app/repository"
	"gateway/app/response"
	"gateway/app/package/redsync"
	"github.com/kataras/iris/v12"
)

func ExampleHandler(ctx iris.Context) (int, response.ApiResponse) {
	return iris.StatusOK, response.Item(ctx, "example handler invoked")
}

func ExampleHandlerArray(ctx iris.Context) (int, response.ApiResponse) {
	return iris.StatusOK, response.Items(ctx, []string("some message","some message"), "")
}

func ExampleHandlerError(ctx iris.Context) (int, response.ApiResponse) {
    
    var i struct {
        Data `json:"data" validate:"required"`
    }

    if err := ctx.ReadBody(&i); err != nil {
        return iris.StatusUnprocessableEntity, response.Error(ctx, "VALIDATION_ERROS", err)
    }
    
	return iris.StatusInternalServerError, response.Error(ctx, "INTERNAL_ERROR", "expected to be failed")
}
```

## Model

For using the base repository you will need to define the interface method which using for the repsitory, example of the model definition as per below :-

```go
package model

import (
	"time"

	"github.com/RevenueMonster/sqlike/types"
)

type Example struct {
	Key             *types.Key
	Name            string
	CreatedDateTime time.Time
	UpdatedDateTime time.Time
}

func (Example) Table() string {
	return "Example"
}
```

After defined the models, you can define your repository struct for example will be create under `app/repository/example.go` and you can go to `app/repository/index.go` and under the `New()` method to add your new defined model and come with basic crud operation, or even extended repository methods.

```go
// app/repository/example.go
package repository

import (
	"context"
	"gateway/app/model"

	"github.com/RevenueMonster/sqlike/sql/expr"
	"github.com/RevenueMonster/sqlike/sqlike/actions"
)

type Example struct {
	baseRepository[model.Example]
}

// Extended Methods
func (ex Example) FindByName(ctx context.Context, name string) (*model.Example, error) {
	result := ex.table.FindOne(ctx, actions.FindOne().Where(
		expr.Equal("Name", name),
	))

	example := new(model.Example)
	if err := result.Decode(example); err != nil {
		return nil, err
	}
	return example, nil
}

```

```go
// app/repository/index.go
package repository

import (
	"gateway/app/model"

	"github.com/RevenueMonster/sqlike/sqlike"
)

type tableContext interface {
	Table(string) *sqlike.Table
}

type Repository struct {
	Example Example
}

func New(db tableContext) *Repository {
	return &Repository{
		Example{newRepository[model.Example](db)},
	}
}
```

## Repository 

There's a abstract repository which come with the default CRUD operation you may refer to `app/repository/repository.go`, or the function signature as per below :-

> `T` type implements `model.TableModel` interface and interface define look like 
> ```go
>  type TableModel interface {
>	Table() string
> }
> ```


1. FIND BY PRIMARY KEY: `Find(ctx context.Context, key string) (*T, error)`
2. INSERT: `Create(ctx context.Context, model *T) (err error)`
3. CREATE TABLE IF NOT EXIST: `Migrate(ctx context.Context) (err error)`
4. UPDATE IF NOT EXIST ELSE INSERT: `Upsert(ctx context.Context, model *T) (err error)`
5. DELETE BY PRIMARY KEY: `Delete(ctx context.Context, model *T) (err error)`
6. PAGINATE WITH CURSOR: `Paginate(ctx context.Context, opts *PaginateOptions) ([]*T, string, error)`, and example of the pagination options definition

```go
type PaginateOptions struct {
	Limit   uint
	Cursor  string
	Queries []interface{}
	Sorts   []interface{}
}

opts := new(PaginateOptions)
opts.Limit = 50
opts.Cursor = ""
opts.Queries = []interface{}{
    expr.Equal("Key", 1),
}

opts.Sorts = []interface{}{
    expr.Desc("UpdatedDateTime")
}
```

## Extra: Background 

Because goroutine is not that friendly to developer who come from other language like Java, C# so here the similar package like other language but work for goroutine task. It come with extended retry mechanism using [avast/retry-go](https://github.com/avast/retry-go). Example usage as per below :-

```go
package main

import (
    "time"
    "github.com/avast/retry-go"
    "gateway/package/background"
)

func main() {
    //`go` is important to perform task in goroutine
    go background.RunTask(
        func(ctx context.Context, s opentracing.Span) error {
            return nil
        },
        // name for the task ( for opentracing tracker )
        background.WithName("risk-detection"), 

        // use when you're passing iris context
        background.WithIrisContext(ctx), 

        // use when you're passing generic context
        background.WithParentContext(ctx), 

        // timeout for the task
        background.WithTimeout(time.Second*60), 

        // opentracing logging fields
        background.WithLogs(
            log.String("some-id", "id"),
            log.String("some-id", "id"),
            log.String("some-id", "id"),
        ),
        
        // retry mechanism: https://github.com/avast/retry-go
        background.WithRetryOptions(
            retry.Attempts(6),
            retry.Delay(5000),
        ),
    )
}
```