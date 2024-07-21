package response

import (
	"gateway/app/config"
	"runtime"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/opentracing/opentracing-go"
)

type ApiResponse struct {
	Cursor string      `json:"cursor,omitempty"`
	Item   interface{} `json:"item,omitempty"`
	Error  *apiError   `json:"error,omitempty"`
}

type apiError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message,omitempty"`
	Validation map[string]string `json:"validation,omitempty"`
}

func newResponse(ctx iris.Context) ApiResponse {
	span := opentracing.SpanFromContext(ctx.Request().Context())
	if span != nil {
		// getPreviousFrame, newResponse, Error/Item/Items
		_, fileName, fileLine := getPreviousFrame(3)
		span.LogKV(
			"response.file", fileName,
			"response.line", fileLine,
		)
	}
	return ApiResponse{}
}

func Error(ctx iris.Context, code string, err interface{}) ApiResponse {
	result := newResponse(ctx)
	result.Error = new(apiError)
	result.Error.Code = code

	switch err.(type) {

	case validator.ValidationErrors:
		result.Error.Validation = make(map[string]string)
		errs := err.(validator.ValidationErrors)
		for _, err := range errs {
			result.Error.Validation[err.Field()] = err.Translate(config.ValidatorTranslator)
		}

	case string:
		result.Error.Message = err.(string)

	case error:
		result.Error.Message = err.(error).Error()

	}
	return result
}

func Item(ctx iris.Context, item interface{}) ApiResponse {
	result := newResponse(ctx)
	result.Item = item
	return result
}

func Items(ctx iris.Context, item interface{}, cursor string) ApiResponse {
	result := newResponse(ctx)
	result.Item = item
	result.Cursor = cursor
	return result
}

func getPreviousFrame(frame int) (string, string, int) {
	pointer, filename, line, ok := runtime.Caller(frame)
	if !ok {
		return "", "", 0
	}
	return runtime.FuncForPC(pointer).Name(), filename, line
}
