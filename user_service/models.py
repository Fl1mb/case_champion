from sqlalchemy import Column, Integer, String, Boolean, DateTime
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime
from passlib.context import CryptContext

#Base class for sqlAlchemy
Base = declarative_base()

class User(Base):
    __tablename__ = "users" # table name in postgres

    #Поля таблицы
    id = Column(Integer, primary_key=True, index=True) #ID пользователя в таблице
    username = Column(String(50), unique=True, nullable=False) #Логин пользователя
    email = Column(String(50), unique=True, nullable=False) #Почта пользователя
    hashed_password = Column(String(255), nullable=False) #Захэшированный пароль
    full_name = Column(String(100), nullable=False) #Полное имя
    is_active = Column(Boolean, default=True)   #Онлайн ли пользователь
    created_at = Column(DateTime, default=datetime.now)  #Дата регистрации

    

