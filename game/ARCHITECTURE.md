# Архитектура проекта "Крестики-нолики"

## Обзор

Проект реализован с использованием **Clean Architecture** (Чистой архитектуры) на языке Go. Это обеспечивает разделение ответственности, тестируемость и масштабируемость кода.

## Структура проекта

```
game/
├── cmd/                    # Точка входа в приложение
├── internal/              # Внутренний код приложения
│   ├── domain/           # Бизнес-логика (ядро)
│   │   ├── entities/     # Бизнес-сущности
│   │   ├── repositories/ # Интерфейсы репозиториев
│   │   └── services/     # Бизнес-сервисы
│   ├── infrastructure/   # Внешние зависимости (БД, API)
│   └── interfaces/       # Интерфейсы (HTTP, CLI)
│       └── http/
│           ├── dto/      # Data Transfer Objects
│           ├── handlers/ # HTTP обработчики
│           └── middleware/ # HTTP middleware
├── api/                  # API документация
├── tests/               # Тесты
└── scripts/             # Скрипты развертывания
```

## Принципы Clean Architecture

### 1. **Domain Layer (Ядро)**
- **Entities**: Бизнес-сущности (Game, User)
- **Repositories**: Интерфейсы для работы с данными
- **Services**: Бизнес-логика

### 2. **Infrastructure Layer**
- Реализация репозиториев
- Работа с базой данных
- Внешние API

### 3. **Interface Layer**
- HTTP handlers
- DTO (Data Transfer Objects)
- Middleware

## Ключевые компоненты

### Entities (Сущности)

#### Game Entity
```go
type Game struct {
    ID         string     `json:"id"`
    Board      [3][3]int  `json:"board"`
    Status     GameStatus `json:"status"`
    Player1ID  string     `json:"player1_id"`
    Player2ID  string     `json:"player2_id"`
    NextTurnID string     `json:"next_turn_id"`
    Mode       GameMode   `json:"mode"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}
```

**Особенности:**
- Инкапсулирует игровую логику
- Проверяет валидность ходов
- Обновляет статус игры
- Определяет победителя

### Services (Сервисы)

#### GameService
- Создание игр
- Обработка ходов
- Управление PvE логикой
- Обновление счетов

#### AIService
- Алгоритм minimax для AI
- Поиск лучшего хода
- Оценка позиции

### DTO (Data Transfer Objects)

#### Входящие DTO
```go
type CreateGameRequest struct {
    Mode string `json:"mode" validate:"required,oneof=pvp pve"`
}

type MoveRequest struct {
    Row int `json:"row" validate:"required,min=0,max=2"`
    Col int `json:"col" validate:"required,min=0,max=2"`
}
```

#### Исходящие DTO
```go
type GameResponse struct {
    ID         string      `json:"id"`
    Board      [3][3]int   `json:"board"`
    Status     string      `json:"status"`
    Player1    PlayerInfo  `json:"player1"`
    Player2    *PlayerInfo `json:"player2,omitempty"`
    NextPlayer string      `json:"next_player"`
    Mode       string      `json:"mode"`
    Result     *GameResult `json:"result,omitempty"`
}
```

## Middleware Stack

### 1. **RecoveryMiddleware**
- Обрабатывает паники
- Логирует ошибки
- Возвращает 500 статус

### 2. **CORSMiddleware**
- Добавляет CORS заголовки
- Обрабатывает preflight запросы

### 3. **LoggingMiddleware**
- Логирует все HTTP запросы
- Записывает время выполнения

### 4. **AuthMiddleware**
- Проверяет Basic Auth
- Извлекает пользователя из БД
- Добавляет userID в контекст

### 5. **ValidationMiddleware**
- Проверяет Content-Type
- Валидирует JSON структуры

## Обработка ошибок

### Иерархия ошибок
```go
// Domain errors
ErrGameNotFound
ErrUserNotFound
ErrGameFull
ErrWrongTurn
ErrInvalidMove
ErrInvalidPosition
ErrGameNotRunning
ErrForbidden
ErrUnauthorized
```

### HTTP статусы
- `400 Bad Request` - неверные данные
- `401 Unauthorized` - не авторизован
- `403 Forbidden` - нет доступа
- `404 Not Found` - ресурс не найден
- `409 Conflict` - конфликт (игра полная)
- `500 Internal Server Error` - внутренняя ошибка

## PvE Логика

### Последовательность ходов
1. Игрок делает ход (X)
2. Система проверяет валидность
3. Обновляется статус игры
4. Если игра не завершена и режим PvE:
   - AI делает ход (O)
   - Обновляется статус игры
5. Возвращается обновленное состояние

### AI Алгоритм
- Использует алгоритм minimax
- Оценивает все возможные ходы
- Выбирает оптимальную стратегию
- Непроходим для игрока

## Валидация

### Уровни валидации
1. **DTO Level** - структура и типы данных
2. **Handler Level** - бизнес-правила
3. **Service Level** - игровая логика
4. **Entity Level** - инварианты домена

### Примеры валидации
```go
// DTO валидация
type MoveRequest struct {
    Row int `json:"row" validate:"required,min=0,max=2"`
    Col int `json:"col" validate:"required,min=0,max=2"`
}

// Бизнес-валидация
if !game.IsPlayerInGame(userID) {
    return ErrForbidden
}

if game.NextTurnID != userID {
    return ErrWrongTurn
}
```

## Рекомендации по улучшению

### 1. **Добавить тесты**
- Unit тесты для entities
- Integration тесты для services
- E2E тесты для API

### 2. **Улучшить логирование**
- Структурированные логи
- Уровни логирования
- Контекстная информация

### 3. **Добавить метрики**
- Prometheus метрики
- Мониторинг производительности
- Алерты

### 4. **Кэширование**
- Redis для сессий
- Кэширование игровых состояний
- Rate limiting

### 5. **WebSocket поддержка**
- Real-time обновления
- Уведомления о ходах
- Чат между игроками

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - регистрация
- `POST /api/v1/auth/login` - вход

### Игры
- `POST /api/v1/games` - создать игру
- `GET /api/v1/games` - список доступных игр
- `GET /api/v1/games/my` - игры пользователя
- `GET /api/v1/games/{id}` - состояние игры
- `POST /api/v1/games/{id}/join` - присоединиться к игре
- `POST /api/v1/games/{id}/move` - сделать ход

### Система
- `GET /health` - проверка здоровья

## Заключение

Данная архитектура обеспечивает:
- ✅ Разделение ответственности
- ✅ Тестируемость кода
- ✅ Масштабируемость
- ✅ Поддерживаемость
- ✅ Расширяемость

Проект готов к продакшену с минимальными доработками.
