from sqlalchemy.orm import Session
from sqlalchemy import desc, func
from models import *
from schemas import *

class TaskRepo:
    @staticmethod
    def create_task(db: Session, task_data: TaskCreate)-> Task:
        """Creating task in db"""
        db_task = Task(
            user_id=task_data.user_id,
            folder_id=task_data.folder_id,
            title=task_data.title,
            description=task_data.description,
            due_time=task_data.due_time,
            priority=task_data.priority
        )
        db.add(db_task)
        db.commit()
        db.refresh(db_task)
        return db_task
    
    @staticmethod
    def get_task(db: Session, user_id:int, task_id: int)->Optional[Task]:
        """Get user task by id"""
        return db.query(Task).filter(
            Task.task_id == task_id,
            Task.user_id == user_id
        ).first()

    @staticmethod
    def get_folder_tasks(db: Session, user_id:int, folder_id: int)->List[Task]:
        """Get all task from folder"""
        return db.query(Task).filter(
            Task.user_id == user_id,
            Task.folder_id == folder_id
        ).order_by(Task.priority.desc(), Task.due_time).all()
    
    @staticmethod
    def update_task(db: Session, task_data: TaskUpdate) -> Optional[Task]:
        """Update task"""
        task = db.query(Task).filter(
            Task.task_id == task_data.task_id,
            Task.user_id == task_data.user_id
        ).first()
        
        if task:
            task.title = task_data.title
            task.description = task_data.description
            task.due_time = task_data.due_time
            task.priority = task_data.priority
            task.folder_id = task_data.folder_id
            db.commit()
            db.refresh(task)
        return task

    @staticmethod
    def delete_task(db: Session, task_data: TaskDelete)->bool:
        """Delete task"""
        task = db.query(Task).filter(
            Task.task_id == task_data.task_id,
            Task.user_id == task_data.user_id,
            Task.folder_id == task_data.folder_id
        ).first()

        if not task:
            return False
    
        db.delete(task)
        db.commit()
        return True

    @staticmethod
    def toggle_task_completion(db: Session, user_id: int, task_id:int)->Optional[Task]:
        """Change status of task"""
        task = db.query(Task).filter(
            Task.task_id == task_id,
            Task.user_id == user_id
        ).first()
        
        if task:
            task.is_completed = not task.is_completed
            db.commit()
            db.refresh(task)
        return task
