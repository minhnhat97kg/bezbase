#!/bin/bash

# Test script for RBAC API endpoints
BASE_URL="http://localhost:8080"
TOKEN=""

echo "=== RBAC API Test Script ==="
echo ""

# Function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=${3:-}
    
    echo "Testing: $method $endpoint"
    if [ -n "$data" ]; then
        curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data" | jq '.' 2>/dev/null || echo "Request failed"
    else
        curl -s -X $method "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $TOKEN" | jq '.' 2>/dev/null || echo "Request failed"
    fi
    echo ""
}

# Test health check first
echo "1. Health Check"
curl -s "$BASE_URL/api/health" | jq '.' 2>/dev/null || echo "Health check failed"
echo ""

# Test user registration
echo "2. Register Test User"
REGISTER_DATA='{
    "email": "admin@test.com",
    "password": "password123",
    "first_name": "Admin",
    "last_name": "User"
}'
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
    -H "Content-Type: application/json" \
    -d "$REGISTER_DATA")
echo "$REGISTER_RESPONSE" | jq '.' 2>/dev/null || echo "Registration failed"
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.token' 2>/dev/null)
echo "Token: $TOKEN"
echo ""

# Test login
echo "3. Login Test User"
LOGIN_DATA='{
    "username": "admin@test.com",
    "password": "password123"
}'
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA")
echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "Login failed"
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token' 2>/dev/null)
echo "Updated Token: $TOKEN"
echo ""

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to get authentication token. Cannot test RBAC endpoints."
    exit 1
fi

# Test RBAC endpoints
echo "4. Get All Roles"
api_call "GET" "/api/v1/rbac/roles"

echo "5. Create New Role"
ROLE_DATA='{
    "name": "editor",
    "display_name": "Content Editor",
    "description": "Can create and edit content",
    "is_active": true
}'
api_call "POST" "/api/v1/rbac/roles" "$ROLE_DATA"

echo "6. Get Role by ID"
api_call "GET" "/api/v1/rbac/roles/1"

echo "7. Update Role"
UPDATE_DATA='{
    "display_name": "Senior Editor",
    "description": "Senior content editor with additional permissions"
}'
api_call "PUT" "/api/v1/rbac/roles/1" "$UPDATE_DATA"

echo "8. Assign Admin Role to User"
ASSIGN_DATA='{
    "user_id": 1,
    "role": "admin"
}'
api_call "POST" "/api/v1/rbac/users/assign-role" "$ASSIGN_DATA"

echo "9. Get User Roles"
api_call "GET" "/api/v1/rbac/users/1/roles"

echo "10. Add Permission to Role"
PERMISSION_DATA='{
    "role": "editor",
    "resource": "posts",
    "action": "create"
}'
api_call "POST" "/api/v1/rbac/permissions" "$PERMISSION_DATA"

echo "11. Get Role Permissions"
api_call "GET" "/api/v1/rbac/roles/editor/permissions"

echo "12. Check User Permission"
api_call "GET" "/api/v1/rbac/users/1/check-permission?resource=posts&action=create"

echo "Test completed!"