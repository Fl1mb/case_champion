from sqlalchemy.orm import Session
from models import User
from schemas import UserCreate
from database import create_tables
from service import serve, logger
from time import sleep

sleep(3)
logger.info("slept")

create_tables()
logger.info("tables are created")

# Пример использования
if __name__ == "__main__":
    serve()
    
