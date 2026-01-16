backend-сервис на Go — ядро системы геооповещений. Сервис интегрируется с новостным порталом (Django) через вебхуки.

КАК ЗАПУСТИТЬ:
- Windows 11
1. установите необходимые приложения:
    - Docker Desktop / Docker
    - TablePlus
    - Postman
    - pgAdmin 4 и расширение PostGIS
    - VS Code
    - Git
    - ngrok

2. выберите папку, куда хотите загрузить проект, и выполните следующие команды:
    - нажмите Win+R и введите cmd
    - cd [путь до вашей папки]
    - git clone https://github.com/BitCoinOffical/geo-announcements.git

3. перейдите в корень проекта и сделайте следующие действия:
    - измените .env.example на .env
    - откройте .env и измените конфиг под себя
    - обратите внимание, что в WEBHOOK_URL используется ngrok для имитации стороннего сервиса, куда будут отправляться webhooks

4. настройка ngrok
    - зайдите в свой личный кабинет ngrok, создайте/скопируйте свой токен и выполните команду:
      ngrok config add-authtoken ВАШ_AUTHTOKEN
    - в терминале запустите:
      ngrok http 9090
      В терминале появится ваш URL — его вставьте в .env в поле WEBHOOK_URL.

5. откройте Docker Desktop / Docker.

6. находясь в корне проекта через терминал VS Code или cmd выполните следующие команды:
    - docker build
    - docker-compose build
    - docker-compose up

7. готово! Если вы всё правильно настроили, можете открыть Postman:
    -для следующих запросов нужно установить в headers ключ и значение: X-API-KEY: [значение из вашего .env, поля API_KEY]
    минимальный и максимальный lat: -90, 90
    минимальный и максимальный lon: -180, 180
---
    POST http://localhost:8080/api/v1/incidents        # добавляет инциденты
    пример запроса:
    {
        "lat": 26.75,
        "lon": 133.25
    }
---
    GET http://localhost:8080/api/v1/incidents/?page={страница}&limit={лимит по кол-ву вывода инцидентов}         # выдает список инцидентов с пагинацией
    пример ответа для http://localhost:8080/api/v1/incidents/?page=1&limit=10:
    [
    {
        "Incident_id": 1,
        "Lat": 26.75,
        "Lon": 133.25,
        "Status": "public",
        "Create_at": "2026-01-16T20:47:34.244514Z",
        "Update_at": "2026-01-16T20:47:34.244514Z",
        "Deleted_at": null
    }
    ]
---
    GET http://localhost:8080/api/v1/incidents?id={id}      # выдает инцидент по его id
    пример ответа для http://localhost:8080/api/v1/incidents?id=1:
    {
    "Incident_id": 1,
    "Lat": 26.75,
    "Lon": 133.25,
    "Status": "public",
    "Create_at": "2026-01-16T20:47:34.244514Z",
    "Update_at": "2026-01-16T20:47:34.244514Z",
    "Deleted_at": null
    }
---
    PUT http://localhost:8080/api/v1/incidents?id={id}           # обновляет инциденты
    пример запроса для http://localhost:8080/api/v1/incidents?id=1:
    {
    "lat": 20,
    "lon": 100
    }
    и теперь при повторном выдаче по id будет следующий вывод:
    {
    "Incident_id": 1,
    "Lat": 20,
    "Lon": 100,
    "Status": "public",
    "Create_at": "2026-01-16T20:47:34.244514Z",
    "Update_at": "2026-01-16T20:55:11.014195Z",
    "Deleted_at": null
    }
---
    DELETE http://localhost:8080/api/v1/incidents?id={id}        # удаляет (скрывает) инциденты
---
    GET http://localhost:8080/api/v1/incidents/stats   # выдает статистику по зонам (сколько людей в какой зоне находится)
    пример вывод данных:
    {
        "Zone_id": 70,
        "UserCount": 1
    }

---
    GET http://localhost:8080/api/v1/system/health    # выдает статистику состояния сервера
    пример вывода:
    {
    "postgres": "ok",
    "redis": "ok",
    "status": 200
    }

    - эти эндпоинты можно вызывать без X-API-KEY. Если хотите, чтобы отправлялись координаты от конкретного человека, то нужно добавить X-Client-Id: [uuid человека]
    POST http://localhost:8080/api/v1/location/check   # проверяет ваши координаты и выдает список ближайших зон
    пример запроса:
    {
        "lat": 26,
        "lon": 33
    }
    пример вывода:
    {
    "dunger zones": [
        {
            "Zone_id": 103,
            "Lat": 29.051129,
            "Lon": 124,
            "Distant": "75"
        }
    ],
    "success": true
}

