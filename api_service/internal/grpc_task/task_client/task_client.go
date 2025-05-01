package taskclient

import (
	task_server "api_service/internal/grpc_task"
	"context"
	"time"

	"google.golang.org/grpc"
)

type TaskServiceClient struct {
	task_server.UnimplementedTaskServiceServer
	Client task_server.TaskServiceClient
	conn   *grpc.ClientConn
}

// Realization of interface
func (t *TaskServiceClient) CreateFolder(ctx context.Context, req *task_server.CreateFolderRequest) (
	*task_server.CreateFolderResponse,
	error,
) {
	return t.Client.CreateFolder(ctx, req)
}

func (t *TaskServiceClient) CreateTask(ctx context.Context, req *task_server.CreateTaskRequest) (
	*task_server.CreateTaskResponse,
	error,
) {
	return t.Client.CreateTask(ctx, req)
}

func (t *TaskServiceClient) DeleteFolder(ctx context.Context, req *task_server.DeleteFolderRequest) (
	*task_server.DeleteFolderResponse,
	error,
) {
	return t.Client.DeleteFolder(ctx, req)
}

func (t *TaskServiceClient) DeleteTask(ctx context.Context, req *task_server.DeleteTaskRequest) (
	*task_server.DeleteTaskResponse,
	error,
) {
	return t.Client.DeleteTask(ctx, req)
}

func (t *TaskServiceClient) GetAllTasks(ctx context.Context, req *task_server.GetAllTasksRequest) (
	*task_server.GetAllTasksResponse,
	error,
) {
	return t.Client.GetAllTasks(ctx, req)
}

func (t *TaskServiceClient) GetFolder(ctx context.Context, req *task_server.GetFolderRequest) (
	*task_server.GetFolderResponse,
	error,
) {
	return t.Client.GetFolder(ctx, req)
}

func (t *TaskServiceClient) GetTask(ctx context.Context, req *task_server.GetTaskRequest) (
	*task_server.GetTaskResponse,
	error,
) {
	return t.Client.GetTask(ctx, req)
}

func (t *TaskServiceClient) GetUserFolders(ctx context.Context, req *task_server.GetFoldersRequest) (
	*task_server.GetFoldersResponse,
	error,
) {
	return t.Client.GetUserFolders(ctx, req)
}

func (t *TaskServiceClient) MoveTaskToFolder(ctx context.Context, req *task_server.MoveTaskRequest) (
	*task_server.TaskResponse,
	error,
) {
	return t.Client.MoveTaskToFolder(ctx, req)
}

func (t *TaskServiceClient) SearchTasks(ctx context.Context, req *task_server.SearchTasksRequest) (
	*task_server.SearchTasksResponse,
	error,
) {
	return t.Client.SearchTasks(ctx, req)
}

func (t *TaskServiceClient) ToggleTaskCompletion(ctx context.Context, req *task_server.ToggleTaskRequest) (
	*task_server.TaskResponse,
	error,
) {
	return t.Client.ToggleTaskCompletion(ctx, req)
}

func (t *TaskServiceClient) UpdateFolder(ctx context.Context, req *task_server.UpdateFolderRequest) (
	*task_server.UpdateFolderResponse,
	error,
) {
	return t.Client.UpdateFolder(ctx, req)
}

func (t *TaskServiceClient) UpdateTask(ctx context.Context, req *task_server.UpdateTaskRequest) (
	*task_server.UpdateTaskResponse,
	error,
) {
	return t.Client.UpdateTask(ctx, req)
}

func NewTaskServiceClient(addr string) (*TaskServiceClient, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		return nil, err
	}
	return &TaskServiceClient{
		Client: task_server.NewTaskServiceClient(conn),
		conn:   conn,
	}, nil
}

func (t *TaskServiceClient) Close() error {
	return t.conn.Close()
}
