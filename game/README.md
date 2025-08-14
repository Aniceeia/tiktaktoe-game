# Tic-Tac-Toe Game API

REST API для игры "Крестики-нолики" с поддержкой PvP и PvE режимов, реализованный с использованием Clean Architecture на Go.

## Особенности

- ✅ **Clean Architecture** - четкое разделение слоев
- ✅ **PvP режим** - игра между двумя игроками
- ✅ **PvE режим** - игра против AI с алгоритмом minimax
- ✅ **Аутентификация** - Basic Auth
- ✅ **PostgreSQL** - хранение данных
- ✅ **Docker** - контейнеризация
- ✅ **REST API** - стандартные HTTP endpoints

## Структура проекта

```
game/
├── cmd/main.go                    # Точка входа
├── internal/
│   ├── domain/                   # Бизнес-логика
│   │   ├── entities/             # Сущности (Game, User)
│   │   ├── repositories/         # Интерфейсы репозиториев
│   │   └── services/             # Бизнес-сервисы
│   ├── infrastructure/           # Внешние зависимости
│   │   └── repositories/         # Реализация репозиториев
│   └── interfaces/               # Интерфейсы
│       └── http/                 # HTTP API
│           ├── dto/              # Data Transfer Objects
│           ├── handlers/         # HTTP обработчики
│           ├── middleware/       # HTTP middleware
│           └── router.go         # Маршрутизация
├── docker/                       # Docker файлы
│   └── docker_files/             # Docker Compose
├── tests/                        # Тесты API
├── Dockerfile                    # Docker образ
└── docker-compose.yaml           # Docker Compose
```

## Быстрый старт

### 1. Запуск с Docker

```bash
# Клонируйте репозиторий
git clone <repository-url>
cd game

# Запустите приложение
cd docker/docker_files
docker-compose up -d

# Проверьте статус
docker-compose ps
```

### 2. Тестирование API

```bash
# Запустите полный тест
cd tests
chmod +x test_all.sh
./test_all.sh

# Или запустите тест с полной игрой
chmod +x test_api.sh
./test_api.sh
```

## API Endpoints

### Аутентификация

- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Вход в систему

### Игры

- `POST /api/v1/games` - Создать новую игру
- `GET /api/v1/games` - Список доступных игр
- `GET /api/v1/games/my` - Игры пользователя
- `GET /api/v1/games/{id}` - Состояние игры
- `POST /api/v1/games/{id}/join` - Присоединиться к игре
- `POST /api/v1/games/{id}/move` - Сделать ход

### Пользователи

- `GET /api/v1/users/{id}` - Информация о пользователе

### Система

- `GET /health` - Проверка здоровья сервиса

## Примеры использования

### Регистрация и вход

```bash
# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}'

# Вход
curl -X POST http://localhost:8080/api/v1/auth/login \
  -u alice:secret1
```

### Создание PvP игры

```bash
# Создание игры
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pvp"}'

# Присоединение к игре
curl -X POST http://localhost:8080/api/v1/games/{game-id}/join \
  -u bob:secret2
```

### Создание PvE игры

```bash
# Создание игры против AI
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pve"}'
```

### Выполнение хода

```bash
# Ход в позицию (1, 1)
curl -X POST http://localhost:8080/api/v1/games/{game-id}/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"row": 1, "col": 1}'
```

## Архитектура

Проект следует принципам Clean Architecture:

### Domain Layer (Ядро)
- **Entities**: Бизнес-сущности (Game, User)
- **Repositories**: Интерфейсы для работы с данными
- **Services**: Бизнес-логика

### Infrastructure Layer
- Реализация репозиториев
- Работа с PostgreSQL
- Внешние зависимости

### Interface Layer
- HTTP handlers
- DTO (Data Transfer Objects)
- Middleware (Auth, CORS, Logging, Recovery)

## Middleware Stack

1. **RecoveryMiddleware** - обработка паник
2. **CORSMiddleware** - CORS заголовки
3. **LoggingMiddleware** - логирование запросов
4. **AuthMiddleware** - аутентификация
5. **ValidationMiddleware** - валидация данных

## База данных

PostgreSQL с таблицами:
- `users` - пользователи
- `games` - игры

## Разработка

### Локальная разработка

```bash
# Установите зависимости
go mod download

# Запустите PostgreSQL
docker run -d --name postgres \
  -e POSTGRES_USER=game \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=game_db \
  -p 5432:5432 \
  postgres:14

# Запустите приложение
go run cmd/main.go
```

### Тестирование

```bash
# Запустите тесты
go test ./...

# Запустите интеграционные тесты
cd tests
./test_all.sh
```

## Переменные окружения

- `DATABASE_URL` - URL подключения к PostgreSQL
- `PORT` - Порт для HTTP сервера (по умолчанию 8080)

## Лицензия

MIT License
