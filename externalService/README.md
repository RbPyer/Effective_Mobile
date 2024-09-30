# External Service

Реализация внешнего сервиса, к которому обращается основной, сделано чисто для удобства.

Конечно, можно было написать мок, но я решил эту проблему иначе :)

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

```yaml
# config/config.yaml

env: "dev" # dev, prod
host: "0.0.0.0:8080"
timeout: 5s # time for request handling
idle_timeout: 60s  # life-time of client-server connection
```

## Запуск

Предварительно удостоверьтесь, что ```docker-compose.yaml``` соответствует ```config/.yam```.

```sh
docker-compose up --build -d
```