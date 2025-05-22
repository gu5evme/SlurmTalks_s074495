# SlurmTalks

## Запуск через docker compose

Для запуска приложения в docker compose нужно:
1. актуализировать значения переменных окружения в файлах:
    - **backend-service/.env** - переменные для Backend, описание в разделе [Backend/Системная информация](#backend)
    - **.env** - переменные для контейнеров с postgresql и redis. Значения одноименных переменных должны совпадать с тем что в backend-service/.env
    - переменная `NEXT_PUBLIC_BACKEND_URL` в docker-compose.yml в сервисе front. Описание переменной в разделе [Frontend/Системная информация](#frontend)
2. запустить приложение командой
    ```bash
    # slurmtalks/
    docker compose up -d
    ```

**Примечание**. Дублирование части переменных в .env файле в корне проекта сделано осознанно, хоть и не является лучшей практикой. Возможно следовало подключать файл backend-service/.env директивой `env_file` в compose файле у сервисов psql и redis, но тогда у них будут лишние переменные

## Развертывание через Ansible

Развернуть приложение можно через Ansible. В результате на указанных серверах запускаются docker контейнеры компонентов приложения

### Зависимости для Ansible controller

Перед использованием установить зависимости на Ansible контроллер:
- установка ролей из galaxy `ansible-galaxy install -r ansible_requirements.yml`
- установка коллекций из galaxy `ansible-galaxy collection install -r ansible_requirements.yml`

### Установка Docker

Плейбук по умолчанию устанавливает docker на все сервера из inventory, если определяет что там не запущен docker

Для настройки параметров установки docker ознакомьтесь с описанием роли [geerlingguy.docker](https://github.com/geerlingguy/ansible-role-docker)

### Настройка inventory

В inventory.ini указать хосты для размещения каждого компонета приложения в соответствующих группах:
  - frontend
  - backend
  - postgresql
  - redis

**Примечание**. При желании можно развернуть все компоненты на одном сервере

Для каждого сервера задать IP адрес в переменной `ansible_host`. Можно объявить как в `inventory.ini`, так и в персональных файлах в `host_vars`

Актуализировать переменные в:
- **all.yml** - общие переменные для всех компонетов приложения
- **postgresql.yml** - переменные для запуска контейнера postgresql
- **redis.yml** - переменные для запуска контейнера redis
- **backend.yml** - переменные для сборки и запуска контейнера backend
- **frontend.yml** - переменные для сборки и запуска контейнера frontend

**Примечание**. .env файлы не используются, значения из них нужно явно указать в переменных

Шифрование не использовано осознанно. При продуктивном использовании зашифровать переменные или файлы переменных через ansible-vault, или использовать стороннее хранилище секретов (например hashicorp vault)

### Развертывание приложения

Выполнить плейбук
```
ansible-playbook playbook.yml
```

# Frontend

## Системная информация

- ❗ ОСОБЕННОСТЬ СБОРКИ: приложение забирает переменные из ФАЙЛА .env, при подготовке образа для докер/кубернетес необходимо
  указать в `NEXT_PUBLIC_BACKEND_URL` URL бэкэнда **(не сервис внутри кубера/стенда, а внешний URL пример: `http://slurm-talks-back.sXXXXXX.edu.slurm.io`)**. Здесь должен быть URL, который будет доступен из браузера на клиентах, запускающих приложение. Используется для XHR запросов
- переменная `NEXT_PUBLIC_BACKEND_URL` должна быть инициализирована на этапе сборки приложения перед вызовом `npm build` - важно учитывать это при использовании подхода multistage сборки образа контейнера
- приложение работает на порту 3000
- `.env`
    - `NEXT_PUBLIC_BACKEND_URL=http://localhost:8080` - URL backend сервиса (при работе на стедне заменить на `http://slurm-talks-back.sXXXXXX.edu.slurm.io`)
- для работы приложения требуется backend
    - все что связано с постами требует коннекта к backend <-> redis
    - все что связано с пользователями требует коннекта к backend <-> postgresql

## Запуск

1. Запустить Backend

1. Установить зависимости

    Добавление репозитория nodesource
    ```bash
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
    ```

    Установка nodejs
    ```bash
    apt install nodejs
    ```

    Обновление npm
    ```bash
    npm install -g npm@11.3.0
    ```

1. Установка необходимых модулей NextJS
    ```bash
    # slurmtalks/frontend-service/
    npm ci
    ```

    **Примечание**. Команда `npm ci` устанавливает версии модулей согласно файлу package-lock.json. Если нужно игнорировать этот файл и установить свежие версии, используйте вместо нее команду `npm install`

    **Дополнительно**

    Обработка замечаний по результатам аудита, который проводится при запуске `npm install`
    ```bash
    # slurmtalks/backend-service/
    npm audit fix --force
    ```

1. Запуск приложения в тестовом режиме
    ```bash
    # slurmtalks/backend-service/
    npm run dev
    ```

    В этом режиме приложение работает медленно, но на консоль выводятся подробные логи

1. Сборка приложения для запуска в продуктивном режиме
    ```bash
    # slurmtalks/backend-service/
    npm run build
    ```

1. Запуск приложения в продуктивном режиме
    ```bash
    # slurmtalks/backend-service/
    npm start
    ```

## Роуты

### `/auth`

- `/sign-up`- регистрация пользователя (далее редирект на `/sign-in`)
    - *(для работы обязателен коннект бэкенда к postgresql)*
- `/sign-in` - аутентификация пользователя  (далее редирект на `/`)
    - *(для работы обязателен коннект бэкенда к postgresql)*

### `/`

- `/` - вывод всех постов
    - *(для работы обязателен коннект бэкенда к redis)*
- `/communities` - вывод списка пользователей
    - *(для работы обязателен коннект бэкенда к postgresql)*
- `/tweets/slurmed` - отсортированый вывод по возрастанию slurm всех постов
    - *(для работы обязателен коннект бэкенда к redis)*
- `/healthz` - health эндпоинт
- `/readyz` - readyz эндпоинт

# Backend

## Системная информация

- приложение работает на порту 8080
- ❗ Redis должен работать на порту 6379
- ❗ Postgres должен работать на порту 5432
- `.env`
    ```
  POSTGRES_USER=postgres # CHANGE ON YOUR POSTGRES USERNAME
  POSTGRES_PASSWORD= # CHANGE ON YOUR POSTGRES PASSWORD
  POSTGRES_DB=postgres # CHANGE ON YOUR POSTGRES DATABASE
  POSTGRES_HOST=localhost # CHANGE ON YOUR POSTGRES HOST
  POSTGRES_PORT=5432 # CHANGE ON YOUR POSTGRES PORT

  REDIS_ADDR=localhost:6379  # CHANGE ON YOUR REDIS HOST:PORT
  REDIS_PASSWORD= # CHANGE ON YOUR REDIS PASSWORD IF EXIST
  REDIS_DB=0

  JWT_SECRET_KEY=your_secret_key
  MIGRATE_IN_CODE=true # MIGRATION IN PGSQL IN CODE
    ```
- для работы приложения требуется backend
    - все что связано с постами требует коннекта к backend <-> redis
    - все что связано с пользователями требует коннекта к backend <-> postgresql

## Запуск

1. Запустить БД в контейнерах
    ```bash
    # slurmtalks/
    mkdir -p psql/data
    mkdir -p redis/data
    docker run --rm -d -v psql/data:/var/lib/postgresql/data -p 5432:5432 --network host --env-file backend-service/.env postgres:17.4
    docker run --rm -d -v redis/data:/data -p 6379:6379 --network host --env-file backend-service/.env redis:7.4.2-alpine
    ```

1. Установить зависимости
    ```bash
    apt install libc6-dev
    snap install go --channel=1.24/stable --classic
    ```

1. Запуск приложения в тестовом режиме
    ```bash
    # slurmtalks/backend-service/
    go run main.go
    ```

1. Компиляция приложения
    ```bash
    # slurmtalks/backend-service/
    go mod download
    go mod verify
    CGO_ENABLED=0 GOOS=linux go build -o /go/bin/backend
    ```

1. Запуск скомпилированного приложения
    ```bash
    /go/bin/backend
    ```

## API


- **POST /register**
  URL: `http://localhost:8080/api/register`
  **Body (raw JSON):**
    ```
    {
        "username": "<username>",      // Замените <username> на имя пользователя
        "email": "<email>",            // Замените <email> на email пользователя
        "password": "<password>"       // Замените <password> на пароль пользователя
    }
    ```

  **сurl:**
  ```bash
  curl -X POST http://localhost:8080/api/register \
    -H "Content-Type: application/json" \
    -d '{"username": "exampleuser", "email": "example@example.com", "password": "examplepassword"}'

- **POST /sign-in**
  URL: `http://localhost:8080/api/login`
  **Body (raw JSON):**
    ```
    {
    "username": "<username>",      // Замените <username> на имя пользователя
    "password": "<password>"       // Замените <password> на пароль пользователя
    }
    ```

  **curl:**
  ```bash
    curl -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"username": "exampleuser", "password": "examplepassword"}'

- **POST**
  URL: `http://localhost:8080/api/tweets`
    - **Authorization: Bearer Token:**
    - **Body (raw JSON):**
      ```
      {
          "text": "<tweet_text>",          // Замените <tweet_text> на текст вашего твита
          "userId": "<user_id>",           // Замените <user_id> на ID пользователя
          "username": "<username>"         // Замените <username> на имя пользователя
      }
      ```
      **curl:**
    - ```bash
      curl -X POST http://localhost:8080/api/tweets \
      -H "Authorization: Bearer <your_token>" \
      -H "Content-Type: application/json" \
      -d '{"text": "Hello, Twitter clone!", "userId": "1", "username": "<exampleuser>"}'

- **GET /get-tweet/:tweetId**
  URL: `http://localhost:8080/api/tweets/:tweetId`
    - **`tweetId=<tweet_id>` // Замените `<tweet_id>` на ID твита**
      **curl:**
  ```bash
        curl http://localhost:8080/api/tweets/<tweet_id>

- **GET /get-all-tweets**
  URL: `http://localhost:8080/api/tweets`
  **curl:**
  ```bash
        curl http://localhost:8080/api/tweets

- **POST /like-tweet/:tweetId/like**
  URL: `http://localhost:8080/api/tweets/:tweetId/like`
    - **Authorization: Bearer Token**
    - - **`tweetId=<tweet_id>` // Замените `<tweet_id>` на ID твита**
        **curl:**
    ```bash
      curl -X POST http://localhost:8080/api/tweets/<tweet_id>/like \
    -H "Authorization: Bearer <your_token>"
- **POST /like-tweet/:tweetId/share**
  URL: `http://localhost:8080/api/tweets/:tweetId/share`
    - **Authorization: Bearer Token**
    - - **`tweetId=<tweet_id>` // Замените `<tweet_id>` на ID твита**
        **curl:**
    ```bash
      curl -X POST http://localhost:8080/api/tweets/<tweet_id>/share \
    -H "Authorization: Bearer <your_token>"
- **POST /like-tweet/:tweetId/slurm**
  URL: `http://localhost:8080/api/tweets/:tweetId/slurm`
    - **Authorization: Bearer Token**
    - - **`tweetId=<tweet_id>` // Замените `<tweet_id>` на ID твита**
        **curl:**
    ```bash
      curl -X POST http://localhost:8080/api/tweets/<tweet_id>/slurm \
    -H "Authorization: Bearer <your_token>"
- **GET  /api/tweets/sorted**
  URL: `http://localhost:8080/api/tweets/sorted`
    - **Authorization: Bearer Token**
      **curl:**
    ```bash
  curl http://localhost:8080/api/tweets/sorted -H "Authorization: Bearer <your_token>"

### users
GET /community
URL: `http://localhost:8080/api/users`

Пример команды curl:
```bash
curl http://localhost:8080/api/users
```

### probe
GET /healthz
URL: `http://localhost:8080/healthz`
📕 Проба проходит после запуска сервиса

Пример команды curl:
```bash
curl http://localhost:8080/healthz
```

GET /readyz
URL: `http://localhost:8080/readyz`
📕 Проба проходит после присоединения сервиса к базам данных

Пример команды curl:
```bash
curl http://localhost:8080/readyz
```

### metrics
GET /metrics
URL: `http://localhost:8080/metrics`

Пример команды curl:
```bash
curl http://localhost:8080/metrics
```
