# Kubernetes

PostgreSQL и Redis развертываются через манифесты k8s

Frontend и Backend оформлены в виде Helm charts

## Сборка и доставка образов

PostgreSQL и Redis основаны на публичных образах, достаточно доступа к репозиторию. Backend и Frontend нужно предварительно собрать и выложить в доступный репозиторий. В рамках данного проекта использовался персональный публичный репозиторий Docker Hub

### Backend

Сборка через Docker, публикация в Docker Hub

1. Скачать репозиторий, перейти в slurmtalks/backend-service

2. Собрать образ

    ```
    docker build -t repo/name:tag .
    ```

    где:
    - **repo** - имя пользователя в Docker Hub
    - **name** - название образа
    - **tag** - версия

    Пример: `docker build -t menich/slurmtalks-backend:1.0 .`

3. Залогиниться в Docker Hub

    ```
    docker login -u username
    ```

4. Опубликовать образ в Docker Hub

    ```
    docker push repo/name:tag
    ```

5. На портале Docker Hub включить публичный режим для создавшегося репозитория

### Frontend

Сборка выполняется аналогично Backend с одним дополнением

На шаге 2 при сборке образа нужно добавить значение переменной **NEXT_PUBLIC_BACKEND_URL** (см. [frontend-service/README.md](../frontend-service/README.md)) в виде аргумента

```
docker build -t repo/name:tag --build-arg NEXT_PUBLIC_BACKEND_URL=http://slurm-talks-back.yourdomain.ru .
```

Из за особенностей фреймворка NextJS и метода сборки образа эту переменную не переопределить на этапе развертывания контейнера

Приложение всегда публикуется на порту 80

## Развертывание приложений

Скачать репозиторий, перейти в slurmtalks/k8s

### PostgreSQL и Redis

Манифесты в db

**Подготовка**

Выбрать узел k8s кластера, на котором будет создано локальное хранилище данных. Создать на этом узле каталоги. Для PostgreSQL и Redis можно использовать разные узлы

По умолчанию в манифесте psql.yaml в PV slurmtalks-psql-pv прописаны:
- путь: /opt/slurmtalks/psql-data
- узел: node-1.home.menich.ru

По умолчанию в манифесте redis.yaml в PV slurmtalks-redis-pv прописаны:
- путь: /opt/slurmtalks/redis-data
- узел: node-2.home.menich.ru

Отредактируйте значения под свой кластер

Также манифестах psql.yaml и redis.yaml можно отредактировать размер томов в PV и PVC. По умолчанию 5Gi

В манифесте env.yaml расположены общие переменные окружения для PostgreSQL и Redis. Значения этих переменных используются и в Backend/Frontend (описаны в values каждого Helm chart), так что они должны совпадать

**Развертывание**

```
kubectl apply -f ./db
```

Дождитесь пока контейнеры создадутся и станут Ready

### Backend

Helm chart в charts/slurmtalks-backend

**Подготовка**

Проверить и при необходимости отредактировать значения в values.yaml. Обратить внимание на:
- **image.repository** - указать имя своего образа
- **image.tag** - указать, если тег отличается от **appVersion** в Chart.yaml
- **configmap** и **secret** - переменные окружения для настройки приложения. Часть из них совпадает с db/env.yaml
- **ingress.className** - точно работает с nginx, другие не проверялись
- **ingress.hosts.host** - изменить на свое DNS имя, по которому приложение будет доступно извне (для NextJS и backend должен быть доступен извне)

**Развертывание**

```
helm install stb charts/slurmtalks-backend/
```

где **stb** - название экземпляра приложения, можно указать любое другое

Дождитесь пока контейнеры создадутся и станут Ready

### Frontend

Helm chart в charts/slurmtalks-frontend

**Подготовка**

Проверить и при необходимости отредактировать значения в values.yaml. Обратить внимание на:
- **image.repository** - указать имя своего образа
- **image.tag** - указать, если тег отличается от **appVersion** в Chart.yaml
- **ingress.className** - точно работает с nginx, другие не проверялись
- **ingress.hosts.host** - изменить на свое DNS имя, по которому приложение будет доступно извне

**Развертывание**

```
helm install stf charts/slurmtalks-frontend/
```

где **stf** - название экземпляра приложения, можно указать любое другое

Дождитесь пока контейнеры создадутся и станут Ready
