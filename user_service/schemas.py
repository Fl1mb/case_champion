from pydantic import BaseModel, EmailStr, Field
from datetime import datetime
from typing import Optional

#Schema for registration
class UserCreate(BaseModel):
    username : str = Field(..., min_length=3, max_length=50)
    email : EmailStr
    password : str = Field(..., min_length=8)
    full_name : str | None = None

class UserLogin(BaseModel):
    username : str = Field(..., min_length=3, max_length=50)
    password  : str = Field(..., min_length=8)

class UserGet(BaseModel):
    username : str = Field(..., min_length=3, max_length=50)

class UserResponse(BaseModel):
    id: int
    username: str
    email: EmailStr
    full_name: str | None
    is_active: bool
    created_at: datetime

    class Config:
        from_attributes = True # For ORM

class Token(BaseModel):
    access_token : str
    token_type : str

class TokenData(BaseModel):
    username : Optional[str] = None