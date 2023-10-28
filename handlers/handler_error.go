package handlers

import "net/http"

func HandlerError(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, 400, "Something went wrong")
}
