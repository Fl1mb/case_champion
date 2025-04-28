package main

import (
	"api_service/internal/grpc/server/user_grpc"
	grpcclient "api_service/internal/grpc_client"
	"api_service/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	time.Sleep(time.Second * 3)
	user_service_client, err := grpcclient.NewUserServiceClient("user_service:50051")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to user_service:50051")
	defer user_service_client.Close()

	http.HandleFunc("/auth/Register", func(w http.ResponseWriter, r *http.Request) {
		//Read json file and send it to user_service
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid http method", http.StatusBadRequest)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Invalid format of data", http.StatusBadRequest)
			return
		}
		var user models.UserModel
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := user_service_client.CreateUser(r.Context(), &user_grpc.CreateUserRequest{
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			FullName: user.Fullname,
		})

		if err != nil {
			http.Error(w, "Error to request Authorization service", http.StatusExpectationFailed)
			return
		}

		if response.Error != "" {
			http.Error(w, response.Error, http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":   response.Id,
			"username":  response.Username,
			"email":     response.Email,
			"full_name": response.FullName,
		})

	})
	http.HandleFunc("/auth/Login", func(w http.ResponseWriter, r *http.Request) {
		//Read json file and send it to user_service
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid http method", http.StatusBadRequest)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Invalid format of data", http.StatusBadRequest)
			return
		}
		var user models.UserModel
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := user_service_client.Login(r.Context(), &user_grpc.LoginRequest{
			Username: user.Username,
			Password: user.Password,
		})

		if err != nil {
			http.Error(w, "Error to request Authorization service", http.StatusExpectationFailed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": response.AccessToken,
		})
	})

	log.Println("Server listening on port 3723:3723")
	if err := http.ListenAndServe(":3723", nil); err != nil {
		log.Fatal(nil)
	}
}
