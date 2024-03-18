# Hezzl-Go
Тестовое задание от Hezzl.com
## О проекте
- Проект создан на Golang, Postgres, Clickhouse, Nats, Redis
- Описаны модели данных и миграции
- Для обращения в БД использованы raw sql запросы
## Инструкция к запуску
Для запуска проекта требуется docker-compose.

`docker-compose up -d` (--volumes) (--remove-orphans)

Далее запустите миграции. 

`make postgres-up`

Запустите проект.

`make run`

При необходимости в Docker-compose.yml могут быть добавлены клиентские службы для Clickhouse
## API
Endpoint = `http://localhost:8080/`
