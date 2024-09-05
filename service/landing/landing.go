package landing

import (
	"net/http"

	"github.com/fayleenpc/tjjeans/views"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		views.Index().Render(r.Context(), w)
	})
	// router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
	// 	// authentication for only 1 user in db , role is unique id and password as admin itself

	// })
}
