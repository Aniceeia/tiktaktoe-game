#!/bin/bash

echo "0. Registering users..."
REG1=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}')
echo "Register alice: $REG1"
REG2=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}')
echo "Register bob: $REG2"

set -e

echo "Cleaning database..."
docker exec -it docker_files-db-1 psql -U game game_db -c "TRUNCATE users, games CASCADE;"

echo "0. Registering users..."
REG1=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}')
echo "Register alice: $REG1"

REG2=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}')
echo "Register bob: $REG2"

echo "0.1 Logging in users..."
LOGIN1=$(curl -s -X POST http://localhost:8080/login \
  -u alice:secret1 | jq -r '.uuid')
echo "Login alice: $LOGIN1"

LOGIN2=$(curl -s -X POST http://localhost:8080/login \
  -u bob:secret2 | jq -r '.uuid')
echo "Login bob: $LOGIN2"

echo "1. Creating a PvP game..."
GAME_RESPONSE=$(curl -s -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pvp"}')
echo "Game response: $GAME_RESPONSE"

GAME_ID=$(echo $GAME_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Game ID: $GAME_ID"

echo "2. Joining game..."
JOIN_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/join \
  -u bob:secret2)
echo "Join response: $JOIN_RESPONSE"

# Test 3: Get game state
echo "3. Getting game state..."
STATE_RESPONSE=$(curl -s -X GET http://localhost:8080/games/$GAME_ID \
  -u alice:secret1)
echo "Game state: $STATE_RESPONSE"
echo ""

# Test 4: Make a move (Bob's turn - he joined first)
echo "4. Making first move (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u bob:secret2 \
  -d '{"board":[[0,0,0],[0,2,0],[0,0,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 4.1: Make a move (Alice's turn)
echo "4.1 Making first move (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[1,0,0],[0,2,0],[0,0,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 5: Make a move (Bob's turn)
echo "5. (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u bob:secret2 \
  -d '{"board":[[1,0,0],[0,2,0],[0,0,2]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 5.1: Make a move (Alice's turn)
echo "5.1 (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[1,0,0],[1,2,0],[0,0,2]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 6: Make a move (Bob's turn)
echo "6. (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u bob:secret2 \
  -d '{"board":[[1,0,0],[1,2,0],[0,2,2]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 6.1: Make a move (Alice's turn)
echo "6.1 (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[1,0,1],[1,2,0],[0,2,2]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 6.2: Make a move (Bob's turn) - Bob wins
echo "6.2 (Bob wins)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u bob:secret2 \
  -d '{"board":[[1,0,1],[1,2,0],[2,2,2]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 7: get user
echo "7. Testing user info endpoint..."
ALICE_INFO=$(curl -s -X GET http://localhost:8080/users/$LOGIN1 \
  -u alice:secret1)
echo "Alice info: $ALICE_INFO"
echo "Alice id: $LOGIN1"

BOB_INFO=$(curl -s -X GET http://localhost:8080/users/$LOGIN2 \
  -u bob:secret2)
echo "Bob info: $BOB_INFO"
echo "Bob id: $LOGIN2"

echo "Test completed!" 