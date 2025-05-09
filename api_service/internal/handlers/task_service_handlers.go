package handlers

import (
	"api_service/internal/grpc/grpc_server"
	grpccache "api_service/internal/grpc_cache"
	task_server "api_service/internal/grpc_task"
	taskclient "api_service/internal/grpc_task/task_client"
	"api_service/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServiceHandler struct {
	Client *taskclient.TaskServiceClient
	Cache  *grpccache.CacheClient
}

func (h *TaskServiceHandler) Close() error {
	return h.Client.Close()
}

func NewTaskServiceClient(addr_task string, cache *grpccache.CacheClient) (*TaskServiceHandler, error) {
	task_service, err := taskclient.NewTaskServiceClient("task_service:50052")
	if err != nil {
		return nil, err
	}
	return &TaskServiceHandler{
		Client: task_service,
		Cache:  cache,
	}, nil
}

func (h *TaskServiceHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			http.Error(w, "Error authorization", http.StatusBadRequest)
			return
		}
		response, err := h.Cache.GetUser(r.Context(), &grpc_server.GetUserRequest{
			JwtKey: authToken,
		})
		if err != nil {
			http.Error(w, "Error in server cache", http.StatusConflict)
			return
		}
		if !response.Success {
			http.Error(w, "Invalid jwt key", http.StatusConflict)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", response.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RegisterRoutes регистрирует все маршруты для задач и папок
func (h *TaskServiceHandler) RegisterRoutes(r chi.Router) {
	r.Use(h.AuthMiddleware)

	r.Route("/folders", func(r chi.Router) {
		r.Get("/", h.GetUserFolders)
		r.Post("/", h.CreateFolder)
		r.Route("/{folderID}", func(r chi.Router) {
			r.Get("/", h.GetFolder)
			r.Put("/", h.UpdateFolder)
			r.Delete("/", h.DeleteFolder)
			r.Get("/tasks", h.GetFolderTasks)
		})
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", h.GetAllTasks)
		r.Post("/", h.CreateTask)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Get("/", h.GetTask)
			r.Put("/", h.UpdateTask)
			r.Delete("/", h.DeleteTask)
			r.Put("/toggle", h.ToggleTaskCompletion)
			r.Put("/move", h.MoveTaskToFolder)
		})
		r.Get("/search", h.SearchTasks)
	})
}

// Folder handlers
func (h *TaskServiceHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user_id").(int32)
	var folder models.FolderModel
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, "invalid json object", http.StatusBadRequest)
		return
	}
	resp, err := h.Client.CreateFolder(r.Context(), &task_server.CreateFolderRequest{
		UserId: user_id,
		Name:   folder.Name,
	})
	if err != nil {
		http.Error(w, "Error to create folder", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"user_id":     strconv.Itoa(int(user_id)),
		"folder_id":   strconv.Itoa(int(resp.Folder.FolderId)),
		"folder_name": resp.Folder.Name,
		"created_at":  resp.Folder.CreatedAt.String(),
	})
}

func (h *TaskServiceHandler) GetUserFolders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)

	resp, err := h.Client.GetUserFolders(r.Context(), &task_server.GetFoldersRequest{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, "Error getting folders", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Folders)
}

func (h *TaskServiceHandler) GetFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	folderID, err := strconv.ParseInt(chi.URLParam(r, "folderID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.GetFolder(r.Context(), &task_server.GetFolderRequest{
		UserId:   userID,
		FolderId: int32(folderID),
	})
	if err != nil {
		http.Error(w, "Error getting folder", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Folder)
}