8. для завершения в Docker Desktop остановите контейнер и выполните команду:
    - docker compose down -v

- Arch Linux
1. установите необходимые приложения для этого можете использовать следующие команды:
    Docker
    - sudo pacman -S --needed docker docker-compose
    - sudo systemctl enable --now docker
    - sudo usermod -aG docker $USER
    - newgrp docker
    PostgreSQL
    - sudo pacman -S postgis postgresql
    Git
    - sudo pacman -S --needed git
    VScode
    - sudo pacman -S --needed code 
    BeeKeeper
    - sudo pacman -S --needed flatpak
    - flatpak install flathub io.beekeeperstudio.BeekeeperStudio
    Ngrok
    - тут вам надо будет скачать bin на сайте ngrok, а затем распаковать его в /usr/local/bin/
    - зайдите в свой личный кабинет ngrok, создайте/скопируйте свой токен и выполните команду:
      ngrok config add-authtoken ВАШ_AUTHTOKEN
    - в терминале запустите:
      ngrok http 9090
      В терминале появится ваш URL — его вставьте в .env в поле WEBHOOK_URL.
2. клонируйтее git репозитория:
    - выберете папку, куда клонировать
    - cd ~/projects
    - git clone https://github.com/BitCoinOffical/geo-announcements.git
    - cd geo-announcements

3. перейдите в корень проекта и сделайте следующие действия:
    - измените .env.example на .env
    - откройте .env и измените конфиг под себя
    - обратите внимание, что в WEBHOOK_URL используется ngrok для имитации стороннего сервиса, куда будут отправляться webhooks

4. запуск через Docker / Docker Compose для этого пропишите эти команды:
    - перейдите в корень проекта
    - запустите демон докера sudo systemctl start docker
    - docker build .
    - docker-compose build
    - docker-compose up
    - ВАЖНО - проверьте логи docker-compose logs -f
5. Готово! можно приступать к тестированию в консоле пропишите следующие команды:
---
    POST http://localhost:8080/api/v1/incidents        # добавляет инциденты
 
    curl -X POST http://localhost:8080/api/v1/incidents \
    -H "Content-Type: application/json" \
    -H "X-API-KEY: <API_KEY>" \
    -d '{"lat":26.75,"lon":133.25}'

---
    GET http://localhost:8080/api/v1/incidents/?page={страница}&limit={лимит по кол-ву вывода инцидентов}         # выдает список инцидентов с пагинацией

    curl "http://localhost:8080/api/v1/incidents/?page=1&limit=10" \
    -H "X-API-KEY: <API_KEY>"

    Пример ответа:
    [
    {
        "Incident_id": 1,
        "Lat": 26.75,
        "Lon": 133.25,
        "Status": "public",
        "Create_at": "2026-01-16T20:47:34.244514Z",
        "Update_at": "2026-01-16T20:47:34.244514Z",
        "Deleted_at": null
    }
    ]
---
    GET http://localhost:8080/api/v1/incidents?id={id}      # выдает инцидент по его id

    curl "http://localhost:8080/api/v1/incidents?id=1" -H "X-API-KEY: <API_KEY>"

    Пример ответа:
    {
    "Incident_id": 1,
    "Lat": 26.75,
    "Lon": 133.25,
    "Status": "public",
    "Create_at": "2026-01-16T20:47:34.244514Z",
    "Update_at": "2026-01-16T20:47:34.244514Z",
    "Deleted_at": null
    }
---
    PUT http://localhost:8080/api/v1/incidents?id={id}           # обновляет инциденты

    curl -X PUT "http://localhost:8080/api/v1/incidents?id=1" \
    -H "Content-Type: application/json" \
    -H "X-API-KEY: <API_KEY>" \
    -d '{"lat":20,"lon":100}'

    и теперь при повторном выдаче по id будет следующий вывод:
    {
    "Incident_id": 1,
    "Lat": 20,
    "Lon": 100,
    "Status": "public",
    "Create_at": "2026-01-16T20:47:34.244514Z",
    "Update_at": "2026-01-16T20:55:11.014195Z",
    "Deleted_at": null
    }
