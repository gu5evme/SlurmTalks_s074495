# docker_container_deploy

Роль запуска контейнеров из собственных образов

Сборка образа выполняется на том же хосте, на котором будет запускаться контейнер

## Требования

На целевом сервере должен быть установлен docker

## Переменные

### Необходимо задать

- **docker_container_deploy_image_name** - целевое имя образа при сборке. Указывать в формате `name:tag`, пример `front:1.3.2`
- **docker_container_deploy_build_src** - папка с Dockerfile и всеми файлами для сборки образа. Не указывать `/` в конце пути

### Можно переопределить

- **docker_container_deploy_container_name** - целевое имя контейнера при запуске. Если явно не указать значение, имя будет сгенерировано автоматически
- **docker_container_deploy_build_args** - словарь параметров ARGs для передачи в Dockerfile на этапе сборки. По умолчанию переменная пуста. Пример:
    ```
    docker_container_deploy_build_args:
      NEXT_PUBLIC_BACKEND_URL: '{{ frontend_public_backend_url }}'
      OTHER_ARGUMENT: '123'
    ```
- **docker_container_deploy_run_env_vars** - словарь переменных окружения для передачи в Dockerfile на этапе запуска. По умолчанию переменная пуста. Пример:
    ```
    docker_container_deploy_run_env_vars:
      POSTGRES_USER: '{{ postgresql_user }}'
      OTHER_ENV_VAR: '321'
    ```
- **docker_container_deploy_run_restart_policy** - политика перезапуска контейнера. По умолчанию `unless-stopped`. [Возможные значения](https://docs.ansible.com/ansible/latest/collections/community/docker/docker_container_module.html#parameter-restart_policy)