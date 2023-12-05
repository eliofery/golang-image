package tpl

import (
	"github.com/eliofery/golang-image/internal/resources"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var (
	assetsDir     = pathJoin("internal/resources/assets")
	assetsPrefix  = "/assets/"
	assetsPattern = assetsPrefix + "*"
)

func AssetsInit(route *chi.Mux) {
	fs := http.FileServer(http.Dir(assetsDir))
	route.Handle(assetsPattern, http.StripPrefix(assetsPrefix, fs))
}

func AssetsFsInit(route *chi.Mux) {
	route.Handle(assetsPattern, http.FileServer(http.FS(resources.Assets)))
}
