package routes

import (
	"github.com/valyala/fasthttp"
	"net/http"
)

func GetHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write([]byte("{\"status\": \"ok\"}"))
}
