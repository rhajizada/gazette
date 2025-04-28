package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	buildDir := "./dist"
	fs := http.FileServer(http.Dir(buildDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if file exists, serve it
		if _, err := os.Stat(filepath.Join(buildDir, r.URL.Path)); err == nil {
			fs.ServeHTTP(w, r)
			return
		}
		// otherwise serve index.html
		http.ServeFile(w, r, filepath.Join(buildDir, "index.html"))
	})

	http.ListenAndServe(":9191", nil)
}
