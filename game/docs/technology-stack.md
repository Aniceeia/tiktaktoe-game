# Технологический стек проекта "Крестики-Нолики"

## Языки программирования

### Go (Golang)
- **Версия**: 1.24+
- **Платформы**: Linux, macOS, Windows
- **Архитектура**: x86_64, ARM64

## База данных

### PostgreSQL
- **Версия**: 15.14 (Alpine)
- **Драйвер**: `github.com/jackc/pgx/v5` v5.5.3
- **Кодировка**: UTF-8
- **Локаль**: C

## Веб-фреймворки и библиотеки

### HTTP Router
- **Gorilla Mux**: v1.8.1
- **Функции**: Роутинг, middleware, CORS

### Валидация
- **go-playground/validator**: v10.27.0
- **Функции**: Валидация структур, кастомные теги

### Криптография
- **golang.org/x/crypto/bcrypt**: v0.41.0
- **Функции**: Хеширование паролей

### UUID
- **github.com/google/uuid**: v1.6.0
- **Функции**: Генерация уникальных идентификаторов

## Контейнеризация

### Docker
- **Версия**: 20.10+
- **Base Image**: Alpine Linux
- **Multi-stage build**: Да
- **Non-root user**: Да

### Docker Compose
- **Версия**: 3.8
- **Сервисы**: PostgreSQL, Application
- **Volumes**: Persistent data
- **Networks**: Isolated

## Тестирование

### Unit Testing
- **Фреймворк**: Go testing
- **Coverage**: Поддерживается
- **Race detection**: Включен

### Integration Testing
- **curl**: HTTP запросы
- **jq**: JSON парсинг
- **bash**: Автоматизация

## Мониторинг и логирование

### Health Checks
- **HTTP endpoint**: `/health`
- **Docker health**: Встроенный
- **Database health**: pg_isready

### Логирование
- **Стандартный логгер**: Go log
- **Уровни**: INFO, ERROR, DEBUG

## Безопасность

### Аутентификация
- **Метод**: Basic Auth
- **Хеширование**: bcrypt
- **Токены**: UUID-based

## Производительность

## Развертывание

### Локальная разработка
- **Docker Compose**: Полный стек
- **Go run**: Только приложение

### Продакшен
- **Docker**: Рекомендуется
- **CI/CD**: Скрипты готовы

## Совместимость

### Операционные системы
- **Linux**: Ubuntu 20.04+, CentOS 8+
- **macOS**: 10.15+
- **Windows**: 10+ (WSL2)


## Зависимости

### Прямые зависимости
```go
github.com/gorilla/mux v1.8.1
github.com/go-playground/validator/v10 v10.27.0
golang.org/x/crypto v0.41.0
github.com/google/uuid v1.6.0
github.com/jackc/pgx/v5 v5.5.3
```

### Косвенные зависимости
```go
github.com/jackc/pgpassfile v1.0.0
github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a
github.com/jackc/puddle/v2 v2.2.1
golang.org/x/sync v0.16.0
```
