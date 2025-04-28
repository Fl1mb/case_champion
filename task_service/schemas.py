from pydantic import BaseModel, Field, conint
from datetime import datetime
from typing import Optional, List

# Общие настройки для полей
USER_ID_DESC = "ID пользователя из сервиса авторизации"
FOLDER_NAME_DESC = "Название папки (1-50 символов)"

class TaskBase(BaseModel):
    title: str = Field(..., min_length=1, max_length=100, description="Название задачи")
    description: str = Field(..., min_length=1, max_length=500, description="Описание задачи")
    due_time: datetime = Field(..., description="Срок выполнения")
    priority: conint(ge=1, le=5) = Field(..., description="Приоритет (1-5)")  # Ограничение 1-5

class TaskCreate(TaskBase):
    user_id: int = Field(..., description=USER_ID_DESC)
    folder_id: int = Field(..., description="ID папки")

class TaskUpdate(TaskBase):
    user_id: int = Field(..., description=USER_ID_DESC)
    folder_id: int = Field(..., description="ID папки")
    task_id: int = Field(..., description="ID задачи")

class TaskDelete(BaseModel):
    user_id: int = Field(..., description=USER_ID_DESC)
    folder_id: int = Field(..., description="ID папки")
    task_id: int = Field(..., description="ID задачи")

class TaskResponse(TaskBase):
    task_id: int
    user_id: int
    folder_id: int
    created_at: datetime
    updated_at: datetime
    is_completed: bool = Field(False, description="Статус выполнения")

    class Config:
        json_encoders = {
            datetime: lambda v: v.isoformat()  # Правильное форматирование datetime
        }

class FolderBase(BaseModel):
    folder_name: str = Field(
        ..., 
        min_length=1, 
        max_length=50,
        description=FOLDER_NAME_DESC
    )

class FolderCreate(FolderBase):
    user_id: int = Field(..., description=USER_ID_DESC)
    folder_name : str = Field(..., min_length=1, max_length=50)

class FolderUpdate(FolderBase):
    user_id: int = Field(..., description=USER_ID_DESC)
    folder_id: int = Field(..., description="ID папки")
    new_folder_name: str = Field(..., min_length=1, max_length=50, description=FOLDER_NAME_DESC)

class FolderDelete(BaseModel):
    user_id: int = Field(..., description=USER_ID_DESC)
    folder_id: int = Field(..., description="ID папки")

class FolderResponse(FolderBase):
    folder_id: int
    user_id: int
    created_at: datetime
    task_ids: List[int] = Field(default_factory=list, description="Список ID задач в папке")

    class Config:
        json_encoders = {
            datetime: lambda v: v.isoformat()
        }