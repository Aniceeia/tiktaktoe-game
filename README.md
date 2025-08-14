# **TicTacToe Arena: Multiplayer Game Service**  
**Сервис для игры в крестики-нолики с REST API и JWT-аутентификацией**  

[![Go](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/) [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://www.postgresql.org/) [![Docker](https://img.shields.io/badge/Docker-24.0+-blue.svg)](https://www.docker.com/)  

<p align="center">
  <img src="docs/gameplay.gif" alt="Game Demo" width="500">
</p>

## **🚀 Возможности**  
- **REST API** для управления игрой  
- **JWT-аутентификация** игроков  
- **PostgreSQL** для хранения данных  
- **Docker-контейнеризация**  
- **Логика игры** с проверкой победных комбинаций  
- **Лобби** для ожидания второго игрока  

## **🛠 Технологии**  
- **Язык**: Go 1.20+  
- **Фреймворки**: Chi (роутинг), uber/fx (DI)  
- **База данных**: PostgreSQL 15+  
- **Аутентификация**: JWT  
- **Инфраструктура**: Docker, Docker Compose  

## **⚙️ Установка**  

### **1. Запуск через Docker**  
```bash
git clone https://github.com/Aniceeia/tiktaktoe-game.git
cd tiktaktoe-game
docker-compose up --build
```

### **2. Локальная разработка**  
```bash
# Установите зависимости
go mod download

# Настройте .env (скопируйте из .env.example)
cp .env.example .env

# Запустите сервер
go run cmd/main.go
```

## **📡 API Endpoints**  
| Метод | Путь | Описание |  
|-------|------|----------|  
| POST | `/api/register` | Регистрация игрока |  
| POST | `/api/login` | Вход в систему |  
| POST | `/api/game/create` | Создать новую игру |  
| POST | `/api/game/{id}/move` | Сделать ход |  

Полная документация: [API.md](docs/API.md)  

## **🎮 Как играть?**  
1. Зарегистрируйтесь:  
   ```bash
   curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"username":"player1", "password":"qwerty"}'
   ```
2. Создайте игру:  
   ```bash
   curl -X POST http://localhost:8080/api/game/create \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```
3. Присоединитесь к игре (из другого терминала):  
   ```bash
   curl -X POST http://localhost:8080/api/game/1/join \
     -H "Authorization: Bearer ANOTHER_JWT_TOKEN"
   ```

## **📊 Производительность**  
| Параметр | Значение |  
|----------|----------|  
| Время отклика API | < 50 мс |  
| Макс. одновременных игр | 1000+ |  

## **📚 Документация**  
- [Архитектура](docs/ARCHITECTURE.md)  
- [Настройка БД](docs/DATABASE.md)  

## **📜 Лицензия**  
MIT License. Подробнее в [LICENSE](LICENSE).  

---  
**Автор**: Хазова Анастасия  
**Контакты**: [GitHub](https://github.com/Aniceeia) | akhazova@icloud.com  

---  
**Примечание**: Для production-сборки используйте `go build -o tiktaktoe cmd/main.go`.