func (h *TaskServiceHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	folderID, err := strconv.ParseInt(chi.URLParam(r, "folderID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	var folder models.FolderModel
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, "invalid json object", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.UpdateFolder(r.Context(), &task_server.UpdateFolderRequest{
		FolderId: int32(folderID),
		UserId:   userID,
		NewName:  folder.Name,
	})
	if err != nil {
		http.Error(w, "Error updating folder", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Folder)
}

func (h *TaskServiceHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	folderID, err := strconv.ParseInt(chi.URLParam(r, "folderID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.DeleteFolder(r.Context(), &task_server.DeleteFolderRequest{
		FolderId: int32(folderID),
		UserId:   userID,
	})
	if err != nil {
		http.Error(w, "Error deleting folder", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": resp.Success,
		"message": resp.Message,
	})
}

func (h *TaskServiceHandler) GetFolderTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	folderID, err := strconv.ParseInt(chi.URLParam(r, "folderID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	folderid := int32(folderID)
	resp, err := h.Client.GetAllTasks(r.Context(), &task_server.GetAllTasksRequest{
		UserId:   userID,
		FolderId: &folderid,
	})
	if err != nil {
		http.Error(w, "Error getting folder tasks", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": resp.Tasks,
	})
}

// Task handlers

func (h *TaskServiceHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	folderID := r.URL.Query().Get("folder_id")

	req := &task_server.GetAllTasksRequest{
		UserId: userID,
	}
	if folderID != "" {
		fID, err := strconv.ParseInt(folderID, 10, 32)
		if err != nil {
			http.Error(w, "Invalid folder_id", http.StatusBadRequest)
			return
		}
		folderid := int32(fID)
		req.FolderId = &folderid
	}

	resp, err := h.Client.GetAllTasks(r.Context(), req)
	if err != nil {
		http.Error(w, "Error getting tasks", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": resp.Tasks,
	})
}

func (h *TaskServiceHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)

	var task models.TaskModel
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid json object", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.CreateTask(r.Context(), &task_server.CreateTaskRequest{
		UserId:      userID,
		FolderId:    task.FolderID,
		Title:       task.Title,
		Description: task.Description,
		DueTime:     timestamppb.New(task.Due_time),
		Priority:    task.Priority,
	})
	if err != nil {
		http.Error(w, "Error creating task", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Task)
}

func (h *TaskServiceHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.GetTask(r.Context(), &task_server.GetTaskRequest{
		UserId: userID,
		TaskId: int32(taskID),
	})
	if err != nil {
		http.Error(w, "Error getting task", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Task)
}

func (h *TaskServiceHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.TaskModel
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid json object", http.StatusBadRequest)
		return
	}

	title := task.Title
	desc := task.Description
	prt := task.Priority
	folder_id := task.FolderID
	req := &task_server.UpdateTaskRequest{
		TaskId:      int32(taskID),
		UserId:      userID,
		FolderId:    &folder_id,
		Title:       &title,
		Description: &desc,
		Priority:    &prt,
	}
	if !task.Due_time.IsZero() {
		req.DueTime = timestamppb.New(task.Due_time)
	}

	resp, err := h.Client.UpdateTask(r.Context(), req)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Task)
}

func (h *TaskServiceHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.DeleteTask(r.Context(), &task_server.DeleteTaskRequest{
		TaskId: int32(taskID),
		UserId: userID,
	})
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": resp.Success,
		"message": resp.Message,
	})
}

func (h *TaskServiceHandler) ToggleTaskCompletion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.ToggleTaskCompletion(r.Context(), &task_server.ToggleTaskRequest{
		TaskId: int32(taskID),
		UserId: userID,
	})
	if err != nil {
		http.Error(w, "Error toggling task completion", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Task)
}

func (h *TaskServiceHandler) MoveTaskToFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var moveReq struct {
		NewFolderID int32 `json:"new_folder_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&moveReq); err != nil {
		http.Error(w, "invalid json object", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.MoveTaskToFolder(r.Context(), &task_server.MoveTaskRequest{
		TaskId:      int32(taskID),
		UserId:      userID,
		NewFolderId: moveReq.NewFolderID,
	})
	if err != nil {
		http.Error(w, "Error moving task", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Task)
}

func (h *TaskServiceHandler) SearchTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int32)
	query := r.URL.Query().Get("query")
	completed := r.URL.Query().Get("completed")
	priority := r.URL.Query().Get("priority")
	dueBefore := r.URL.Query().Get("due_before")

	req := &task_server.SearchTasksRequest{
		UserId: userID,
		Query:  query,
	}

	if completed != "" {
		comp, err := strconv.ParseBool(completed)
		if err == nil {
			req.Completed = &comp
		}
	}

	if priority != "" {
		prio, err := strconv.ParseInt(priority, 10, 32)
		if err == nil {
			prt := int32(prio)
			req.Priority = &prt
		}
	}

	if dueBefore != "" {
		dueTime, err := time.Parse(time.RFC3339, dueBefore)
		if err == nil {
			req.DueBefore = timestamppb.New(dueTime)
		}
	}

	resp, err := h.Client.SearchTasks(r.Context(), req)
	if err != nil {
		http.Error(w, "Error searching tasks", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks":       resp.Tasks,
		"total_count": resp.TotalCount,
	})
}
