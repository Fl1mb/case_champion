import subprocess
import time
import yaml 
from typing import Dict, Any

def is_compose_running(compose_file: str) -> bool:
    """Проверка, запущен ли docker-compose проект"""
    try:
        result = subprocess.run(
            ['docker-compose', '-f', compose_file, 'ps', '--services'],
            capture_output=True, 
            text=True,
            check=True
        )
        running_services = [s for s in result.stdout.splitlines() if s.strip()]
        return len(running_services) > 0
    except subprocess.CalledProcessError:
        return False
    except Exception as e:
        print(f"Ошибка проверки состояния docker-compose: {e}")
        return False

def deploy_docker_compose(compose_file: str = 'docker-compose.yml') -> bool:
    """Полное переразвертывание docker-compose"""
    try:
        # Проверка текущего состояния
        if is_compose_running(compose_file):
            print("Обнаружены запущенные сервисы, останавливаю...")
            subprocess.run(
                ['docker-compose', '-f', compose_file, 'down'],
                check=True,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.STDOUT
            )
            print("Все сервисы успешно остановлены")

        print(f"Запуск docker-compose из файла {compose_file}...")
        subprocess.run(
            ['docker-compose', '-f', compose_file, 'up', '-d'],
            check=True,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.STDOUT
        )
        print("Docker-compose успешно развернут")
        return True
        
    except subprocess.CalledProcessError as e:
        print(f"Ошибка при работе с docker-compose: {e}")
        return False
    except Exception as e:
        print(f"Неожиданная ошибка: {e}")
        return False

def get_running_services(compose_file: str) -> Dict[str, Dict[str, Any]]:
    """Получение списка сервисов из docker-compose и их статуса"""
    try:
        # Получаем список сервисов из docker-compose.yml
        with open(compose_file, 'r') as f:
            compose_config = yaml.safe_load(f)
            services = compose_config.get('services', {})
        
        # Получаем статус контейнеров
        result = subprocess.run(
            ['docker-compose', '-f', compose_file, 'ps', '--services', '--filter', 'status=running'],
            capture_output=True, 
            text=True
        )
        running_services = set(result.stdout.splitlines())
        
        service_info = {}
        for service in services:
            service_info[service] = {
                'should_run': True,
                'is_running': service in running_services,
                'restart_policy': services[service].get('restart', 'no')
            }
        
        return service_info
    except Exception as e:
        print(f"Ошибка при получении информации о сервисах: {e}")
        return {}

def restart_service(compose_file: str, service_name: str) -> bool:
    """Перезапуск указанного сервиса"""
    try:
        print(f"Перезапуск сервиса {service_name}...")
        subprocess.run(
            ['docker-compose', '-f', compose_file, 'restart', service_name],
            check=True,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.STDOUT
        )
        return True
    except subprocess.CalledProcessError as e:
        print(f"Ошибка при перезапуске сервиса {service_name}: {e}")
        return False

def monitor_services(compose_file: str, check_interval: int = 30):
    """Мониторинг сервисов и их автоматический перезапуск"""
    print(f"Начало мониторинга сервисов (интервал: {check_interval} сек)...")
    
    try:
        while True:
            services = get_running_services(compose_file)
            
            for service, info in services.items():
                if not info['is_running'] and info['should_run']:
                    print(f"Сервис {service} не работает!")
                    if restart_service(compose_file, service):
                        print(f"Сервис {service} успешно перезапущен")
                    else:
                        print(f"Не удалось перезапустить {service}")
            
            time.sleep(check_interval)
            
    except KeyboardInterrupt:
        print("\nМониторинг остановлен пользователем")
    except Exception as e:
        print(f"Критическая ошибка мониторинга: {e}")

if __name__ == "__main__":
    COMPOSE_FILE = 'docker-compose.yml'
    CHECK_INTERVAL = 5
    
    if deploy_docker_compose(COMPOSE_FILE):
        monitor_services(COMPOSE_FILE, CHECK_INTERVAL)