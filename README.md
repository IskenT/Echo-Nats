[![Github CI/CD](https://img.shields.io/github/workflow/status/IskenT/Echo-Nats/Go)](https://github.com/IskenT/Echo-Nats/actions)
[![Go Report](https://goreportcard.com/badge/github.com/IskenT/Echo-Nats)](https://goreportcard.com/report/github.com/IskenT/Echo-Nats)
![Repository Top Language](https://img.shields.io/github/languages/top/IskenT/Echo-Nats)
[![Scrutinizer Code Quality](https://img.shields.io/scrutinizer/quality/g/IskenT/Echo-Nats)](https://scrutinizer-ci.com/g/IskenT/Echo-Nats)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/IskenT/Echo-Nats)
[![Codacy Grade](https://img.shields.io/codacy/grade/c9467ed47e064b1981e53862d0286d65)](https://www.codacy.com/gh/IskenT/Echo-Nats/dashboard)
![Github Repository Size](https://img.shields.io/github/repo-size/IskenT/Echo-Nats)
![Github Open Issues](https://img.shields.io/github/issues/IskenT/Echo-Nats)
![Lines of code](https://img.shields.io/tokei/lines/github/IskenT/Echo-Nats)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/IskenT/Echo-Nats)
![GitHub contributors](https://img.shields.io/github/contributors/IskenT/Echo-Nats)
![Simply the best ;)](https://img.shields.io/badge/simply-the%20best%20%3B%29-orange)

# Echo-Nats
Тестовое задание от Hezzl.com
## О проекте
- Проект создан на Golang, Postgres, Clickhouse, Nats, Redis
- Описаны модели данных и миграции
- Для обращения в БД использованы raw sql запросы
## Инструкция к запуску
Для запуска проекта требуется docker-compose.

`docker-compose up --build --remove-orphans` или `make run`

Далее запустите миграции. 

`make postgres-up`

При необходимости в Docker-compose.yml могут быть добавлены клиентские службы для Clickhouse
## API
Endpoint = `http://localhost:8080/`

## A picture is worth a thousand words

<img src="./images/hezzl-run.PNG">

