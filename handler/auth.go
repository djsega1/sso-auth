package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/djsega1/sso-auth/config"
	"github.com/djsega1/sso-auth/service"
	"github.com/djsega1/sso-auth/utils"
)

type AuthHandler struct {
	UserService *service.UserService
	Config      *config.Config
}

func NewAuthHandler(userService *service.UserService, config *config.Config) *AuthHandler {
	return &AuthHandler{UserService: userService, Config: config}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := h.UserService.RegisterUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Registration failed: %v", err), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"user_id": id.String(),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ok, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(req.Username, h.Config.JWTSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate token: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ValidateJWT(tokenStr, h.Config.JWTSecret)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
		})
		return
	}

	username := claims["username"].(string)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":    true,
		"username": username,
	})
}
