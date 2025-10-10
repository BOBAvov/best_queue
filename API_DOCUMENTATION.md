# API Documentation

## Аутентификация

### POST /auth/sign-up
Регистрация нового пользователя
```json
{
  "username": "string",
  "password": "string",
  "tg_nick": "string",
  "group": "string"
}
```

### POST /auth/sign-in
Вход в систему
```json
{
  "tg_nick": "string",
  "password": "string"
}
```

## Пользователи (требует аутентификации)

### GET /api/profile
Получить профиль текущего пользователя

### PUT /api/profile
Обновить данные текущего пользователя
```json
{
  "username": "string",
  "tg_nick": "string",
  "group_id": 1,
  "is_admin": false
}
```

## Административные функции

### GET /api/admin/users
Получить список всех пользователей (только для админов)

### DELETE /api/admin/users/:id
Удалить пользователя (только для админов)

## Факультеты

### GET /api/faculties
Получить список всех факультетов

### GET /api/faculties/:code
Получить факультет по коду

### POST /api/faculties
Создать новый факультет (только для админов)
```json
{
  "code": "string",
  "name": "string",
  "comments": "string"
}
```

## Очереди

### GET /api/queues
Получить все очереди

### GET /api/queues/:id
Получить очередь по ID

### POST /api/queues
Создать новую очередь (только для админов)
```json
{
  "title": "string",
  "group_name": "string",
  "time_add": 5
}
```

### PUT /api/queues/:id
Обновить очередь (только для админов)

### DELETE /api/queues/:id
Удалить очередь (только для админов)

### POST /api/queues/:id/join
Присоединиться к очереди
```json
{
  "queue_id": 1
}
```

### DELETE /api/queues/:id/leave
Покинуть очередь

### GET /api/queues/:id/participants
Получить список участников очереди

### POST /api/queues/:id/shift
Сдвинуть очередь (удалить первого пользователя) - только для админов

## Заголовки авторизации

Для всех защищенных эндпоинтов используйте заголовок:
```
Authorization: Bearer <your_jwt_token>
```