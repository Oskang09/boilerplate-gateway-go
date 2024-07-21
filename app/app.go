package app

import (
	"fmt"
	"gateway/app/api/middleware"
	v1 "gateway/app/api/v1"
	"gateway/app/bootstrap"
	"gateway/app/config"
	"gateway/app/response"
	"os"
	"reflect"
	"strings"

	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/middleware/recover"
)

func getLogger() *accesslog.AccessLog {
	ac := accesslog.New(os.Stdout)
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = true
	ac.BytesSent = true
	ac.BodyMinify = false
	ac.RequestBody = false
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler
	ac.SetFormatter(&accesslog.JSON{HumanTime: true})
	return ac
}

func Start(port string) {
	bs := bootstrap.New()

	// register validator translation
	validator := validator.New()
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	en_translations.RegisterDefaultTranslations(validator, config.ValidatorTranslator)

	app := iris.New()
	app.Use(getLogger().Handler)
	app.Validator = validator

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.StopWithJSON(iris.StatusNotFound, response.Error(ctx, "ROUTE_NOT_EXIST", "route doens't exists"))
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.StopWithStatus(iris.StatusOK)
	})

	// injecting depdedencies to container
	app.RegisterDependency(
		bs.Database,
		bs.Redis,
		bs.Redsync,
		bs.Repository,
	)

	authorizedApp := app.Party(
		"/",
		recover.New(),
		middleware.RequestID(),
	)

	/* for versioning you will just need to define as per below, and configure using injection container */
	authorizedApp.Party("/v1", middleware.Opentracing("example-v1")).ConfigureContainer(v1.RegisterRoutes)

	app.Listen(
		fmt.Sprintf(":%s", port),
		iris.WithOptimizations,
		iris.WithoutPathCorrectionRedirection,
	)
}
