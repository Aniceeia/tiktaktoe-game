#!/bin/bash
set -e

YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${YELLOW}[JWT-TEST]${NC} $1"; }

log "Cleaning database..."
docker exec -it game-db psql -U game game_db -c "TRUNCATE users, games CASCADE;" >/dev/null

log "Register user..."
REG=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"jwt_user","password":"secretpw"}')
echo "Register: $REG"

log "Login (get tokens)..."
LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"jwt_user","password":"secretpw"}')
echo "Login: $LOGIN"
AT=$(echo "$LOGIN" | jq -r '.accessToken')
RT=$(echo "$LOGIN" | jq -r '.refreshToken')

if [ -z "$AT" ] || [ "$AT" = "null" ]; then
  echo "accessToken not received"
  exit 1
fi

log "Access protected endpoint with accessToken"
ME=$(curl -s -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $AT")
echo "Me: $ME"

log "Refresh access token"
NEW_ACCESS=$(curl -s -X POST http://localhost:8080/api/v1/auth/refresh/access \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\":\"$RT\"}")
echo "New access: $NEW_ACCESS"
NAT=$(echo "$NEW_ACCESS" | jq -r '.accessToken')

if [ -z "$NAT" ] || [ "$NAT" = "null" ]; then
  echo "new accessToken not received"
  exit 1
fi

log "Rotate refresh token"
ROTATED=$(curl -s -X POST http://localhost:8080/api/v1/auth/refresh/rotate \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\":\"$RT\"}")
echo "Rotated refresh: $ROTATED"
NRT=$(echo "$ROTATED" | jq -r '.refreshToken')

if [ -z "$NRT" ] || [ "$NRT" = "null" ]; then
  echo "new refreshToken not received"
  exit 1
fi

log "Check unauthorized with invalid token"
UNAUTH=$(curl -s -o /dev/null -w "%{http_code}" -X GET http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer invalid.token")
if [ "$UNAUTH" != "401" ]; then
  echo "Expected 401 for invalid token, got $UNAUTH"
  exit 1
fi

log "JWT tests completed successfully!"
