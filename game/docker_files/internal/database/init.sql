CREATE TABLE users (
    uuid VARCHAR(36) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    score INT DEFAULT 0
);

CREATE TABLE games (
    id VARCHAR(36) PRIMARY KEY,
    board JSONB NOT NULL,
    status VARCHAR(50) NOT NULL,
    player1 VARCHAR(36) REFERENCES users(uuid) ON DELETE SET NULL,
    player2 VARCHAR(36),
    next_turn VARCHAR(36),
    mode VARCHAR(10) NOT NULL
);