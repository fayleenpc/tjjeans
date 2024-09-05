package products

import (
	"net/http"

	"github.com/fayleenpc/tjjeans/service/auth"
	"github.com/fayleenpc/tjjeans/types"
	"github.com/fayleenpc/tjjeans/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/product", auth.WithLogger(auth.WithRateLimiter(h.handleGetProducts))).Methods("GET")
	// router.HandleFunc("/product", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods("POST")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, ps)
}

// func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {

// 	var payload types.Product
// 	if err := utils.ParseJSON(r, &payload); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	id, err := h.store.CreateProduct(payload)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	payload.ID = id
// 	utils.WriteJSON(w, http.StatusOK, payload)
// }