---
    DELETE http://localhost:8080/api/v1/incidents?id={id}        # удаляет (скрывает) инциденты
 
    curl -X DELETE "http://localhost:8080/api/v1/incidents?id=1" -H "X-API-KEY: <API_KEY>"

---
    GET http://localhost:8080/api/v1/incidents/stats   # выдает статистику по зонам (сколько людей в какой зоне находится)

    curl "http://localhost:8080/api/v1/incidents/stats" -H "X-API-KEY: <API_KEY>"

    пример вывод данных:
    {
        "Zone_id": 70,
        "UserCount": 1
    }
---
    GET http://localhost:8080/api/v1/system/health    # выдает статистику состояния сервера

    curl "http://localhost:8080/api/v1/system/health"

    пример вывода:
    {
    "postgres": "ok",
    "redis": "ok",
    "status": 200
    }

---
    эти эндпоинты можно вызывать без X-API-KEY. Если хотите, чтобы отправлялись координаты от конкретного человека, то нужно добавить X-Client-Id: [uuid человека]

    curl -X POST http://localhost:8080/api/v1/location/check \
    -H "Content-Type: application/json" \
    -H "X-Client-Id: <UUID_USER> \
    -d '{"lat":26,"lon":33}'

    пример вывода:
    {
    "dunger zones": [
        {
        "Zone_id": 103,
        "Lat": 29.051129,
        "Lon": 124,
        "Distant": "75"
        }
    ],
    "success": true
    }

6. для завершения в консоли выполните команду:
    - docker compose down -v

- macOS
Прошу прощения, я никогда не работал на macOS и не знаю, как запустить проект на macOS

## Архитектура проекта

```text

app-1
├── cmd                      # Точки входа приложения
│   ├── server               # HTTP-сервер
│   └── worker               # Воркер для асинхронных задач
│
├── config                   # Конфигурация backend-сервиса
│
├── internal                 # Внутренняя логика приложения
│   ├── adapters             # Адаптеры внешних зависимостей
│   │   └── secondary
│   │       ├── migration    # Миграции БД
│   │       ├── postgres     # Работа с PostgreSQL / PostGIS
│   │       └── redis        # Работа с Redis
│   │
│   ├── api                  # HTTP-уровень
│   │   ├── handlers         # HTTP-обработчики
│   │   │   └── mocks
│   │   ├── middleware       # Middleware
│   │   └── response         # Формирование HTTP-ответов
│   │
│   ├── domain               # Доменные правила
│   │   └── rules            # Кастомные валидаторы (lat / lon)
│   │
│   ├── interfaces           # Бизнес-логика и интерфейсы
│   │   └── http
│   │       ├── cache        # Кэширование
│   │       ├── dto          # DTO-структуры
│   │       ├── models       # Модели данных
│   │       ├── queue        # Очереди
│   │       ├── repo         # Репозитории
│   │       └── services     # Сервисы бизнес-логики
│   │           └── mocks
│   │
│   ├── pkg                  # Вспомогательные пакеты
│   ├── retry                # Механизм повторной отправки
│   └── worker               # Асинхронная отправка webhook (worker pool)
│
└── migrations               # SQL-миграции

Используется для тестирования и имитации внешнего сервиса, принимающего webhook-уведомления
app-2                        # сервер-заглушка
├── cmd
└── internal
    └── interfaces
        └── http
            ├── dto
            └── handlers
```

Литература:
- https://postgis.net/docs/manual-3.1/PostGIS_Special_Functions_Index.html?utm_source=chatgpt.com
- https://pkg.go.dev/github.com/data-dog/go-sqlmock#section-readme
- https://www.reddit.com/r/gis/comments/ush76v/how_to_work_out_if_point_is_within_polygon/?tl=ru
- https://habr.com/ru/articles/460535/
- https://youtu.be/ZEd0giJegVI?si=E9L7pINWFaQXcFtv
- https://www.reddit.com/r/golang/comments/1c7rnfp/golang_migrations_best_practices/?tl=ru
- https://habr.com/ru/companies/timeweb/articles/810857/
- https://www.reddit.com/r/golang/comments/v05rjw/how_to_write_unit_test_for_gin_golang/?tl=ru
- https://purpleschool.ru/knowledge-base/article/clean-architecture-go
- https://gin-gonic.com/ru/docs/testing/
- https://habr.com/ru/companies/first/articles/927460/
- https://www.reddit.com/r/golang/comments/1bk3hap/go_validator_or_gin_binding_validation_on_custom/?tl=ru