package handler

import (
	"gateway/app/repository"
	"gateway/app/response"

	"github.com/kataras/iris/v12"
)

func ExampleHandler(ctx iris.Context, repository *repository.Repository) (int, response.ApiResponse) {
	return iris.StatusOK, response.Item(ctx, "example handler invoked")
}
