#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_header() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}



print_header "Cleaning database..."
docker exec -it game-db psql -U game game_db -c "TRUNCATE users, games CASCADE;"

print_header "0. Registering users..."
REG1=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}')
echo "Register alice: $REG1"

REG2=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}')
echo "Register bob: $REG2"

print_header "0.1 Logging in users (JWT)..."
LOGIN1=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}')
ALICE_UUID=$(echo $LOGIN1 | jq -r '.uuid // ""')
ALICE_AT=$(echo $LOGIN1 | jq -r '.accessToken')
ALICE_RT=$(echo $LOGIN1 | jq -r '.refreshToken')
echo "Login alice: $LOGIN1"
echo "Alice UUID: $ALICE_UUID"

LOGIN2=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}')
BOB_UUID=$(echo $LOGIN2 | jq -r '.uuid // ""')
BOB_AT=$(echo $LOGIN2 | jq -r '.accessToken')
BOB_RT=$(echo $LOGIN2 | jq -r '.refreshToken')
echo "Login bob: $LOGIN2"
echo "Bob UUID: $BOB_UUID"

print_header "========== PvP TESTS =========="
print_header "1. Creating a PvP game..."
GAME_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"mode": "pvp"}')
echo "Game response: $GAME_RESPONSE"

GAME_ID=$(echo $GAME_RESPONSE | jq -r '.id')
echo "Game ID: $GAME_ID"

print_header "2. Joining PvP game..."
JOIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/join \
  -H "Authorization: Bearer $BOB_AT")
echo "Join response: $JOIN_RESPONSE"

print_header "3. Getting PvP game state..."
STATE_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/games/$GAME_ID \
  -H "Authorization: Bearer $ALICE_AT")
echo "Game state: $STATE_RESPONSE"

print_header "3.1 Creating another PvP game (for available games test)..."
GAME2_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_AT" \
  -d '{"mode": "pvp"}')
echo "Second PvP game response: $GAME2_RESPONSE"

print_header "========== PvE TESTS =========="
print_header "4. Creating a PvE game..."
GAME_PVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"mode": "pve"}')
echo "PvE Game response: $GAME_PVE_RESPONSE"

GAME_PVE_ID=$(echo $GAME_PVE_RESPONSE | jq -r '.id')
echo "PvE Game ID: $GAME_PVE_ID"

print_header "5. Getting initial PvE game state..."
PVE_STATE1=$(curl -s -X GET http://localhost:8080/api/v1/games/$GAME_PVE_ID \
  -H "Authorization: Bearer $ALICE_AT")
echo "Initial PvE game state: $PVE_STATE1"

print_header "6. Making first move in PvE (player)..."
PVE_MOVE1=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_PVE_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"row": 1, "col": 1}')
echo "Player move response: $PVE_MOVE1"

print_header "7. Getting PvE game state after player move..."
PVE_STATE2=$(curl -s -X GET http://localhost:8080/api/v1/games/$GAME_PVE_ID \
  -H "Authorization: Bearer $ALICE_AT")
echo "PvE game state after player move: $PVE_STATE2"

print_header "8. Checking if AI made a move..."
# Проверяем, что AI сделал ход (должна быть хотя бы одна 2 в поле)
AI_MOVE_PRESENT=$(echo $PVE_STATE2 | grep -o '2' | wc -l)
if [ $AI_MOVE_PRESENT -gt 0 ]; then
  echo "AI successfully made a move"
else
  echo "ERROR: AI didn't make a move"
  echo $AI_MOVE_PRESENT
  exit 1
fi

print_header "9. Making second move in PvE (player)..."
PVE_MOVE2=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_PVE_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"row": 0, "col": 0}')
echo "Player move response: $PVE_MOVE2"

print_header "10. Final PvE game state..."
PVE_STATE3=$(curl -s -X GET http://localhost:8080/api/v1/games/$GAME_PVE_ID \
  -H "Authorization: Bearer $ALICE_AT")
echo "Final PvE game state: $PVE_STATE3"

print_header "11. Testing invalid PvE join attempt..."
PVE_JOIN_ATTEMPT=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_PVE_ID/join \
  -H "Authorization: Bearer $BOB_AT")
echo "PvE join attempt response (should fail): $PVE_JOIN_ATTEMPT"

print_header "========== GAME LIST TESTS =========="
print_header "12. Checking available games..."
AVAILABLE_GAMES=$(curl -s -X GET http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer $BOB_AT")
echo "Available games: $AVAILABLE_GAMES"

print_header "13. Creating another PvE game..."
GAME_PVE2_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_AT" \
  -d '{"mode": "pve"}')
GAME_PVE2_ID=$(echo $GAME_PVE2_RESPONSE | jq -r '.id')
echo "Second PvE game created with ID: $GAME_PVE2_ID"

print_header "13.1 Creating additional PvP game from Alice (for available games test)..."
GAME3_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"mode": "pvp"}')
echo "Third PvP game response: $GAME3_RESPONSE"

print_header "14. Checking available games after creation (Alice's view)..."
AVAILABLE_GAMES2=$(curl -s -X GET http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer $ALICE_AT")
echo "Available games now (Alice): $AVAILABLE_GAMES2"

print_header "14.1 Checking available games (Bob's view)..."
AVAILABLE_GAMES3=$(curl -s -X GET http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer $BOB_AT")
echo "Available games now (Bob): $AVAILABLE_GAMES3"

print_header "15. Checking user games..."
USER_GAMES_ALICE=$(curl -s -X GET http://localhost:8080/api/v1/games/my \
  -H "Authorization: Bearer $ALICE_AT")
echo "Alice's games: $USER_GAMES_ALICE"

USER_GAMES_BOB=$(curl -s -X GET http://localhost:8080/api/v1/games/my \
  -H "Authorization: Bearer $BOB_AT")
echo "Bob's games: $USER_GAMES_BOB"

print_header "========== USER TESTS =========="
print_header "16. Testing user info endpoint..."
ALICE_INFO=$(curl -s -X GET "http://localhost:8080/api/v1/users/$ALICE_UUID" \
  -H "Authorization: Bearer $ALICE_AT")
echo "Alice info: $ALICE_INFO"

BOB_INFO=$(curl -s -X GET "http://localhost:8080/api/v1/users/$BOB_UUID" \
  -H "Authorization: Bearer $BOB_AT")
echo "Bob info: $BOB_INFO"

print_header "========== ERROR TESTS =========="
print_header "17. Testing invalid move..."
INVALID_MOVE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"row": 1, "col": 1}')
echo "Invalid move response: $INVALID_MOVE"

print_header "18. Testing wrong turn..."
WRONG_TURN=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"row": 0, "col": 0}')
echo "Wrong turn response: $WRONG_TURN"

print_header "19. Testing unauthorized access..."
UNAUTHORIZED=$(curl -s -X GET http://localhost:8080/api/v1/games/$GAME_ID)
echo "Unauthorized response: $UNAUTHORIZED"

print_header "20. Testing health check..."
HEALTH=$(curl -s -X GET http://localhost:8080/health)
echo "Health check: $HEALTH"
