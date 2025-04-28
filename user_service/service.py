from concurrent import futures
import grpc
from sqlalchemy.orm import Session
import user_pb2_grpc as user_grpc
import user_pb2 as pb2
import database
import logging
from models import User
from schemas import UserCreate, UserLogin, UserGet
from datetime import datetime, timedelta
import jwt
import bcrypt
import os
from typing import Optional

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Конфигурация JWT
SECRET_KEY = os.getenv("SECRET_KEY", "default-secret-key")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

class AuthService:
    @staticmethod
    def hash_password(password: str) -> str:
        """Хеширование пароля"""
        return bcrypt.hashpw(password.encode(), bcrypt.gensalt()).decode()

    @staticmethod
    def verify_password(plain_password: str, hashed_password: str) -> bool:
        """Проверка пароля"""
        return bcrypt.checkpw(plain_password.encode(), hashed_password.encode())

    @staticmethod
    def create_access_token(data: dict, expires_delta: Optional[timedelta] = None) -> str:
        """Создание JWT токена"""
        to_encode = data.copy()
        if expires_delta:
            expire = datetime.utcnow() + expires_delta
        else:
            expire = datetime.utcnow() + timedelta(minutes=15)
        to_encode.update({"exp": expire})
        return jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)

    @staticmethod
    def decode_token(token: str) -> Optional[dict]:
        """Декодирование JWT токена"""
        try:
            return jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        except jwt.PyJWTError as e:
            logger.error(f"Token decode error: {str(e)}")
            return None

class UserService(user_grpc.UserServiceServicer):
    def __init__(self, db: Session):
        super().__init__()
        self.db = db
    
    def CreateUser(self, request: pb2.CreateUserRequest, context):
        try:
            logger.info(f"Creating user: {request.username}")
            
            # Хешируем пароль перед сохранением
            hashed_password = AuthService.hash_password(request.password)
            logger.info(hashed_password)
            
            user = database.CreateUser(self.db, UserCreate(
                username=request.username,
                email=request.email,
                password=hashed_password,
                full_name=request.full_name if request.full_name else None
            ))
            
            return pb2.UserResponse(
                error="",
                id=user.id,
                username=user.username,
                email=user.email,
                full_name=user.full_name or ""
            )
            
        except ValueError as e:
            logger.error(f"Error creating user: {str(e)}")
            return pb2.UserResponse(error=str(e))
    
    def Login(self, request: pb2.LoginRequest, context):
        try:
            logger.info(f"Login attempt for user: {request.username}")
            
            user = database.GetUser(self.db, UserGet(username=request.username))
            if not user:
                logger.warning(f"User not found: {request.username}")
                return pb2.LoginResponse(access_token="")
            
            if not AuthService.verify_password(request.password, user.hashed_password):
                logger.warning(f"Invalid password for user: {request.username}")
                return pb2.LoginResponse(access_token="")
            
            access_token = AuthService.create_access_token(
                data={"sub": user.username},
                expires_delta=timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
            )
            
            logger.info(f"User {request.username} authenticated successfully")
            return pb2.LoginResponse(
                access_token=access_token,
            )
            
        except Exception as e:
            logger.error(f"Login error: {str(e)}")
            return pb2.LoginResponse(access_token="")
    
    def GetUser(self, request: pb2.GetUserRequest, context):
        pass

def serve():
    db = database.SessionLocal()
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_grpc.add_UserServiceServicer_to_server(UserService(db), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    logger.info("Server started on port 50051")
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down server...")
        server.stop(0)
        db.close()

if __name__ == '__main__':
    serve()