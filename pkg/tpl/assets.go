package tpl

import (
	"github.com/eliofery/golang-image/internal/resources"
	"github.com/go-chi/chi/v5"
	"net/http"
	"path/filepath"
)

var (
	assetsDir     = filepath.Join("internal/resources/assets")
	assetsPrefix  = "/assets/"
	assetsPattern = assetsPrefix + "*"
)

func AssetsInit(route *chi.Mux) {
	fs := http.FileServer(http.Dir(assetsDir))
	route.Handle(assetsPattern, http.StripPrefix(assetsPrefix, fs))
}

func AssetsFsInit(route *chi.Mux) {
	fs := http.FileServer(http.FS(resources.Assets))
	route.Handle(assetsPattern, fs)
}
