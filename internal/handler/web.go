package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) WebHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: load this from environment variables
	buildDir := "data/public"
	fs := http.FileServer(http.Dir(buildDir))
	if _, err := os.Stat(filepath.Join(buildDir, r.URL.Path)); err == nil {
		fs.ServeHTTP(w, r)
		return
	}
	// otherwise serve index.html
	http.ServeFile(w, r, filepath.Join(buildDir, "index.html"))
}
