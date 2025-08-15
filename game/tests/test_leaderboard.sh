#!/bin/bash
set -e

YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${YELLOW}[LEADERBOARD-TEST]${NC} $1"; }

log "Cleaning database..."
docker exec -it game-db psql -U game game_db -c "TRUNCATE users, games CASCADE;" >/dev/null

log "Register users..."
REG1=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"player1","password":"secret1"}')
echo "Register player1: $REG1"

REG2=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"player2","password":"secret2"}')
echo "Register player2: $REG2"

REG3=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"player3","password":"secret3"}')
echo "Register player3: $REG3"

log "Login users..."
LOGIN1=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"player1","password":"secret1"}')
AT1=$(echo "$LOGIN1" | jq -r '.accessToken')

LOGIN2=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"player2","password":"secret2"}')
AT2=$(echo "$LOGIN2" | jq -r '.accessToken')

LOGIN3=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"player3","password":"secret3"}')
AT3=$(echo "$LOGIN3" | jq -r '.accessToken')

log "Create and play games to generate stats..."

# Game 1: player1 wins (horizontal line)
GAME1=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"mode": "pvp"}')
GAME1_ID=$(echo "$GAME1" | jq -r '.id')

curl -s -X POST http://localhost:8080/api/v1/games/$GAME1_ID/join \
  -H "Authorization: Bearer $AT2" >/dev/null

# player1 wins with horizontal line
curl -s -X POST http://localhost:8080/api/v1/games/$GAME1_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"row": 0, "col": 0}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME1_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT2" \
  -d '{"row": 1, "col": 1}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME1_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"row": 0, "col": 1}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME1_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT2" \
  -d '{"row": 2, "col": 2}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME1_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"row": 0, "col": 2}' >/dev/null

log "Game 1 completed, player1 wins"

# Game 2: player2 wins (vertical line)
GAME2=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT2" \
  -d '{"mode": "pvp"}')
GAME2_ID=$(echo "$GAME2" | jq -r '.id')

curl -s -X POST http://localhost:8080/api/v1/games/$GAME2_ID/join \
  -H "Authorization: Bearer $AT3" >/dev/null

# player2 wins with vertical line
curl -s -X POST http://localhost:8080/api/v1/games/$GAME2_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT2" \
  -d '{"row": 0, "col": 0}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME2_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT3" \
  -d '{"row": 1, "col": 1}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME2_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT2" \
  -d '{"row": 1, "col": 0}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME2_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT3" \
  -d '{"row": 2, "col": 2}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME2_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT2" \
  -d '{"row": 2, "col": 0}' >/dev/null

log "Game 2 completed, player2 wins"

# Game 3: player1 wins again (diagonal)
GAME3=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"mode": "pvp"}')
GAME3_ID=$(echo "$GAME3" | jq -r '.id')

curl -s -X POST http://localhost:8080/api/v1/games/$GAME3_ID/join \
  -H "Authorization: Bearer $AT3" >/dev/null

# player1 wins with diagonal
curl -s -X POST http://localhost:8080/api/v1/games/$GAME3_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"row": 0, "col": 0}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME3_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT3" \
  -d '{"row": 0, "col": 1}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME3_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"row": 1, "col": 1}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME3_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT3" \
  -d '{"row": 0, "col": 2}' >/dev/null

curl -s -X POST http://localhost:8080/api/v1/games/$GAME3_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AT1" \
  -d '{"row": 2, "col": 2}' >/dev/null

log "Game 3 completed, player1 wins again"

log "Test leaderboard endpoint..."
LEADERBOARD=$(curl -s -X GET "http://localhost:8080/api/v1/leaderboard?limit=5" \
  -H "Authorization: Bearer $AT1")
echo "Leaderboard: $LEADERBOARD"

# Parse and verify leaderboard
PLAYER1_WINS=$(echo "$LEADERBOARD" | jq -r '.players[] | select(.login=="player1") | .wins')
PLAYER2_WINS=$(echo "$LEADERBOARD" | jq -r '.players[] | select(.login=="player2") | .wins')
PLAYER3_WINS=$(echo "$LEADERBOARD" | jq -r '.players[] | select(.login=="player3") | .wins')

echo "Player1 wins: $PLAYER1_WINS"
echo "Player2 wins: $PLAYER2_WINS" 
echo "Player3 wins: $PLAYER3_WINS"

if [ "$PLAYER1_WINS" != "2" ]; then
  echo "ERROR: Expected player1 to have 2 wins, got $PLAYER1_WINS"
  exit 1
fi

if [ "$PLAYER2_WINS" != "1" ]; then
  echo "ERROR: Expected player2 to have 1 win, got $PLAYER2_WINS"
  exit 1
fi

if [ "$PLAYER3_WINS" != "0" ]; then
  echo "ERROR: Expected player3 to have 0 wins, got $PLAYER3_WINS"
  exit 1
fi

log "Test leaderboard with different limit..."
LEADERBOARD2=$(curl -s -X GET "http://localhost:8080/api/v1/leaderboard?limit=2" \
  -H "Authorization: Bearer $AT1")
echo "Leaderboard (limit=2): $LEADERBOARD2"

TOTAL_PLAYERS=$(echo "$LEADERBOARD2" | jq -r '.total')
if [ "$TOTAL_PLAYERS" != "2" ]; then
  echo "ERROR: Expected 2 players with limit=2, got $TOTAL_PLAYERS"
  exit 1
fi

log "Test unauthorized access..."
UNAUTH=$(curl -s -o /dev/null -w "%{http_code}" -X GET "http://localhost:8080/api/v1/leaderboard")
if [ "$UNAUTH" != "401" ]; then
  echo "Expected 401 for unauthorized access, got $UNAUTH"
  exit 1
fi

log "Leaderboard tests completed successfully!"
