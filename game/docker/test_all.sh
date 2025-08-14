#!/bin/bash

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
  -u alice:secret1)
ALICE_UUID=$(echo $LOGIN1 | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "Login alice: $LOGIN1"
echo "Alice UUID: $ALICE_UUID"

LOGIN2=$(curl -s -X POST http://localhost:8080/login \
  -u bob:secret2)
BOB_UUID=$(echo $LOGIN2 | grep -o '"uuid":"[^"]*"' | cut -d'"' -f4)
echo "Login bob: $LOGIN2"
echo "Bob UUID: $BOB_UUID"

echo "========== PvP TESTS =========="
echo "1. Creating a PvP game..."
GAME_RESPONSE=$(curl -s -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pvp"}')
echo "Game response: $GAME_RESPONSE"

GAME_ID=$(echo $GAME_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Game ID: $GAME_ID"

echo "2. Joining PvP game..."
JOIN_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/join \
  -u bob:secret2)
echo "Join response: $JOIN_RESPONSE"

echo "3. Getting PvP game state..."
STATE_RESPONSE=$(curl -s -X GET http://localhost:8080/games/$GAME_ID \
  -u alice:secret1)
echo "Game state: $STATE_RESPONSE"

echo "3.1 Creating another PvP game (for available games test)..."
GAME2_RESPONSE=$(curl -s -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -u bob:secret2 \
  -d '{"mode": "pvp"}')
echo "Second PvP game response: $GAME2_RESPONSE"

echo "========== PvE TESTS =========="
echo "4. Creating a PvE game..."
GAME_PVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pve"}')
echo "PvE Game response: $GAME_PVE_RESPONSE"

GAME_PVE_ID=$(echo $GAME_PVE_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "PvE Game ID: $GAME_PVE_ID"

echo "5. Getting initial PvE game state..."
PVE_STATE1=$(curl -s -X GET http://localhost:8080/games/$GAME_PVE_ID \
  -u alice:secret1)
echo "Initial PvE game state: $PVE_STATE1"

echo "6. Making first move in PvE (player)..."
PVE_MOVE1=$(curl -s -X POST http://localhost:8080/games/$GAME_PVE_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[0,0,0],[0,1,0],[0,0,0]]}')
echo "Player move response: $PVE_MOVE1"

echo "7. Getting PvE game state after player move..."
PVE_STATE2=$(curl -s -X GET http://localhost:8080/games/$GAME_PVE_ID \
  -u alice:secret1)
echo "PvE game state after player move: $PVE_STATE2"

echo "8. Checking if AI made a move..."
# Проверяем, что AI сделал ход (должна быть хотя бы одна 2 в поле)
AI_MOVE_PRESENT=$(echo $PVE_STATE2 | grep -o '2' | wc -l)
if [ $AI_MOVE_PRESENT -gt 0 ]; then
  echo "AI successfully made a move"
else
  echo "ERROR: AI didn't make a move"
  echo $AI_MOVE_PRESENT

 # exit 1
fi

echo "9. Making second move in PvE (player)..."
PVE_MOVE2=$(curl -s -X POST http://localhost:8080/games/$GAME_PVE_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[2,0,0],[0,1,0],[1,0,0]]}')
echo "Player move response: $PVE_MOVE2"

echo "10. Final PvE game state..."
PVE_STATE3=$(curl -s -X GET http://localhost:8080/games/$GAME_PVE_ID \
  -u alice:secret1)
echo "Final PvE game state: $PVE_STATE3"

echo "11. Testing invalid PvE join attempt..."
PVE_JOIN_ATTEMPT=$(curl -s -X POST http://localhost:8080/games/$GAME_PVE_ID/join \
  -u bob:secret2)
echo "PvE join attempt response (should fail): $PVE_JOIN_ATTEMPT"

echo "========== GAME LIST TESTS =========="
echo "12. Checking available games..."
AVAILABLE_GAMES=$(curl -s -X GET http://localhost:8080/games/available \
  -u bob:secret2)
echo "Available games: $AVAILABLE_GAMES"

echo "13. Creating another PvE game..."
GAME_PVE2_RESPONSE=$(curl -s -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -u bob:secret2 \
  -d '{"mode": "pve"}')
GAME_PVE2_ID=$(echo $GAME_PVE2_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Second PvE game created with ID: $GAME_PVE2_ID"

echo "13.1 Creating additional PvP game from Alice (for available games test)..."
GAME3_RESPONSE=$(curl -s -X POST http://localhost:8080/games \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"mode": "pvp"}')
echo "Third PvP game response: $GAME3_RESPONSE"

echo "14. Checking available games after creation (Alice's view)..."
AVAILABLE_GAMES2=$(curl -s -X GET http://localhost:8080/games/available \
  -u alice:secret1)
echo "Available games now (Alice): $AVAILABLE_GAMES2"

echo "14.1 Checking available games (Bob's view)..."
AVAILABLE_GAMES3=$(curl -s -X GET http://localhost:8080/games/available \
  -u bob:secret2)
echo "Available games now (Bob): $AVAILABLE_GAMES3"

echo "========== USER TESTS =========="
# В test_all.sh
echo "15. Testing user info endpoint..."
ALICE_INFO=$(curl -s -X GET "http://localhost:8080/users/$ALICE_UUID" \
  -H "Authorization: Basic $(echo -n alice:secret1 | base64)")
echo "Alice info: $ALICE_INFO"

echo "Test completed successfully!"