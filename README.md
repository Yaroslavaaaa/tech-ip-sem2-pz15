## Практическая работа №15. Вуйко Ярослава, ЭФМО-01-25
### Деплой приложения на VPS. Настройка system. 28.05.2026

### Подключение к VPS


Подключение  по SSH:

```bash
ssh root@XXX.XXX.XXX.XXX
```
<img width="762" height="280" alt="2026-05-28_22-34-15" src="https://github.com/user-attachments/assets/4f36f73e-2966-4ca3-83c2-32c54b497a87" />


После подключения была выполнена проверка доступа и обновление системы:

```bash
sudo apt update && sudo apt upgrade -y
```
<img width="875" height="267" alt="2026-05-28_23-14-52" src="https://github.com/user-attachments/assets/302352df-114e-42e4-a468-a4fb59dab46c" />



### Структура размещения приложения

На сервере приложение размещено в следующих директориях:

```text
/opt/tasks/tasks          — бинарный файл сервиса
/etc/tasks/tasks.env      — конфигурационный файл 
/etc/systemd/system/tasks.service — unit-файл systemd
```
<img width="938" height="158" alt="2026-05-28_23-20-23" src="https://github.com/user-attachments/assets/814e3323-45a2-4fa8-a156-c061a1814e4f" />




### Конфигурация systemd (tasks.service)


```ini
[Unit]
Description=Tasks Service
After=network.target

[Service]
Type=simple
User=tasksuser
WorkingDirectory=/opt/tasks

EnvironmentFile=/etc/tasks/tasks.env

ExecStart=/opt/tasks/tasks

Restart=always
RestartSec=2

NoNewPrivileges=true
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```


### Объяснение ключевых параметров:

- User=tasksuser
  Сервис запускается не от root, а от обычного пользователя — это безопаснее

- WorkingDirectory=/opt/tasks
  Папка, в которой запускается приложение

- EnvironmentFile
  Файл с переменными окружения

- ExecStart
  Команда, которая запускает бинарник

- Restart=always
  Если сервис упал — systemd автоматически его перезапустит

- NoNewPrivileges=true
  Запрещает сервису получать дополнительные права


### Статус сервиса systemd

Проверка состояния сервиса:

```bash
sudo systemctl status tasks
```


<img width="1024" height="402" alt="2026-05-28_23-23-30" src="https://github.com/user-attachments/assets/39d164ee-ff9c-4cc6-b864-935a487bea4a" />


### Логи сервиса (journalctl)

Просмотр логов:

```bash
sudo journalctl -u tasks -n 30 --no-pager
```

<img width="935" height="94" alt="2026-05-28_23-23-52" src="https://github.com/user-attachments/assets/ed510d69-2ae6-4c3d-a1a8-0e66c0bd062e" />



### Проверка доступности сервиса

Проверка health endpoint:

```bash
curl http://127.0.0.1:8082/health
```
<img width="553" height="155" alt="2026-05-28_23-26-44" src="https://github.com/user-attachments/assets/23a5a92c-3a23-4f8e-8855-c9c86ab2f8a2" />

<img width="534" height="84" alt="2026-05-28_23-24-55" src="https://github.com/user-attachments/assets/51451a52-4a1d-4fe2-94f5-09ab7fdfffbb" />

Демонстрация работы:

<img width="1365" height="788" alt="2026-05-28_23-25-48" src="https://github.com/user-attachments/assets/51f8dcad-e7d2-42e5-b74b-2615d0d31fe2" />

<img width="1375" height="813" alt="2026-05-28_23-25-54" src="https://github.com/user-attachments/assets/da02030e-a37d-4044-b3ba-af6cab817b88" />

## Процедура обновления и отката версии

## Обновление версии

Обновление выполняется следующим образом:

1. Сборка новой версии локально:

```bash
go build -o tasks ./cmd/tasks
```

2. Копирование на сервер:

```bash
scp tasks root@VPS_IP:/tmp/tasks
```

3. Остановка сервиса:

```bash
sudo systemctl stop tasks
```

4. Сохранение старой версии:

```bash
sudo mv /opt/tasks/tasks /opt/tasks/tasks.old
```

5. Установка новой версии:

```bash
sudo mv /tmp/tasks /opt/tasks/tasks
sudo chmod 755 /opt/tasks/tasks
sudo chown tasksuser:tasksuser /opt/tasks/tasks
```

6. Запуск сервиса:

```bash
sudo systemctl start tasks
```

<img width="953" height="673" alt="2026-05-28_23-34-42" src="https://github.com/user-attachments/assets/866d85a5-b946-4d6d-be4c-ee537cbd7b94" />


### Откат версии (rollback)

В случае ошибки выполняется откат:

```bash
sudo systemctl stop tasks
sudo mv /opt/tasks/tasks /opt/tasks/tasks.bad
sudo mv /opt/tasks/tasks.old /opt/tasks/tasks
sudo systemctl start tasks
```

<img width="957" height="375" alt="2026-05-28_23-35-36" src="https://github.com/user-attachments/assets/12d4129b-9463-4241-bf72-218076d008db" />


### Итог

В ходе работы был:

- развернут Go-сервис на VPS
- настроен systemd unit
- организовано безопасное хранение конфигурации
- реализован автоматический рестарт сервиса
- освоены команды управления и диагностики (systemctl, journalctl)
- выполнена процедура обновления и отката версии


## Контрольные вопросы

1. Зачем нужен systemd и чем он лучше screen/tmux?
systemd обеспечивает управление сервисами как системными демонами: автоматический запуск при старте системы, рестарт при сбоях, централизованные логи.
В отличие от screen/tmux, systemd гарантирует стабильность и управляемость сервиса.

2. Почему не стоит запускать сервис от root?
Запуск от root создаёт угрозу безопасности: при уязвимости в приложении злоумышленник получает полный доступ к системе.
Использование отдельного пользователя ограничивает потенциальный ущерб.

3. Зачем хранить env в /etc/..., а не в репозитории?
Конфигурация может содержать чувствительные данные и должна быть отделена от кода:

4. Как посмотреть логи сервиса, если он упал?
Используется systemd journal:

```bash
journalctl -u tasks -n 100
journalctl -u tasks -f
```

5. Что даёт Restart=always и RestartSec?
- Restart=always — автоматический перезапуск сервиса при любом падении
- RestartSec — задержка перед повторным запуском, чтобы избежать циклических рестартов




