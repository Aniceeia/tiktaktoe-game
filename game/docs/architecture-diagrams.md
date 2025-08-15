# Схемы архитектуры проекта "Крестики-Нолики"

## 1. Общая архитектура системы

```mermaid
graph TB
    subgraph "Клиент"
        A[Web Browser]
        B[Mobile App]
        C[API Client]
    end
    
    subgraph "Load Balancer"
        LB[NGINX/HAProxy]
    end
    
    subgraph "Application Layer"
        APP1[Game Service 1]
        APP2[Game Service 2]
        APP3[Game Service N]
    end
    
    subgraph "Database Layer"
        DB[(PostgreSQL)]
        DB_REPLICA[(PostgreSQL Replica)]
    end
    
    subgraph "Monitoring"
        PROM[Prometheus]
        GRAF[Grafana]
        LOGS[ELK Stack]
    end
    
    A --> LB
    B --> LB
    C --> LB
    LB --> APP1
    LB --> APP2
    LB --> APP3
    APP1 --> DB
    APP2 --> DB
    APP3 --> DB
    DB --> DB_REPLICA
    APP1 --> PROM
    APP2 --> PROM
    APP3 --> PROM
    PROM --> GRAF
    APP1 --> LOGS
    APP2 --> LOGS
    APP3 --> LOGS
```

## 2. Clean Architecture - Слои

```mermaid
graph TD
    subgraph "External Layer"
        HTTP[HTTP Handlers]
        GRPC[gRPC Handlers]
        CLI[CLI Commands]
    end
    
    subgraph "Interface Layer"
        DTO[Data Transfer Objects]
        MID[Middleware]
        VAL[Validators]
    end
    
    subgraph "Application Layer"
        US[User Service]
        GS[Game Service]
        AS[Auth Service]
        AI[AI Service]
    end
    
    subgraph "Domain Layer"
        UE[User Entity]
        GE[Game Entity]
        UR[User Repository Interface]
        GR[Game Repository Interface]
        BL[Business Logic]
    end
    
    subgraph "Infrastructure Layer"
        UREP[User Repository Implementation]
        GREP[Game Repository Implementation]
        DB[(PostgreSQL)]
        CACHE[(Redis)]
    end
    
    HTTP --> DTO
    GRPC --> DTO
    CLI --> DTO
    DTO --> US
    DTO --> GS
    DTO --> AS
    US --> UE
    GS --> GE
    AS --> UE
    US --> UR
    GS --> GR
    UR --> UREP
    GR --> GREP
    UREP --> DB
    GREP --> DB
    GREP --> CACHE
```

## 3. Бизнес-процессы

### 3.1 Регистрация и авторизация

```mermaid
sequenceDiagram
    participant C as Клиент
    participant A as Auth Handler
    participant S as Auth Service
    participant U as User Service
    participant R as User Repository
    participant DB as Database
    
    C->>A: POST /auth/register
    A->>S: Register(login, password)
    S->>U: CreateUser(login, password)
    U->>R: Create(user)
    R->>DB: INSERT INTO users
    DB-->>R: User created
    R-->>U: User entity
    U-->>S: User entity
    S-->>A: User response
    A-->>C: 201 Created
```

### 3.2 Создание игры

```mermaid
sequenceDiagram
    participant C as Клиент
    participant G as Game Handler
    participant S as Game Service
    participant R as Game Repository
    participant DB as Database
    
    C->>G: POST /games {mode: "pvp"}
    G->>S: CreateGame(playerID, mode)
    S->>S: NewGame(creator, mode)
    S->>R: Create(game)
    R->>DB: INSERT INTO games
    DB-->>R: Game created
    R-->>S: Game entity
    S-->>G: Game response
    G-->>C: 201 Created
```

### 3.3 Игровой процесс

```mermaid
sequenceDiagram
    participant P1 as Player 1
    participant P2 as Player 2
    participant G as Game Handler
    participant S as Game Service
    participant AI as AI Service
    participant R as Game Repository
    participant DB as Database
    
    P1->>G: POST /games/{id}/move {row: 1, col: 1}
    G->>S: MakeMove(gameID, playerID, row, col)
    S->>S: Validate move
    S->>S: Update board
    S->>S: Check win condition
    S->>R: Update(game)
    R->>DB: UPDATE games
    DB-->>R: Game updated
    R-->>S: Game entity
    S-->>G: Game response
    G-->>P1: 200 OK
    
    alt PvE Mode
        S->>AI: FindBestMove(board)
        AI-->>S: Best move
        S->>S: Make AI move
        S->>R: Update(game)
        R->>DB: UPDATE games
        DB-->>R: Game updated
    end
```

