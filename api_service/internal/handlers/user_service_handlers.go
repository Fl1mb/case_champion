package handlers

import (
	"api_service/internal/grpc/grpc_server"
	"api_service/internal/grpc/server/user_grpc"
	grpccache "api_service/internal/grpc_cache"
	grpcclient "api_service/internal/grpc_client"
	"api_service/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserAuthHandler struct {
	Client *grpcclient.UserServiceClient
	cache  *grpccache.CacheClient
}

func (h *UserAuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

func (h *UserAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user models.UserModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Username, email and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.CreateUser(r.Context(), &user_grpc.CreateUserRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		FullName: user.Fullname,
	})

	if err != nil {
		h.handleGRPCError(w, err)
		return
	}

	if resp.Error != "" {
		http.Error(w, resp.Error, http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   resp.Id,
		"username":  resp.Username,
		"email":     resp.Email,
		"full_name": resp.FullName,
	})
}

func (h *UserAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user models.UserModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	loginResp, err := h.Client.Login(r.Context(), &user_grpc.LoginRequest{
		Username: user.Username,
		Password: user.Password,
	})

	if err != nil {
		h.handleGRPCError(w, err)
		return
	}

	resp, err := h.cache.Write(r.Context(), &grpc_server.WriteRequest{
		UserId:    loginResp.UserId,
		UserLogin: user.Username,
		JwtKey:    loginResp.AccessToken,
	})

	if err != nil || !resp.Success {
		http.Error(w, "Error to write to cache", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": loginResp.AccessToken,
	})
}

func (h *UserAuthHandler) handleGRPCError(w http.ResponseWriter, err error) {
	// Здесь можно добавить обработку различных gRPC ошибок
	// и преобразование их в соответствующие HTTP статусы
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func (h *UserAuthHandler) Close() error {
	return h.Client.Close()
}

func NewUserAuthHandler(addr string, chc *grpccache.CacheClient) (*UserAuthHandler, error) {
	user_service_client, err := grpcclient.NewUserServiceClient("user_service:50051")
	if err != nil {
		return nil, err
	}
	return &UserAuthHandler{
		Client: user_service_client,
		cache:  chc,
	}, nil
}
