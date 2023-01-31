package controller

import (
	"net/http"
	"psyWeb/web/views"
)

func HandleStaticResource(w http.ResponseWriter, r *http.Request) {
	views.RenderStaticResources(w, r.URL.Path)
}
