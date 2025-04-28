from sqlalchemy.orm import Session
from sqlalchemy import desc, func
from models import *
from schemas import *

class FolderRepo:
    @staticmethod
    def create_folder(db : Session, folder_data: FolderCreate)->TaskFolder:
        """Create new folder"""
        db_folder = TaskFolder(
            user_id = folder_data.user_id,
            name = folder_data.folder_name
        )
        db.add(db_folder)
        db.commit() 
        db.refresh(db_folder)
        return db_folder
    
    @staticmethod
    def get_folder(db: Session, user_id: int, folder_id:int)->Optional[TaskFolder]:
        """Get folder by id and check user"""
        return db.query(TaskFolder).filter(
            TaskFolder.folder_id == folder_id,
            TaskFolder.user_id == user_id
        ).first()

    @staticmethod
    def get_user_folders(db: Session, user_id: int)->List[TaskFolder]:
        """Get all user's folders"""
        return db.query(TaskFolder).filter(
            TaskFolder.user_id == user_id
        ).order_by(TaskFolder.created_at).all()
    
    @staticmethod
    def update_folder(db: Session, folder_data : FolderUpdate)->Optional[TaskFolder]:
        """Update folder's name"""
        folder = db.query(TaskFolder).filter(
            TaskFolder.user_id == folder_data.user_id,
            TaskFolder.folder_id == folder_data.folder_id
        ).first()

        if folder:
            folder.name = folder_data.new_folder_name
            db.commit()
            db.refresh(folder)
        return folder
    
    @staticmethod
    def delete_folder(db: Session, delete_data: FolderDelete)->bool:
        """Deleting all tasks in folder plus folder"""
        folder = db.query(TaskFolder).filter(
            TaskFolder.user_id == delete_data.user_id,
            TaskFolder.folder_id == delete_data.folder_id
        ).first()

        if not folder:
            return False

        db.query(Task).filter(
            Task.folder_id == delete_data.folder_id
        ).delete()

        db.delete(folder)
        db.commit()

        return True

