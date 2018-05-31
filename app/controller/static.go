package controller

import (
	"net/http"
)

// Static maps static files
func Static(w http.ResponseWriter, r *http.Request) {
	// Disable listing directories
	// if strings.HasSuffix(r.URL.Path, "/") {
	// 	Error404(w, r)
	// 	return
	// }
	http.ServeFile(w, r, r.URL.Path[1:])
}
