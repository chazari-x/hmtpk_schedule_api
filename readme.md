# API ХМТПК Расписание

[![Go Reference](https://pkg.go.dev/badge/github.com/chazari-x/hmtpk_schedule_api.svg)](https://pkg.go.dev/github.com/chazari-x/hmtpk_schedule_api)
[![Go Report Card](https://goreportcard.com/badge/github.com/chazari-x/hmtpk_schedule_api)](https://goreportcard.com/report/github.com/chazari-x/hmtpk_schedule_api)
![License](https://img.shields.io/github/license/chazari-x/hmtpk_schedule_api)
[![Group](https://img.shields.io/badge/VK-Subscripe-blue)](https://vk.com/club218199825)

Этот репозиторий представляет собой API для получения списка групп, преподавателей и расписания [Ханты-Мансийского технолого-педагогического колледжа (ХМТПК)](https://hmtpk.ru/ru/).

### Описание

API предоставляет функциональность HTTP и gRPC серверов для получения информации о группах, преподавателях и расписании занятий.

### Особенности

- **HTTP сервер**: Предоставляет конечные точки для получения списка групп, преподавателей и общего расписания.
- **gRPC сервер**: Реализует методы gRPC для получения списка групп, преподавателей и деталей расписания.
- **Докеризованное развертывание**: Включает конфигурацию Docker для удобного развертывания.

### Использование

Перед запуском необходимо создать файл [config.docker.yaml](etc/config.docker.yaml) в директории [etc](etc) на основе файла [config.docker.example.yaml](etc/config.docker.example.yaml), который находится в той же директории. Вот пример содержимого файла [config.docker.yaml](etc/config.docker.example.yaml):

```yaml
redis:
  host: redis
  port: :6379
  password: password
server:
  http:
    domain: hmtpksch.chazari.ru
    http_address: :80
    https_address: :443
  grpc:
    domain: hmtpksch.chazari.ru
    address: :50051
```

#### Запуск HTTP сервера

Чтобы запустить HTTP сервер, выполните следующую команду:

```bash
docker-compose up http-schedule-server
```

Это запустит HTTP сервер на порту 80 для HTTP и порту 443 для HTTPS.

#### Запуск gRPC сервера

Чтобы запустить gRPC сервер, выполните следующую команду:

```bash
docker-compose up grpc-schedule-server
```

Это запустит gRPC сервер на порту 50051.

#### Build Docker Images

Чтобы создать Docker образы при обновлении кода, выполните следующие команды:

```bash
docker-compose build http-schedule-server
```

```bash
docker-compose build grpc-schedule-server
```

### API Endpoints

#### HTTP Конечные Точки

- `GET /groups`: Получить список групп.
- `GET /teachers`: Получить список преподавателей.
- `GET /schedule`: Получить расписание на основе предоставленных параметров.

#### Методы gRPC

- `GetGroups`: Получить список групп.
- `GetTeachers`: Получить список преподавателей.
- `GetSchedule`: Получить расписание на основе предоставленных параметров.

### Определения Protobuf

```protobuf
syntax = "proto3";

option go_package = "../protobuf";

package Schedule;

message Request {
  string token = 1;
}

message Response {
  string message = 1;
}

message ScheduleRequest {
  string token = 1;
  string date = 2;
  string group = 3;
  string teacher = 4;
}

message ScheduleResponse {
  string message = 1;
}

service Schedule {
  rpc GetGroups (Request) returns (Response);
  rpc GetTeachers (Request) returns (Response);
  rpc GetSchedule (ScheduleRequest) returns (ScheduleResponse);
}
```

### Конфигурация Docker Compose

```yaml
version: '3'
services:
  redis:
    container_name: redis
    image: redis:6.2-alpine
    restart: always
    ports:
      - '6379:6379'
    command: redis-server --save 20 1 --loglevel warning --requirepass password
  http-schedule-server:
    container_name: http-schedule-server
    build:
      context: .
      dockerfile: Dockerfile-http
    ports:
      - "80:80"
      - "443:443"
    restart: always
    command: [ "/app/main", "http" ]
  grpc-schedule-server:
    container_name: grpc-schedule-server
    build:
      context: .
      dockerfile: Dockerfile-grpc
    ports:
      - "50051:50051"
    restart: always
    command: [ "/app/main", "grpc" ]
```

### OpenAPI Specification

Файл [openAPI.yaml](openAPI.yaml) содержит спецификацию API в формате OpenAPI.

### Лицензия

Этот проект лицензирован под [лицензией MIT](LICENSE).