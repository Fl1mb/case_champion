from sqlalchemy import Column, Integer, String, Boolean, DateTime, ForeignKey
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime

#Base class for sqlAlchemy
Base = declarative_base()

class TaskFolder(Base):
    __tablename__ = "folders"

    folder_id = Column(Integer, primary_key=True, autoincrement=True)
    user_id = Column(Integer, nullable=False)  # Без ForeignKey!
    name = Column(String, nullable=False)
    created_at = Column(DateTime, default=datetime.now())

class Task(Base):
    __tablename__ = "tasks"

    # table fields
    task_id = Column(Integer, primary_key=True, autoincrement=True)
    folder_id = Column(Integer, ForeignKey('folders.folder_id', ondelete='CASCADE'))
    user_id = Column(Integer, nullable=False)
    title = Column(String, nullable=False)
    description = Column(String)
    due_time = Column(DateTime)
    priority = Column(Integer, default=1)
    is_completed = Column(Boolean, default=False)
    created_at = Column(DateTime, default=datetime.now())
    updated_at = Column(DateTime, default=datetime.now(), onupdate=datetime.now())

    # relations
    folder = relationship("TaskFolder", back_populates="tasks")

    

