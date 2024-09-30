# Go Songs Service (требуемый сервис)

Реализация требуемого сервиса.
 
## Требования

- docker
- docker-compose

## Установка

1. **Клонируйте репозиторий:**

    ```sh
    git clone git@github.com:RbPyer/Effective_Mobile.git
    cd Effective_Mobile
    ```

2. **Создайте файл окружения `.env` в директории `config` и добавьте необходимые переменные:**

```dotenv
# config/.env

# Environment (dev, prod)
ENV=

# Database connection settings
DB_HOST=        # Host and port of the Postgres server
DB_NAME=        # Name of the database
DB_USER=        # Database user
DB_PASSWORD=    # Password for the database user

# HTTP server configuration
HTTP_SERVER_ADDRESS=               # Address and port where the HTTP server will run
HTTP_SERVER_TIMEOUT=               # Timeout for handling requests
HTTP_SERVER_IDLE_TIMEOUT=          # Timeout for keeping connections open

# Cache server configuration (Redis)
CACHE_ADDRESS=   # Address and port of the Redis server
```

## Запуск

Предварительно удостоверьтесь, что ```docker-compose.yaml``` соответствует ```config/.env```.

```sh
docker-compose up --build -d
```