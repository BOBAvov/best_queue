# API Система Электорнных Очередей

RESTful API для системы единого входа (Single Sign-On) с управлением пользователями, очередями и группами.

## Быстрый старт

```bash
# Запуск приложения
go run cmd/main.go

# Приложение будет доступно на http://localhost:8080
```

## Аутентификация

Все защищенные endpoints требуют JWT токен в заголовке `Authorization`:
```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### 🔐 Аутентификация

#### POST `/auth/sign-up`
Регистрация нового пользователя.

**Request Body:**
```json
{
  "username": "string",
  "password": "string", 
  "tg_nick": "string",
  "group": "string"
}
```

**Response (200):**
```json
{
  "id": 1,
  "message": "successfully signed up"
}
```

**Response (400):**
```json
{
  "error": "all fields are required"
}
```

#### POST `/auth/sign-in`
Вход пользователя в систему.

**Request Body:**
```json
{
  "tg_nick": "string",
  "password": "string"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (401):**
```json
{
  "error": "invalid credentials"
}
```

### 👤 Пользователи

#### GET `/api/profile`
Получение профиля текущего пользователя.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "user": {
    "id": 1,
    "username": "john_doe",
    "tg_nick": "@johndoe",
    "group_id": 1,
    "is_admin": false
  }
}
```

#### PUT `/api/profile`
Обновление профиля пользователя.

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "username": "new_username",
  "tg_nick": "@new_nick"
}
```

**Response (200):**
```json
{
  "message": "user updated successfully"
}
```

#### GET `/api/admin`
Проверка статуса администратора.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "admin": true
}
```

**Response (401):**
```json
{
  "message": "unauthorized"
}
```

#### GET `/api/admin/users`
Получение списка всех пользователей (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Response (200):**
```json
{
  "users": [
    {
      "id": 1,
      "username": "john_doe",
      "tg_nick": "@johndoe",
      "group_id": 1,
      "is_admin": false
    }
  ]
}
```

**Response (403):**
```json
{
  "error": "admin access required"
}
```

#### DELETE `/api/admin/users/:id`
Удаление пользователя (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Response (200):**
```json
{
  "message": "user deleted successfully"
}
```

**Response (400):**
```json
{
  "error": "invalid user id"
}
```

### 📋 Очереди

#### GET `/api/queues`
Получение списка всех очередей.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "queues": [
    {
      "id": 1,
      "title": "Консультация по математике",
      "time_start": "2024-01-15T10:00:00Z",
      "time_end": "2024-01-15T12:00:00Z"
    }
  ]
}
```

#### POST `/api/queues`
Создание новой очереди (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:**
```json
{
  "title": "string",
  "time_start": "2024-01-15T10:00:00Z",
  "time_end": "2024-01-15T12:00:00Z"
}
```

**Response (200):**
```json
{
  "id": 1,
  "message": "queue created successfully"
}
```

#### GET `/api/queues/:id`
Получение очереди по ID.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "queue": {
    "id": 1,
    "title": "Консультация по математике",
    "time_start": "2024-01-15T10:00:00Z",
    "time_end": "2024-01-15T12:00:00Z"
  }
}
```

**Response (404):**
```json
{
  "error": "queue not found"
}
```

#### PUT `/api/queues/:id`
Обновление очереди (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:**
```json
{
  "title": "Обновленное название",
  "time_start": "2024-01-15T11:00:00Z",
  "time_end": "2024-01-15T13:00:00Z"
}
```

**Response (200):**
```json
{
  "message": "queue updated successfully"
}
```

#### DELETE `/api/queues/:id`
Удаление очереди (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Response (200):**
```json
{
  "message": "queue deleted successfully"
}
```

#### POST `/api/queues/:id/join`
Присоединение к очереди.

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "queue_id": 1
}
```

**Response (200):**
```json
{
  "participant_id": 1,
  "message": "joined queue successfully"
}
```

#### DELETE `/api/queues/:id/leave`
Покидание очереди.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "left queue successfully"
}
```

#### GET `/api/queues/:id/participants`
Получение участников очереди.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "participants": [
    {
      "id": 1,
      "queue_id": 1,
      "user_id": 1,
      "position": 1,
      "joined_at": "2024-01-15T10:30:00Z",
      "is_active": true
    }
  ]
}
```

#### POST `/api/queues/:id/shift`
Сдвиг очереди - удаление первого участника (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Response (200):**
```json
{
  "message": "queue shifted successfully"
}
```

### 👥 Группы

#### GET `/api/groups`
Получение списка всех групп.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "code": "ИУ7-12Б",
    "comment": "Группа информационных технологий"
  }
]
```

#### POST `/api/groups`
Создание новой группы (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:**
```json
{
  "code": "ИУ7-13Б",
  "comment": "Новая группа"
}
```

**Response (200):**
```json
{
  "id": 1,
  "message": "group created"
}
```

#### GET `/api/groups/:id`
Получение группы по ID.

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "code": "ИУ7-12Б",
  "comment": "Группа информационных технологий"
}
```

#### PUT `/api/groups/:id`
Обновление группы (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Request Body:**
```json
{
  "code": "ИУ7-12Б-Updated",
  "comment": "Обновленное описание группы"
}
```

**Response (200):**
```json
{
  "message": "group updated"
}
```

#### DELETE `/api/groups/:id`
Удаление группы (только для админов).

**Headers:** `Authorization: Bearer <admin_token>`

**Response (200):**
```json
{
  "message": "group deleted"
}
```

### 📊 Статус API

#### GET `/`
Проверка статуса API.

**Response (200):**
```json
{
  "message": "API is running"
}
```

## Коды ошибок

| Код | Описание |
|-----|----------|
| 200 | Успешный запрос |
| 400 | Неверные данные запроса |
| 401 | Неавторизованный доступ |
| 403 | Недостаточно прав доступа |
| 404 | Ресурс не найден |
| 500 | Внутренняя ошибка сервера |

## Примеры использования

### Регистрация и вход
```bash
# Регистрация
curl -X POST http://localhost:8080/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123",
    "tg_nick": "@johndoe",
    "group": "ИУ7-12Б"
  }'

# Вход
curl -X POST http://localhost:8080/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "tg_nick": "@johndoe",
    "password": "password123"
  }'
```

### Работа с очередями
```bash
# Создание очереди (требует админ токен)
curl -X POST http://localhost:8080/api/queues \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "title": "Консультация по математике",
    "time_start": "2024-01-15T10:00:00Z",
    "time_end": "2024-01-15T12:00:00Z"
  }'

# Присоединение к очереди
curl -X POST http://localhost:8080/api/queues/1/join \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "queue_id": 1
  }'
```

## Технические детали

- **Фреймворк:** Gin (Go)
- **База данных:** PostgreSQL
- **Аутентификация:** JWT токены
- **Формат данных:** JSON
- **Порт:** 8080 (по умолчанию)

## Разработка

```bash
# Установка зависимостей
go mod download

# Запуск в режиме разработки
go run cmd/main.go

# Запуск тестов
go test ./...
```
