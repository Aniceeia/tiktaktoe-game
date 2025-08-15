## Требования

- Go 1.24+
- Docker и Docker Compose
- Git

## Доступные скрипты

### Сборка проекта
```bash
./scripts/build.sh
```
### Деплой
```bash
./scripts/deploy.sh
```

### Тестирование
```bash
./scripts/test.sh
```
### Очистка
```bash
./scripts/cleanup.sh
```

### Ручная сборка


### 1. Запуск с помощью Docker

```bash
./scripts/deploy.sh
docker-compose -f docker/docker-compose.yaml ps
docker-compose -f docker/docker-compose.yaml logs -f
```

### 2. Локальная разработка

```bash
go mod tidy
cd docker && docker-compose up -d db && cd ..
go build -o game cmd/main.go
./game
```

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Авторизация пользователя

### Игры
- `POST /api/v1/games` - Создание новой игры
- `GET /api/v1/games` - Получение списка доступных игр
- `GET /api/v1/games/{id}` - Получение состояния игры
- `POST /api/v1/games/{id}/join` - Присоединение к игре
- `POST /api/v1/games/{id}/move` - Выполнение хода
- `GET /api/v1/games/my` - Получение игр пользователя

### Пользователи
- `GET /api/v1/users/{id}` - Получение информации о пользователе

### Система
- `GET /health` - Проверка здоровья сервиса

## Примеры использования

### Регистрация и авторизация
```bash
# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}'

# Авторизация
curl -X POST http://localhost:8080/api/v1/auth/login \
  -u alice:secret1
```

### Создание и игра
```bash
# Создание PvP игры
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pvp"}'

# Присоединение к игре
curl -X POST http://localhost:8080/api/v1/games/{game_id}/join \
  -u bob:secret2

# Выполнение хода
curl -X POST http://localhost:8080/api/v1/games/{game_id}/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"row": 1, "col": 1}'
```

## Структура проекта

```
game/
├── cmd/
│   └── main.go              # Точка входа приложения
├── internal/
│   ├── domain/              # Доменный слой (Clean Architecture)
│   │   ├── entities/        # Бизнес-сущности
│   │   ├── repositories/    # Интерфейсы репозиториев
│   │   └── services/        # Бизнес-логика
│   ├── infrastructure/      # Инфраструктурный слой
│   │   ├── repositories/    # Реализации репозиториев
│   │   └── database/        # Схема базы данных
│   └── interfaces/          # Интерфейсный слой
│       └── http/            # HTTP handlers и middleware
├── docker/                  # Docker конфигурация
├── scripts/                 # CI/CD скрипты
├── tests/                   # Тесты
└── docs/                    # Документация
```

## Переменные окружения

- `DATABASE_URL` - URL подключения к PostgreSQL
- `PORT` - Порт для HTTP сервера (по умолчанию: 8080)

## Мониторинг

### Логи
```bash
docker-compose -f docker/docker-compose.yaml logs -f app
docker-compose -f docker/docker-compose.yaml logs -f db
```

## Устранение неполадок

### Проблемы с базой данных
```bash
docker-compose -f docker/docker-compose.yaml down -v
docker-compose -f docker/docker-compose.yaml up -d
```

### Проблемы с портами
```bash
lsof -i :8080
lsof -i :5432
sudo lsof -ti:8080 | xargs sudo kill -9
```

## Разработка

### Добавление новых тестов
1. Создайте тест в папке `tests/`
2. Добавьте вызов в `scripts/test.sh`

### Изменение схемы БД
1. Обновите `internal/infrastructure/database/init.sql`
2. Пересоздайте контейнеры: `./scripts/cleanup.sh && ./scripts/deploy.sh`
