package main

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type responseRecorder struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func newResponseRecorder() *responseRecorder {
	return &responseRecorder{
		header: make(http.Header),
		status: http.StatusOK,
	}
}

func (r *responseRecorder) Header() http.Header {
	return r.header
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *responseRecorder) writeTo(w http.ResponseWriter) {
	for key, values := range r.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	if r.status != 0 {
		w.WriteHeader(r.status)
	}
	_, _ = io.Copy(w, &r.body)
}

func newAssetHandler(assets fs.FS) http.Handler {
	inner := application.AssetFileServerFS(assets)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldTrySPAFallback(r) {
			inner.ServeHTTP(w, r)
			return
		}

		rec := newResponseRecorder()
		inner.ServeHTTP(rec, r)
		if rec.status == http.StatusNotFound {
			fallback := r.Clone(r.Context())
			fallback.URL.Path = "/"
			fallback.RequestURI = "/"
			inner.ServeHTTP(w, fallback)
			return
		}

		rec.writeTo(w)
	})
}

func shouldTrySPAFallback(r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}

	requestPath := r.URL.Path
	if requestPath == "/" || requestPath == "/index.html" {
		return false
	}
	if strings.HasPrefix(requestPath, "/wails") {
		return false
	}

	ext := path.Ext(requestPath)
	return ext == "" || ext == ".html"
}
