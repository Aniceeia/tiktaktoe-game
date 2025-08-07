#!/bin/bash

# set -e

# echo "Testing Two-Player Game API"
# echo "============================"

# # Test 0: Register users
echo "0. Registering users..."
REG1=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login":"alice","password":"secret1"}')
echo "Register alice: $REG1"
REG2=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login":"bob","password":"secret2"}')
echo "Register bob: $REG2"
# REG3=$(curl -s -X POST http://localhost:8080/register \
#   -H "Content-Type: application/json" \
#   -d '{"login":"charlie","password":"secret3"}')
# echo "Register charlie: $REG3"
# echo ""

# # Test 0.1: Register duplicate user
# echo "0.1 Register duplicate user (should fail)..."
# REG_DUP=$(curl -s -X POST http://localhost:8080/register \
#   -H "Content-Type: application/json" \
#   -d '{"login":"alice","password":"secret1"}')
# echo "Register duplicate alice: $REG_DUP"
# echo ""

# # Test 0.2: Login with valid credentials
# echo "0.2 Login with valid credentials..."
# LOGIN1=$(curl -s -X POST http://localhost:8080/login \
#   -u alice:secret1)
# echo "Login alice: $LOGIN1"
# LOGIN2=$(curl -s -X POST http://localhost:8080/login \
#   -u bob:secret2)
# echo "Login bob: $LOGIN2"
# echo ""

# # Test 0.3: Login with invalid credentials
# echo "0.3 Login with invalid credentials (should fail)..."
# LOGIN_FAIL=$(curl -s -X POST http://localhost:8080/login \
#   -u alice:wrongpass)
# echo "Login fail: $LOGIN_FAIL"
# echo ""

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
echo "Login alice: $LOGIN1"

LOGIN2=$(curl -s -X POST http://localhost:8080/login \
  -u bob:secret2)
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


# Test 4: Make a move (Alice's turn)
echo "4. Making first move (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[0,0,0],[0,1,0],[0,0,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 4.1: Make a move (Bob's turn)
echo "4.1 Making first move (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u  bob:secret2 \
  -d '{"board":[[2,0,0],[0,1,0],[0,0,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 5: Make a move (Alice's turn)
echo "5. (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[2,0,1],[0,1,0],[0,0,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 5.1: Make a move (Bob's turn)
echo "5.1  (Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u  bob:secret2 \
  -d '{"board":[[2,0,1],[0,1,0],[2,0,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 6: Make a move (Alice's turn)
echo "6. (Alice)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u alice:secret1 \
  -d '{"board":[[2,0,1],[0,1,0],[2,1,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""

# Test 6.1: Make a move (Bob's turn)
echo 6.1  "(Bob)..."
MOVE_RESPONSE=$(curl -s -X POST http://localhost:8080/games/$GAME_ID/move \
  -H "Content-Type: application/json" \
  -u  bob:secret2 \
  -d '{"board":[[2,0,1],[2,1,0],[2,1,0]]}')
echo "Move response: $MOVE_RESPONSE"
echo ""
echo "Test completed!" 
