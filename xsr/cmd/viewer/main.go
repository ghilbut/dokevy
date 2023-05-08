package main

import (
	// external
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	// project
	. "github.com/ultary/oss/xsr/internal/viewer/config"
	. "github.com/ultary/oss/xsr/internal/viewer/storage"
)

func main() {
	cfg := NewConfig()
	log.SetLevel(cfg.LogLevel)

	storage := NewStorage(cfg)

	log.Info("RUN")
	if err := fasthttp.ListenAndServe(cfg.Address, handle(storage, cfg.Bucket)); err != nil {
		log.Fatal(err)
	}
}

// handle returns a fasthttp request handler
func handle(storage Storage, bucket string) fasthttp.RequestHandler {
	return fasthttp.CompressHandlerLevel(
		func(ctx *fasthttp.RequestCtx) {
			host := string(ctx.Host())
			path := string(ctx.Path())

			defer func() {
				if r := recover(); r != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					ctx.SetContentType("text/html; charset=utf-8")
					ctx.SetBody([]byte("Internal Server Error"))
				}
			}()

			key := storage.Find(host, path)
			if key == "" {
				ctx.NotFound()
				return
			}

			storage.Serve(ctx, key)
		},
		fasthttp.CompressBestCompression,
	)
}
