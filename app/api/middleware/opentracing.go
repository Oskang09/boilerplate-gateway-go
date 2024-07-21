package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/kataras/iris/v12"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
)

func Opentracing(componentName string) iris.Handler {
	return func(ctx iris.Context) {
		tracer := opentracing.GlobalTracer()
		req := ctx.Request()

		operation := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
		span := tracer.StartSpan(operation)
		defer span.Finish()

		// request id
		{
			span.LogKV("request.id", ctx.Request().Header.Get("x-request-id"))
		}

		// headers loggign
		{
			requestHeaders, err := json.Marshal(ctx.Request().Header)
			if err == nil {
				span.LogKV("request.headers", requestHeaders)
			}
		}

		// body logging
		{
			requestBody, ok := sortJSON(ctx)
			if ok {
				span.LogKV("request.body", requestBody.String())
				ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(requestBody.Bytes()))
			}
		}

		carrierContext := injectCarrier(req.Context(), span, tracer)
		ctx.ResetRequest(req.WithContext(carrierContext))

		ctx.Next()

		responseStatus := ctx.ResponseWriter().StatusCode()
		if responseStatus > 299 {
			span.LogKV("error", true)
		}

		// http standard variables
		{

			ext.HTTPStatusCode.Set(span, uint16(responseStatus))
			ext.HTTPMethod.Set(span, req.Method)
			ext.HTTPUrl.Set(span, req.URL.Path)
			ext.Component.Set(span, componentName)
		}

		// response header
		{
			responseHeaders, err := json.Marshal(ctx.ResponseWriter().Header())
			if err == nil {
				span.LogKV("response.headers", responseHeaders)
			}
		}

		// response body
		{
			response := ctx.Recorder()
			span.LogKV("response.body", string(response.Body()))
		}
	}
}

func sortJSON(c iris.Context) (*bytes.Buffer, bool) {
	body := new(bytes.Buffer)
	ok := false

	if c.Request().Body != http.NoBody {
		reqBodyByte, _ := ioutil.ReadAll(c.Request().Body)
		mxjData, err := mxj.NewMapJson(reqBodyByte)
		if err == nil {
			data, err := mxjData.Json(true)
			if err == nil {
				json.Compact(body, data)
				ok = true
			}
		}
	}

	return body, ok
}

func injectCarrier(ctx context.Context, span opentracing.Span, tracer opentracing.Tracer) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	mdWriter := newCarrierMetadata(md)
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, mdWriter)
	if err != nil {
		md = metadata.MD{}
	}

	return metadata.NewOutgoingContext(
		opentracing.ContextWithSpan(ctx, span),
		md,
	)
}

type carrierMetadata struct {
	metadata.MD
}

func newCarrierMetadata(md metadata.MD) carrierMetadata {
	return carrierMetadata{md}
}

// Set :
func (w carrierMetadata) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

// ForeachKey :
func (w carrierMetadata) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
