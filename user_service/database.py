from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from sqlalchemy.exc import SQLAlchemyError
from dotenv import load_dotenv
import os
from models import User
from schemas import *
import logging


load_dotenv() # Get .env file var's

# Данные для подключения (создайте файл .env)
DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://user:password@localhost/dbname")

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Создаем таблицы в базе (выполнить один раз)
def create_tables():
    from models import Base
    Base.metadata.create_all(bind=engine)



def CreateUser(db: Session, user_data: UserCreate):
    # Проверяем, нет ли такого пользователя
    existing_user = db.query(User).filter(
        (User.username == user_data.username) | 
        (User.email == user_data.email)
    ).first()
    
    if existing_user:
        raise ValueError("User already exists")
    
    # Создаем объект пользователя
    db_user = User(
        username=user_data.username,
        hashed_password=user_data.password,
        email=user_data.email,
        full_name=user_data.full_name
        )
    
    # Сохраняем в базу
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    
    return db_user

def GetUser(db: Session, user: UserGet):
    try:
        # Ищем пользователя по username
        existing_user = db.query(User).filter(
            User.username == user.username
        ).first()

        if not existing_user:
            logging.warning(f"User not found: {user.username}")
            return None

        logging.info(f"User retrieved: {existing_user.username}")
        return existing_user

    except SQLAlchemyError as e:
        logging.error(f"Database error while fetching user: {str(e)}")
        db.rollback()
        return None


    
