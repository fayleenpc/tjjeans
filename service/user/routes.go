package user

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/fayleenpc/ecom/config"
	"github.com/fayleenpc/ecom/service/auth"
	"github.com/fayleenpc/ecom/types"
	"github.com/fayleenpc/ecom/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
	sync.WaitGroup
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.WithLogger(auth.WithRateLimiter(h.handleLogin))).Methods("POST")
	router.HandleFunc("/register", auth.WithLogger(auth.WithRateLimiter(h.handleRegister))).Methods("POST")
	router.HandleFunc("/verify", auth.WithLogger(auth.WithRateLimiter(h.handleVerify))).Methods("GET")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate the payload

	if err := utils.Validate.Struct(payload); err != nil {
		errv := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errv))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	// fmt.Fprintf(w, "Registration successful. Please check your email to verify.")

	// if u.Email == "jackwdy1@gmail.com" {
	// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token, "admin": "true"})
	// } else {
	// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	// }

}

func (h *Handler) handleVerify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email, exists := auth.EmailVerificationTokens[token]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid or expired token"))
		return
	}

	delete(auth.EmailVerificationTokens, token)

	// Here you can update the user's status in your database or mark the email as verified.
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Email %s has been verified successfully!\n", email))
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {

	// get json payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", error))
		return
	}

	// check if the user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// get token verification
	tokenVerification, err := auth.GenerateToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error generating token verification"))
		return
	}
	// verification here
	auth.EmailVerificationTokens[tokenVerification] = payload.Email
	// if err := auth.SendVerificationEmail(payload.Email, tokenVerification); err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error sending verification email"))
	// 	return
	// }

	// if it doesn't create new user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
