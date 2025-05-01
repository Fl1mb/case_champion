# Используем официальный образ Python
FROM python:3.11-slim

# Устанавливаем системные зависимости
RUN apt-get update && apt-get install -y \
    gcc \
    python3-dev \
    curl \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY requirements.txt .

# Устанавливаем Python-зависимости и плагины
RUN pip install --no-cache-dir -r requirements.txt && \
    pip install protoc-gen-validate && \
    export PATH=$PATH:$(python -m site --user-base)/bin

# Создаем структуру директорий для proto файлов
RUN mkdir -p /app/proto/validate

# Загружаем validate.proto в правильную директорию
RUN curl -o /app/proto/validate/validate.proto \
    https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/main/validate/validate.proto

# Копируем исходный код
COPY . .

# Компилируем proto-файлы
RUN python -m grpc_tools.protoc \
    -I/app/proto \
    -I/app \
    --python_out=/app \
    --grpc_python_out=/app \
    /app/proto/task.proto

RUN python -m grpc_tools.protoc \
    -I/app/proto \
    -I/app \
    --python_out=/app \
    --grpc_python_out=/app \
    /app/proto/validate/validate.proto

# Запускаем приложение
CMD ["python", "./main.py"]