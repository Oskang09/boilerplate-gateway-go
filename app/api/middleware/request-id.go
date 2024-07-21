package middleware

import (
	"math/rand"
	"time"

	"github.com/kataras/iris/v12"
)

func numeric(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letterRunes = []rune("1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	s := string(b)
	return s
}

func generateRequestID() string {
	return "W" + time.Now().Format("20060102150405") + "E" + numeric(9)
}

func RequestID() iris.Handler {
	return func(ctx iris.Context) {
		requestID := ctx.Request().Header.Get("x-request-id")
		if requestID == "" {
			requestID = generateRequestID()
		}
		ctx.Request().Header.Add("x-request-id", requestID)
		ctx.ResponseWriter().Header().Add("x-request-id", requestID)
		ctx.Next()
	}
}
