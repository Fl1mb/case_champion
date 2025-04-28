# Используем официальный образ Python
FROM python:3.11-slim

# Устанавливаем системные зависимости (для компиляции grpcio и др.)
RUN apt-get update && apt-get install -y \
    gcc \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY requirements.txt .

# Устанавливаем Python-зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Копируем исходный код
COPY . .

# Генерируем gRPC-код из .proto (если нужно)
RUN python -m grpc_tools.protoc \
    -I. \
    --python_out=. \
    --grpc_python_out=. \
    user.proto

# Запускаем приложение (можно заменить на ваш скрипт)
CMD ["python", "./main.py"]