package handler

import (
	"net/http"
	backend "restgeld-backend"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	backend.Handler(w, r)
}
