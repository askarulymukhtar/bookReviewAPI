package handlers

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, 200, struct{}{})
}
