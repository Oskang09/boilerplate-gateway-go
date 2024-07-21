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
