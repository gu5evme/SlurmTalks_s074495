# SlurmTalks

- [Frontend](frontend-service)
- [Backend](backend-service)

## Запуск через docker compose

Для запуска приложения в docker compose нужно:
1. актуализировать значения переменных в:

    - **backend-service/.env** - переменные для Backend, описание в разделе [Backend/Системная информация](#backend)
    - переменная `NEXT_PUBLIC_BACKEND_URL` в docker-compose.yml в сервисе front. Описание переменной в разделе [Frontend/Системная информация](#frontend)

1. создать подпапки `data` в папках `psql` и `redis` для хранения данных контейнеров БД

    ```bash
    # slurmtalks/
    mkdir -p psql/data
    mkdir -p redis/data
    ```

1. запустить приложение командой

    ```bash
    # slurmtalks/
    docker compose up -d
    ```

## Развертывание через Ansible

Развернуть приложение можно через Ansible. В результате на указанных серверах запускаются docker контейнеры компонентов приложения

> .env файлы при развертывании через Ansible не используются, значения из них нужно явно указать в переменных

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

> При желании можно развернуть все компоненты на одном сервере

Для каждого сервера задать IP адрес в переменной `ansible_host`. Можно объявить как в `inventory.ini`, так и в персональных файлах в `host_vars`

Актуализировать переменные в:
- **all.yml** - общие переменные для всех компонетов приложения
- **postgresql.yml** - переменные для запуска контейнера postgresql
- **redis.yml** - переменные для запуска контейнера redis
- **backend.yml** - переменные для сборки и запуска контейнера backend
- **frontend.yml** - переменные для сборки и запуска контейнера frontend

Шифрование не использовано осознанно. При продуктивном использовании зашифровать переменные или файлы переменных через ansible-vault, или использовать стороннее хранилище секретов (например hashicorp vault)

### Развертывание приложения

Выполнить плейбук
```bash
ansible-playbook playbook.yml
```

## Развертывание в Kubernetes

> Для работы приложения нужны PostgreSQL и Redis. В тестовой среде для их развертывания можно воспользоваться манифестами в папке [k8s_mainfests](k8s_mainfests). Там используется local-storage и все сильно захардкожено, так что нужно отредактировать под свое окружение

Для развертывания компонетов Frontend и Backend приложения созданы helm чарты

Хранение образов контейнеров предполагается в Docker Hub

Настроен GitLab CI пайплайн. Чтобы им пользоваться, нужно задать переменные CI/CD

- **K8S_CA_CERT** - публичная часть CA сертификата кластера k8s. Для получения выполнить:

  ```bash
  kubectl get cm -o jsonpath="{.items[?(@.metadata.name==\"kube-root-ca.crt\")].data['ca\.crt']}"
  ```

- **K8S_CI_TOKEN** - токен сервисной УЗ в кластере k8s для деплоя. Создать УЗ и получить ее токен можно так:

  ```bash
  kubectl create serviceaccount ci-runner
  kubectl create clusterrolebinding ci-runner-binding --clusterrole=cluster-admin --serviceaccount=default:ci-runner
  kubectl apply -f - <<EOF
  apiVersion: v1
  kind: Secret
  metadata:
    name: ci-runner-token
    annotations:
      kubernetes.io/service-account.name: ci-runner
  type: kubernetes.io/service-account-token
  EOF
  kubectl get secret ci-runner-token -o jsonpath='{.data.token}' | base64 --decode
  ```

- **K8S_API_URL** - адрес API кластера k8s. Получить можно так:

  ```bash
  kubectl config view -ojsonpath='{.clusters[0].cluster.server}' ; echo
  ```

- **CI_REGISTRY_USER** - имя пользователя Docker Hub
- **CI_REGISTRY_TOKEN** - токен доступа для пользователя Docker Hub
