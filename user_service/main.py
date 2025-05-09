import threading
from prometheus_client import start_http_server, Counter, Gauge, Histogram
import service
import time
from database import create_tables

# Метрики
REQUEST_COUNT = Counter(
    'http_requests_total',
    'Total HTTP Requests',
    ['method', 'endpoint', 'status_code']
)

REQUEST_LATENCY = Histogram(
    'http_request_duration_seconds',
    'HTTP request latency',
    ['endpoint']
)

ERROR_COUNT = Counter(
    'http_errors_total',
    'Total HTTP Errors',
    ['error_type']
)

DB_QUERY_TIME = Histogram(
    'db_query_duration_seconds',
    'Database query duration',
    ['query_type']
)

def run_metrics_server():
    """Запуск сервера метрик в отдельном потоке"""
    start_http_server(8053)

if __name__ == '__main__':
    # Запуск сервера метрик в фоновом режиме
    metrics_thread = threading.Thread(target=run_metrics_server, daemon=True)
    metrics_thread.start()
    
    time.sleep(3)  # Ожидание инициализации
    create_tables()
    service.serve()
    
    