package vitgo

import (
	"log"
	"net/http"
	"strings"
)

// Redirector for dev server
func (vg *VitGo) DevServerRedirector() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		original := r.URL.Path
		prefix := "/dev/"

		if len(original) < len(prefix) || original[:len(prefix)] != prefix {
			http.NotFound(w, r)
			return
		}

		rest := original[len(prefix)-1:]

		escapedRest := strings.Replace(rest, "\n", "", -1)
		escapedRest = strings.Replace(escapedRest, "\r", "", -1)

		log.Println("rest: ", escapedRest)

		w.Header().Set("Content-Type", "application/javascript")
		http.Redirect(w, r, vg.DevServer+rest, http.StatusPermanentRedirect)
	}

	return http.HandlerFunc(handler)
}
