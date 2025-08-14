# Инструкции по тестированию

## Подготовка к тестированию

### 1. Запуск приложения

```bash
# Перейдите в папку с Docker файлами
cd docker/docker_files

# Запустите приложение
docker-compose up -d

# Проверьте статус
docker-compose ps
```

### 2. Проверка работы сервисов

```bash
# Проверьте, что база данных запущена
docker exec -it docker_files-db-1 psql -U game game_db -c "SELECT version();"

# Проверьте health check
curl http://localhost:8080/health
```

## Запуск тестов

### Полный тест (test_all.sh)

Этот тест проверяет все основные функции API:

```bash
cd tests
./test_all.sh
```

**Что тестируется:**
- ✅ Регистрация и аутентификация пользователей
- ✅ Создание PvP игр
- ✅ Присоединение к играм
- ✅ Создание PvE игр
- ✅ Выполнение ходов в PvE (игрок → AI)
- ✅ Получение списка доступных игр
- ✅ Получение игр пользователя
- ✅ Обработка ошибок
- ✅ Health check

### Тест с полной игрой (test_api.sh)

Этот тест проводит полную игру до победы:

```bash
cd tests
./test_api.sh
```

**Что тестируется:**
- ✅ Регистрация пользователей
- ✅ Создание PvP игры
- ✅ Последовательные ходы до победы
- ✅ Обновление счета победителя
- ✅ Получение информации о пользователях

## Ожидаемые результаты

### Успешный тест

При успешном выполнении вы увидите:

```
========== PvP TESTS ==========
1. Creating a PvP game...
Game response: {"id":"...","status":"player1_turn",...}

========== PvE TESTS ==========
4. Creating a PvE game...
PvE Game response: {"id":"...","status":"player1_turn",...}

8. Checking if AI made a move...
AI successfully made a move

========== GAME LIST TESTS ==========
12. Checking available games...
Available games: {"games":[...],"total":3}

Test completed successfully!
```

### Проверка AI в PvE

В PvE режиме после хода игрока AI должен автоматически сделать ответный ход:

```bash
# Создайте PvE игру
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pve"}'

# Сделайте ход игрока
curl -X POST http://localhost:8080/api/v1/games/{game-id}/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"row": 1, "col": 1}'

# Проверьте состояние - должен быть ход AI
curl -X GET http://localhost:8080/api/v1/games/{game-id} \
  -u alice:secret1
```

## Ручное тестирование

### 1. Регистрация пользователей

```bash
# Регистрация Alice
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}'

# Регистрация Bob
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}'
```

### 2. Создание PvP игры

```bash
# Alice создает игру
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pvp"}'

# Bob присоединяется к игре
curl -X POST http://localhost:8080/api/v1/games/{game-id}/join \
  -u bob:secret2
```

### 3. Выполнение ходов

```bash
# Bob делает первый ход
curl -X POST http://localhost:8080/api/v1/games/{game-id}/move \
  -H "Content-Type: application/json" \
  -u bob:secret2 \
  -d '{"row": 1, "col": 1}'

# Alice делает ответный ход
curl -X POST http://localhost:8080/api/v1/games/{game-id}/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"row": 0, "col": 0}'
```

### 4. Проверка состояния игры

```bash
# Получить текущее состояние
curl -X GET http://localhost:8080/api/v1/games/{game-id} \
  -u alice:secret1
```

## Проверка ошибок

### 1. Неверные данные

```bash
# Неверный режим игры
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "invalid"}'

# Неверные координаты хода
curl -X POST http://localhost:8080/api/v1/games/{game-id}/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"row": 5, "col": 5}'
```

### 2. Неавторизованный доступ

```bash
# Попытка доступа без аутентификации
curl -X GET http://localhost:8080/api/v1/games/{game-id}
```

### 3. Неверный ход

```bash
# Попытка хода не в свою очередь
curl -X POST http://localhost:8080/api/v1/games/{game-id}/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"row": 0, "col": 0}'
```

## Очистка данных

Для очистки базы данных между тестами:

```bash
# Очистить все данные
docker exec -it docker_files-db-1 psql -U game game_db -c "TRUNCATE users, games CASCADE;"
```

## Логи

Для просмотра логов приложения:

```bash
# Логи приложения
docker-compose logs app

# Логи базы данных
docker-compose logs db

# Все логи
docker-compose logs
```

## Остановка

```bash
# Остановить приложение
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```

## Устранение неполадок

### 1. Порт занят

Если порт 8080 занят, измените его в `docker-compose.yaml`:

```yaml
ports:
  - "8081:8080"  # Измените 8081 на свободный порт
```

### 2. База данных не запускается

```bash
# Проверьте статус контейнеров
docker-compose ps

# Перезапустите базу данных
docker-compose restart db

# Проверьте логи
docker-compose logs db
```

### 3. Приложение не отвечает

```bash
# Проверьте health check
curl http://localhost:8080/health

# Проверьте логи приложения
docker-compose logs app

# Перезапустите приложение
docker-compose restart app
```

## Результаты тестирования

При успешном прохождении всех тестов вы должны увидеть:

- ✅ Все API endpoints работают корректно
- ✅ PvP игры функционируют правильно
- ✅ PvE игры с AI работают автоматически
- ✅ Обработка ошибок работает корректно
- ✅ Аутентификация защищает ресурсы
- ✅ База данных сохраняет данные
- ✅ Счетчики побед обновляются

Проект готов к продакшену! 🎉
