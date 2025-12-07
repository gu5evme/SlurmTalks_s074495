# PostgreSQL и Redis

PostgreSQL и Redis развертываются через манифесты k8s

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

Задать нужный namespace

**Развертывание**

```
kubectl apply -f k8s_mainfests
```

Дождитесь пока контейнеры создадутся и станут Ready
