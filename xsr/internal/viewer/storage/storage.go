package storage

import (
	"fmt"

	// external
	"github.com/valyala/fasthttp"

	// project
	. "github.com/ultary/oss/xsr/internal/viewer/config"
)

var empty = make([]*string, 0)

type Storage interface {
	Find(host, path string) string
	Serve(ctx *fasthttp.RequestCtx, key string)
}

func NewStorage(cfg *Config) Storage {
	if cfg.StorageType == "" {
		const m = "storage type is empty"
		panic(m)
	}

	switch cfg.StorageType {
	case StorageTypeGCS:
		return NewGCS(cfg)
	case StorageTypeS3:
		return NewS3(cfg)
	default:
		const f = "storage type (%s) is invalid"
		m := fmt.Sprintf(f, cfg.StorageType)
		panic(m)
	}
}
