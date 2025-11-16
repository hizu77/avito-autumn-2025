package health

import (
	"net/http"

	"github.com/go-chi/render"
)

func Liveness(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]any{"status": "ok"})
}
