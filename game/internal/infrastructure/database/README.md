# Database Schema

## Обзор

Схема базы данных для игры "Крестики-Нолики" реализована в соответствии с принципами Clean Architecture.

## Структура

- `init.sql` - Основная схема базы данных
- `migration.sql` - Миграции для обновления схемы

## Таблицы

### users
Хранит информацию о пользователях системы.

```sql
CREATE TABLE users (
    uuid VARCHAR(36) PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    score INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### games
Хранит информацию об играх.

```sql
CREATE TABLE games (
    id VARCHAR(36) PRIMARY KEY,
    board JSONB NOT NULL,
    status VARCHAR(50) NOT NULL,
    player1_id VARCHAR(36) REFERENCES users(uuid) ON DELETE SET NULL,
    player2_id VARCHAR(36),
    next_turn_id VARCHAR(36),
    mode VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Индексы

```sql
CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_games_player1_id ON games(player1_id);
CREATE INDEX idx_games_player2_id ON games(player2_id);
CREATE INDEX idx_games_created_at ON games(created_at);
```

## Миграции

Для обновления схемы используйте файл `migration.sql`.

## Использование

### Локальная разработка
```bash
# Подключение к БД
docker exec -it game-db psql -U game game_db

# Применение схемы
docker exec -it game-db psql -U game game_db -f /docker-entrypoint-initdb.d/init.sql
```

### Docker
Схема автоматически применяется при первом запуске контейнера PostgreSQL.
