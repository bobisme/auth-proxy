package main

import (
	"strings"
	"github.com/valyala/fasthttp"
	"github.com/yeqown/fasthttp-reverse-proxy"
	"github.com/rs/zerolog/log"
)

var (
	proxyServer = proxy.NewReverseProxy("localhost:8081")
)

func userInfoForAuthorization(header string) []byte {
	if header == "" {
		return nil
	}
	if !strings.HasPrefix(header, "Bearer ") {
		log.Debug().Str("header", header).Msg("invalid authorization")
		return nil
	}
	// token := header[8:len(header)-1]
	// DEBUG
	return []byte(`{"uid":"123","sub":"bob@8o8.me"}`)
}

func ProxyHandler(ctx *fasthttp.RequestCtx) {
	ctx.Request.Header.Del("X-Userinfo")
	authorization := string(ctx.Request.Header.Peek("Authorization"))
	if userInfo := userInfoForAuthorization(authorization); userInfo != nil {
		ctx.Request.Header.AddBytesV("X-Userinfo", userInfo)
	}
	proxyServer.ServeHTTP(ctx)
}

type logger struct {}

func (l logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func main() {
	log.Info().Int("port", 8080).Msg("starting proxy")
	s := fasthttp.Server{
		Handler: ProxyHandler,
		Logger: logger{},
	}
	if err := s.ListenAndServe(":8080"); err != nil {
		log.Fatal().Err(err).Msg("fatal")
	}
}