## 4. Модель данных

```mermaid
erDiagram
    USERS {
        string uuid PK
        string login UK
        string password_hash
        int score
        timestamp created_at
        timestamp updated_at
    }
    
    GAMES {
        string id PK
        jsonb board
        string status
        string player1_id FK
        string player2_id
        string next_turn_id
        string mode
        timestamp created_at
        timestamp updated_at
    }
    
    USERS ||--o{ GAMES : "creates"
    USERS ||--o{ GAMES : "joins"
    GAMES }o--|| USERS : "player1"
    GAMES }o--|| USERS : "player2"
```

## 5. Компонентная архитектура

```mermaid
graph LR
    subgraph "API Gateway"
        GW[API Gateway]
    end
    
    subgraph "Microservices"
        AUTH[Auth Service]
        GAME[Game Service]
        USER[User Service]
    end
    
    subgraph "Data Layer"
        PG[(PostgreSQL)]
        REDIS[(Redis Cache)]
    end
    
    subgraph "External Services"
        AI[AI Service]
        NOTIFY[Notification Service]
    end
    
    GW --> AUTH
    GW --> GAME
    GW --> USER
    AUTH --> PG
    GAME --> PG
    USER --> PG
    GAME --> REDIS
    GAME --> AI
    GAME --> NOTIFY
```

## 6. Схема развертывания

```mermaid
graph TB
    subgraph "Production Environment"
        subgraph "Load Balancer"
            LB1[NGINX 1]
            LB2[NGINX 2]
        end
        
        subgraph "Application Servers"
            APP1[Game App 1]
            APP2[Game App 2]
            APP3[Game App 3]
        end
        
        subgraph "Database Cluster"
            DB_MASTER[(PostgreSQL Master)]
            DB_SLAVE1[(PostgreSQL Slave 1)]
            DB_SLAVE2[(PostgreSQL Slave 2)]
        end
        
        subgraph "Monitoring"
            PROM[Prometheus]
            GRAF[Grafana]
            ALERT[AlertManager]
        end
    end
    
    LB1 --> APP1
    LB1 --> APP2
    LB1 --> APP3
    LB2 --> APP1
    LB2 --> APP2
    LB2 --> APP3
    APP1 --> DB_MASTER
    APP2 --> DB_MASTER
    APP3 --> DB_MASTER
    DB_MASTER --> DB_SLAVE1
    DB_MASTER --> DB_SLAVE2
    APP1 --> PROM
    APP2 --> PROM
    APP3 --> PROM
    PROM --> GRAF
    PROM --> ALERT
```

## 7. Схема безопасности

```mermaid
graph TB
    subgraph "Client Layer"
        C[Client]
    end
    
    subgraph "Network Security"
        FW[Firewall]
        WAF[WAF]
        SSL[SSL/TLS]
    end
    
    subgraph "Application Security"
        AUTH[Authentication]
        AUTHZ[Authorization]
        VAL[Input Validation]
        RATE[Rate Limiting]
    end
    
    subgraph "Data Security"
        ENC[Encryption at Rest]
        BACKUP[Backup Encryption]
        AUDIT[Audit Logging]
    end
    
    C --> FW
    FW --> WAF
    WAF --> SSL
    SSL --> AUTH
    AUTH --> AUTHZ
    AUTHZ --> VAL
    VAL --> RATE
    RATE --> ENC
    ENC --> BACKUP
    BACKUP --> AUDIT
```

## 8. Схема масштабирования

```mermaid
graph TB
    subgraph "Auto Scaling Group"
        subgraph "Instance 1"
            APP1[Game App]
            MON1[Monitoring]
        end
        
        subgraph "Instance 2"
            APP2[Game App]
            MON2[Monitoring]
        end
        
        subgraph "Instance N"
            APPN[Game App]
            MONN[Monitoring]
        end
    end
    
    subgraph "Database Scaling"
        subgraph "Read Replicas"
            RR1[(Read Replica 1)]
            RR2[(Read Replica 2)]
            RRN[(Read Replica N)]
        end
        
        subgraph "Sharding"
            SHARD1[(Shard 1)]
            SHARD2[(Shard 2)]
            SHARDN[(Shard N)]
        end
    end
    
    APP1 --> RR1
    APP2 --> RR2
    APPN --> RRN
    RR1 --> SHARD1
    RR2 --> SHARD2
    RRN --> SHARDN
```
