package main

import (
	"api_service/internal/grpc/server/user_grpc"
	grpccache "api_service/internal/grpc_cache"
	grpcclient "api_service/internal/grpc_client"
	task_server "api_service/internal/grpc_task"
	taskclient "api_service/internal/grpc_task/task_client"
	"api_service/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	time.Sleep(time.Second * 3)
	user_service_client, err := grpcclient.NewUserServiceClient("user_service:50051")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to user_service:50051")
	defer user_service_client.Close()

	task_service, err := taskclient.NewTaskServiceClient("task_service:50052")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to task_service:50052")
	defer task_service.Close()

	cache_service, err := grpccache.NewCacheClient("cache_service:50053")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to cache_service:50053")
	defer cache_service.Close()

	http.HandleFunc("/task/CreateFolder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid http method", http.StatusBadRequest)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Invalid format of data", http.StatusBadRequest)
			return
		}
		var folder models.FolderModel
		err := json.NewDecoder(r.Body).Decode(&folder)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := task_service.CreateFolder(r.Context(), &task_server.CreateFolderRequest{
			UserId: folder.UserID,
			Name:   folder.Name,
		})

		if err != nil {
			http.Error(w, "Error to req task service", http.StatusBadGateway)
			return
		}

		if !resp.Success {
			http.Error(w, "Error to create task", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"folder_id":  resp.Folder.FolderId,
			"user_id":    resp.Folder.UserId,
			"name":       resp.Folder.Name,
			"created_at": resp.Folder.CreatedAt.AsTime(),
		})

	})

	http.HandleFunc("/task/CreateTask", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid http method", http.StatusBadRequest)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Invalid format of data", http.StatusBadRequest)
			return
		}
		var task models.TaskModel
		err := json.NewDecoder(r.Body).Decode(&task)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := task_service.CreateTask(r.Context(), &task_server.CreateTaskRequest{
			UserId:      task.UserID,
			FolderId:    task.FolderID,
			Title:       task.Title,
			Description: task.Description,
			DueTime:     timestamppb.New(task.Due_time),
			Priority:    task.Priority,
		})

		if err != nil {
			http.Error(w, "Error to req task service", http.StatusBadGateway)
			return
		}

		if !resp.Success {
			http.Error(w, "Error to create task", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"task_id":      resp.Task.TaskId,
			"folder_id":    resp.Task.FolderId,
			"user_id":      resp.Task.UserId,
			"title":        resp.Task.Title,
			"description":  resp.Task.Description,
			"due_time":     resp.Task.DueTime.AsTime(),
			"priority":     resp.Task.Priority,
			"is_completed": resp.Task.IsCompleted,
			"created_at":   resp.Task.CreatedAt.AsTime(),
			"updated_at":   resp.Task.UpdatedAt.AsTime(),
		})
	})

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
