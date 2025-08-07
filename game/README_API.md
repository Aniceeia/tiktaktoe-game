# Tic-Tac-Toe API Documentation

## Game Modes

- **PvP (Player vs Player)**: Two human players take turns
- **PvE (Player vs Environment)**: One human player vs AI algorithm

## Authentication

Currently uses simple header-based authentication:
- `X-User` header with username
- Or `?user=username` query parameter

## API Endpoints

### 1. Create Game
```bash
POST /games
Content-Type: application/json
X-User: alice

{
  "mode": "pvp"  # or "pve"
}
```

Response:
```json
{
  "id": "game-uuid",
  "board": [[0,0,0],[0,0,0],[0,0,0]],
  "status": "waiting",
  "player1": "alice",
  "player2": "",
  "next_turn": "alice",
  "mode": "pvp"
}
```

### 2. Join Game (PvP only)
```bash
POST /games/{game_id}/join
X-User: bob
```

### 3. Get Game State
```bash
GET /games/{game_id}
```

### 4. Make Move
```bash
POST /games/{game_id}/move
Content-Type: application/json
X-User: alice

{
  "board": [[1,0,0],[0,0,0],[0,0,0]]
}
```

## Game Flow Examples

### PvP Game Flow:
1. Alice creates game with mode "pvp"
2. Bob joins the game
3. Alice makes first move
4. Bob makes second move
5. Continue until game ends

### PvE Game Flow:
1. Alice creates game with mode "pve"
2. Alice makes move
3. AI automatically responds
4. Continue until game ends

## Board Representation
- `0`: Empty cell
- `1`: X (Player 1)
- `2`: O (Player 2 or AI)

## Game Status
- `waiting`: Game in progress
- `x_won`: X player won
- `o_won`: O player won  
- `draw`: Game ended in draw

## Error Responses
```json
{
  "error": "Error message"
}
```

Common errors:
- `"not your turn"`: Trying to move when it's not your turn
- `"game not found"`: Invalid game ID
- `"game is already completed"`: Game has ended
- `"invalid move detected"`: Invalid board state
- `"game is full"`: Trying to join a full PvP game 