# Backend

## Системная информация

- приложение работает на порту 8080
- ❗ Redis должен работать на порту 6379
- ❗ Postgres должен работать на порту 5432
- 📌 Необходимо изменить значения `.env ` файле, если вы запускаете базы данных с другими переменными окружения
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
    - ```
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
    -
        - **`tweetId=<tweet_id>` // Замените `<tweet_id>` на ID твита**
          **curl:**
    ```bash         
      curl -X POST http://localhost:8080/api/tweets/<tweet_id>/like \
    -H "Authorization: Bearer <your_token>"
- **POST /like-tweet/:tweetId/share**  
  URL: `http://localhost:8080/api/tweets/:tweetId/share`
    - **Authorization: Bearer Token**
    -
        - **`tweetId=<tweet_id>` // Замените `<tweet_id>` на ID твита**
          **curl:**
    ```bash         
      curl -X POST http://localhost:8080/api/tweets/<tweet_id>/share \
    -H "Authorization: Bearer <your_token>"
- **POST /like-tweet/:tweetId/slurm**  
  URL: `http://localhost:8080/api/tweets/:tweetId/slurm`
    - **Authorization: Bearer Token**
    -
        - **`tweetId=<tweet_id>` // Замените `<tweet_id>` на ID твита**
          **curl:**
    ```bash         
      curl -X POST http://localhost:8080/api/tweets/<tweet_id>/slurm \
    -H "Authorization: Bearer <your_token>"
- **GET /api/tweets/sorted**  
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
