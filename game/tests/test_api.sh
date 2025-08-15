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


print_header "0. Registering users..."
REG1=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}')
print_header "Register alice: $REG1"
REG2=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}')
echo "Register bob: $REG2"

set -e

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
ALICE_AT=$(echo $LOGIN1 | jq -r '.accessToken')
ALICE_RT=$(echo $LOGIN1 | jq -r '.refreshToken')
echo "Login alice: $LOGIN1"

LOGIN2=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}')
BOB_AT=$(echo $LOGIN2 | jq -r '.accessToken')
BOB_RT=$(echo $LOGIN2 | jq -r '.refreshToken')
echo "Login bob: $LOGIN2"

print_header "1. Creating a PvP game..."
GAME_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"mode": "pvp"}')
echo "Game response: $GAME_RESPONSE"

GAME_ID=$(echo $GAME_RESPONSE | jq -r '.id')
echo "Game ID: $GAME_ID"

print_header "2. Joining game..."
JOIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/join \
  -H "Authorization: Bearer $BOB_AT")
echo "Join response: $JOIN_RESPONSE"

print_header "3. Getting game state..."
STATE_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/games/$GAME_ID \
  -H "Authorization: Bearer $ALICE_AT")
echo "Game state: $STATE_RESPONSE"

echo ""

print_header "4. Making first move (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_AT" \
  -d '{"row": 1, "col": 1}')
echo "Move response: $MOVE_RESPONSE"

echo ""

print_header "4.1 Making first move (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"row": 0, "col": 0}')
echo "Move response: $MOVE_RESPONSE"

echo ""

print_header "5. (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_AT" \
  -d '{"row": 2, "col": 2}')
echo "Move response: $MOVE_RESPONSE"

echo ""

print_header "5.1 (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"row": 1, "col": 0}')
echo "Move response: $MOVE_RESPONSE"

echo ""

print_header "6. (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_AT" \
  -d '{"row": 2, "col": 1}')
echo "Move response: $MOVE_RESPONSE"

echo ""

print_header "6.1 (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_AT" \
  -d '{"row": 0, "col": 2}')
echo "Move response: $MOVE_RESPONSE"

echo ""

print_header "6.2 (Bob wins)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_AT" \
  -d '{"row": 2, "col": 0}')
echo "Move response: $MOVE_RESPONSE"

echo ""

print_header "7. Testing user info endpoint..."
ALICE_INFO=$(curl -s -X GET http://localhost:8080/api/v1/users/$(echo $LOGIN1 | jq -r '.uuid') \
  -H "Authorization: Bearer $ALICE_AT")
echo "Alice info: $ALICE_INFO"

BOB_INFO=$(curl -s -X GET http://localhost:8080/api/v1/users/$(echo $LOGIN2 | jq -r '.uuid') \
  -H "Authorization: Bearer $BOB_AT")
echo "Bob info: $BOB_INFO"
