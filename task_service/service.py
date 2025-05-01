from concurrent import futures
import grpc
from sqlalchemy.orm import Session
import task_pb2
import task_pb2_grpc
import database
import logging
from datetime import datetime, timedelta
import os
from typing import Optional
from repos.folderRepo import *
from repos.TaskRepo import *
from google.protobuf.timestamp_pb2 import Timestamp

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def serve():
    db = database.SessionLocal()
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    task_pb2_grpc.add_TaskServiceServicer_to_server(TaskService(db), server)
    server.add_insecure_port('[::]:50052')
    server.start()
    logger.info("Server started on port 50052")
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down server...")
        server.stop(0)
        db.close()

class TaskService(task_pb2_grpc.TaskServiceServicer):
    def __init__(self, database: Session):
        self.db = database

    # ========== Folder Methods ==========
    
    def CreateFolder(self, request, context):
        """Create new folder"""
        try:
            folder = FolderRepo.create_folder(
                self.db,
                FolderCreate(
                    user_id=request.user_id,
                    folder_name=request.name
                )
            )
            return task_pb2.CreateFolderResponse(
                success=True,
                folder=self._folder_to_proto(folder)
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"CreateFolder error: {e}")
            return task_pb2.CreateFolderResponse(success=False)

    def UpdateFolder(self, request, context):
        """Update folder name"""
        try:
            folder = FolderRepo.update_folder(
                self.db,
                FolderUpdate(
                    user_id=request.user_id,
                    folder_id=request.folder_id,
                    new_folder_name=request.new_name
                )
            )
            if not folder:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return task_pb2.UpdateFolderResponse(success=False)
                
            return task_pb2.UpdateFolderResponse(
                success=True,
                folder=self._folder_to_proto(folder)
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"UpdateFolder error: {e}")
            return task_pb2.UpdateFolderResponse(success=False)

    def GetFolder(self, request, context):
        """Get folder by ID"""
        try:
            folder = FolderRepo.get_folder(
                self.db, 
                request.user_id, 
                request.folder_id
            )
            if not folder:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return task_pb2.GetFolderResponse(success=False)
                
            return task_pb2.GetFolderResponse(
                success=True,
                folder=self._folder_to_proto(folder)
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"GetFolder error: {e}")
            return task_pb2.GetFolderResponse(success=False)

    def DeleteFolder(self, request, context):
        """Delete folder and its tasks"""
        try:
            success = FolderRepo.delete_folder(
                self.db,
                FolderDelete(
                    user_id=request.user_id,
                    folder_id=request.folder_id
                )
            )
            return task_pb2.DeleteFolderResponse(
                success=success,
                message="Folder deleted" if success else "Folder not found"
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"DeleteFolder error: {e}")
            return task_pb2.DeleteFolderResponse(success=False)

    # ========== Task Methods ==========

    def CreateTask(self, request, context):
        """Create new task"""
        try:
            task = TaskRepo.create_task(
                self.db,
                TaskCreate(
                    user_id=request.user_id,
                    folder_id=request.folder_id,
                    title=request.title,
                    description=request.description,
                    due_time=self._proto_to_datetime(request.due_time),
                    priority=request.priority
                )
            )
            return task_pb2.CreateTaskResponse(
                success=True,
                task=self._task_to_proto(task)
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"CreateTask error: {e}")
            return task_pb2.CreateTaskResponse(success=False)

    def GetTask(self, request, context):
        """Get task by ID"""
        try:
            task = TaskRepo.get_task(
                self.db, 
                request.user_id, 
                request.task_id
            )
            if not task:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return task_pb2.GetTaskResponse(success=False)
                
            return task_pb2.GetTaskResponse(
                success=True,
                task=self._task_to_proto(task)
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"GetTask error: {e}")
            return task_pb2.GetTaskResponse(success=False)

    def UpdateTask(self, request, context):
        """Update task"""
        try:
            task = TaskRepo.update_task(
                self.db,
                TaskUpdate(
                    task_id=request.task_id,
                    user_id=request.user_id,
                    folder_id=request.folder_id,
                    new_title=request.title,
                    new_description=request.description,
                    due_time=self._proto_to_datetime(request.due_time),
                    priority=request.priority
                )
            )
            if not task:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return task_pb2.UpdateTaskResponse(success=False)
                
            return task_pb2.UpdateTaskResponse(
                success=True,
                task=self._task_to_proto(task)
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"UpdateTask error: {e}")
            return task_pb2.UpdateTaskResponse(success=False)

    def DeleteTask(self, request, context):
        """Delete task"""
        try:
            # Сначала получаем задачу, чтобы узнать folder_id
            task = TaskRepo.get_task(self.db, request.user_id, request.task_id)
            if not task:
                return task_pb2.DeleteTaskResponse(success=False)
                
            success = TaskRepo.delete_task(
                self.db,
                TaskDelete(
                    user_id=request.user_id,
                    task_id=request.task_id,
                    folder_id=task.folder_id
                )
            )
            return task_pb2.DeleteTaskResponse(
                success=success,
                message="Task deleted" if success else "Task not found"
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"DeleteTask error: {e}")
            return task_pb2.DeleteTaskResponse(success=False)

    def SearchTasks(self, request, context):
        """Search tasks with filters"""
        try:
            # Реализация поиска задач
            query = self.db.query(Task).filter(
                Task.user_id == request.user_id
            )
            
            if request.query:
                query = query.filter(
                    (Task.title.ilike(f"%{request.query}%")) |
                    (Task.description.ilike(f"%{request.query}%"))
                )
                
            if request.completed is not None:
                query = query.filter(Task.is_completed == request.completed)
                
            if request.priority:
                query = query.filter(Task.priority == request.priority)
                
            if request.due_before:
                due_date = self._proto_to_datetime(request.due_before)
                query = query.filter(Task.due_time <= due_date)
                
            # Применяем пагинацию
            tasks = query.offset(request.pagination.offset)\
                       .limit(request.pagination.limit)\
                       .all()
                       
            return task_pb2.SearchTasksResponse(
                tasks=[self._task_to_proto(t) for t in tasks],
                total_count=query.count()
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"SearchTasks error: {e}")
            return task_pb2.SearchTasksResponse()

    # ========== Universal methods :D =========

    def GetUserFolders(self, request, context):
        """Get all folders for user"""
        try:
            folders = FolderRepo.get_user_folders(self.db, request.user_id)
            return task_pb2.GetFoldersResponse(
                folders=[self._folder_to_proto(f) for f in folders]
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"GetUserFolders error: {e}")
            return task_pb2.GetFoldersResponse()

    def GetAllTasks(self, request, context):
        """Get all tasks for user (optionally filtered by folder)"""
        try:
            if request.HasField('folder_id'):
                tasks = TaskRepo.get_folder_tasks(
                    self.db, 
                    request.user_id, 
                    request.folder_id
                )
            else:
                # Получаем задачи из всех папок пользователя
                tasks = []
                folders = FolderRepo.get_user_folders(self.db, request.user_id)
                for folder in folders:
                    tasks.extend(
                        TaskRepo.get_folder_tasks(
                            self.db, 
                            request.user_id, 
                            folder.folder_id
                        )
                    )
            
            return task_pb2.GetAllTasksResponse(
                tasks=[self._task_to_proto(t) for t in tasks]
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"GetAllTasks error: {e}")
            return task_pb2.GetAllTasksResponse()

    def ToggleTaskCompletion(self, request, context):
        """Toggle task completion status"""
        try:
            task = TaskRepo.toggle_task_completion(
                self.db, 
                request.user_id, 
                request.task_id
            )
            if not task:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return task_pb2.TaskResponse()
            
            return task_pb2.TaskResponse(
                task=self._task_to_proto(task)
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"ToggleTaskCompletion error: {e}")
            return task_pb2.TaskResponse()

    def MoveTaskToFolder(self, request, context):
        """Move task to another folder"""
        try:
            # Проверяем существование задачи
            task = TaskRepo.get_task(
                self.db, 
                request.user_id, 
                request.task_id
            )
            if not task:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return task_pb2.TaskResponse()
            
            # Проверяем существование целевой папки
            folder = FolderRepo.get_folder(
                self.db, 
                request.user_id, 
                request.new_folder_id
            )
            if not folder:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                return task_pb2.TaskResponse()
            
            # Обновляем папку задачи
            updated_task = TaskRepo.update_task(
                self.db,
                TaskUpdate(
                    task_id=request.task_id,
                    user_id=request.user_id,
                    folder_id=request.new_folder_id,
                    new_title=task.title,
                    new_description=task.description,
                    due_time=task.due_time,
                    priority=task.priority
                )
            )
            
            return task_pb2.TaskResponse(
                task=self._task_to_proto(updated_task)
            )
        except Exception as e:
            self.db.rollback()
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"MoveTaskToFolder error: {e}")
            return task_pb2.TaskResponse()

    def SearchTasks(self, request, context):
        """Search tasks with filters"""
        try:
            # Базовый запрос
            query = self.db.query(Task).filter(
                Task.user_id == request.user_id
            )
            
            # Применяем фильтры
            if request.query:
                query = query.filter(
                    (Task.title.ilike(f"%{request.query}%")) |
                    (Task.description.ilike(f"%{request.query}%"))
                )
            
            if request.completed is not None:
                query = query.filter(
                    Task.is_completed == request.completed
                )
            
            if request.priority:
                query = query.filter(
                    Task.priority == request.priority
                )
            
            if request.HasField('due_before'):
                due_date = self._proto_to_datetime(request.due_before)
                query = query.filter(
                    Task.due_time <= due_date
                )
            
            # Применяем пагинацию
            tasks = query.order_by(Task.due_time.asc())\
                       .offset(request.pagination.offset)\
                       .limit(request.pagination.limit)\
                       .all()
            
            return task_pb2.SearchTasksResponse(
                tasks=[self._task_to_proto(t) for t in tasks],
                total_count=query.count()
            )
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            logger.error(f"SearchTasks error: {e}")
            return task_pb2.SearchTasksResponse()

    # ========== Utility Methods ==========
    
    def _folder_to_proto(self, folder: TaskFolder) -> task_pb2.Folder:
        """Convert SQLAlchemy Folder to protobuf message"""
        return task_pb2.Folder(
            folder_id=folder.folder_id,
            user_id=folder.user_id,
            name=folder.name,
            created_at=self._datetime_to_proto(folder.created_at),
        )
    
    def _task_to_proto(self, task: Task) -> task_pb2.Task:
        """Convert SQLAlchemy Task to protobuf message"""
        return task_pb2.Task(
            task_id=task.task_id,
            folder_id=task.folder_id,
            user_id=task.user_id,
            title=task.title,
            description=task.description,
            due_time=self._datetime_to_proto(task.due_time),
            priority=task.priority,
            is_completed=task.is_completed,
            created_at=self._datetime_to_proto(task.created_at),
            updated_at=self._datetime_to_proto(task.updated_at)
        )
    
    def _datetime_to_proto(self, dt: datetime) -> Optional[Timestamp]:
        """Convert datetime to protobuf Timestamp"""
        if not dt:
            return None
        timestamp = Timestamp()
        timestamp.FromDatetime(dt)
        return timestamp

    def _proto_to_datetime(self, ts: Timestamp) -> Optional[datetime]:
        """Convert protobuf Timestamp to datetime"""
        if not ts:
            return None
        return ts.ToDatetime()

    

            



