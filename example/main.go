package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

const port = 8081

type UserInfo struct {
	Subject      string   `json:"sub"`
	ClientID     string   `json:"cid"`
	UserID       string   `json:"uid"`
	ClientScopes []string `json:"scopes"`
}

func (u *UserInfo) ClientHasScope(scope string) bool {
	for _, s := range u.ClientScopes {
		if s == scope {
			return true
		}
	}
	return false
}

func logger(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
		log.Info().Str("RequestURI", string(ctx.RequestURI())).
			IPAddr("remoteIP", ctx.RemoteIP()).
			Int("status", ctx.Response.StatusCode()).
			Msg("request")
	}
}

func userInfoMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		unauthorized := func() {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.WriteString(http.StatusText(http.StatusUnauthorized))
		}
		userInfoHeader := ctx.Request.Header.Peek("X-Userinfo")
		if len(userInfoHeader) == 0 {
			unauthorized()
			log.Debug().Msg("no X-Userinfo header")
			return
		}
		userInfo := new(UserInfo)
		if err := json.Unmarshal(userInfoHeader, &userInfo); err != nil {
			unauthorized()
			log.Debug().
				Err(err).
				Bytes("X-Userinfo", userInfoHeader).
				Msg("unauthorized")
			return
		}
		if userInfo.Subject == "" {
			unauthorized()
			log.Debug().Msg("missing subject claim")
			return
		}
		ctx.SetUserValue("userInfo", userInfo)
		h(ctx)
	}
}

func hello(ctx *fasthttp.RequestCtx) {
	userInfo := ctx.UserValue("userInfo").(*UserInfo)
	fmt.Fprintf(ctx, "Hi there, %s!\n", userInfo.Subject)
	fmt.Fprintf(ctx, "user id = %s\n", userInfo.UserID)
	fmt.Fprintf(ctx, "client id = %s\n", userInfo.ClientID)
	fmt.Fprintf(ctx, "scopes = %v\n", userInfo.ClientScopes)
	if userInfo.ClientHasScope("fun") {
		ctx.WriteString("client can have fun")
	} else {
		ctx.WriteString("client can not have fun")
	}
}

func main() {
	log.Info().Int("port", port).Msg("starting server")
	if err := fasthttp.ListenAndServe(
		fmt.Sprintf("localhost:%d", port), logger(userInfoMiddleware(hello))); err != nil {
		log.Fatal().Err(err).Msg("fatal")
	}
}
