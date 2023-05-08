package storage

import (
	// external
	"github.com/valyala/fasthttp"

	// project
	. "github.com/ultary/oss/xsr/internal/viewer/config"
)

type GCS struct {
}

func NewGCS(cfg *Config) *GCS {
	return &GCS{}
}

// Find is not implemented yet
func (s *GCS) Find(host, path string) string {
	panic("gcs find function is not implemented yet")
}

// Serve is not implemented yet
func (s *GCS) Serve(ctx *fasthttp.RequestCtx, key string) {
	panic("gcs serve function is not implemented yet")
}
