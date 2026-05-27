#!/usr/bin/env zsh
# Requires zsh

set -e

content_type_header="Content-Type: application/json"
regular_user_body='{"email":"user@example.com", "password":"password123"}'
admin_user_body='{"email":"admin@example.com", "password":"admin123"}'
base_url='http://localhost:8080/api'
auth_header='Authorization: Bearer'
x_admin_header='X-Admin-Key:'

test_count=1

# Positive: login as regular user
res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$regular_user_body" \
    --location "$base_url/auth/login")

token=$(echo "$res" | jq -r '.token')

if [ -n "$token" ] && [ "$token" != "null" ]; then
  echo "TEST $test_count: ✅ PASS - Login as a regular user successful"
else
  echo "TEST $test_count: ❌ FAIL - Login as a regular user failed"
fi

((++test_count))

# Positive: checking current user
res=$(curl -X GET -sS \
    -H "$content_type_header" \
    --location "$base_url/users/me" \
    -H "$auth_header $token")

id=$(echo "$res" | jq -r '.id')
name=$(echo "$res" | jq -r '.name')
role=$(echo "$res" | jq -r '.role')
bio=$(echo "$res" | jq -r '.bio')
err=$(echo "$res" | jq -r '.error')


if [[ -n "$id" && "$id" != "null" &&
      -n "$name" && "$name" != "null" &&
      -n "$role" && "$role" != "null" &&
      "$role" == "user" &&
      -n "$bio" && "$bio" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - /users/me as regular user success"
elif [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - /users/me as regular user failed"

fi

((++test_count))

# Positive: changing regular user's bio

rand=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 16)

newBio="This is a regular user new bio with no XSS $rand"

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$newBio"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')
err=$(echo "$res" | jq -r '.error')

if [[ "$bio" == "$newBio" ]]; then
  echo "TEST $test_count: ✅ PASS - Bio of regular user changed successfully!"
elif [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - Bio of regular user not changed"
fi

((++test_count))

# Test XSS-001: Script tag in bio (from xss-tests.md)
xss_payload='<script>alert("XSS")</script>'
body=$(jq -nc --arg bio "$xss_payload" '{bio:$bio}')

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d "$body" \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')
err=$(echo "$res" | jq -r '.error')
if [[ "$bio" != *"<script>"* ]] && [[ "$bio" != *"alert"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-001: Script tag sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-001: Script tag NOT sanitized"
fi

((++test_count))

# Test XSS-002: Image tag with onerror (from xss-tests.md)
xss_payload='<img src=x onerror=alert("XSS")>'

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$xss_payload"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')

if [[ "$bio" != *"onerror"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-002: Image onerror sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-002: Image onerror NOT sanitized"
fi

((++test_count))

# Test XSS-003: SVG tag (from xss-tests.md)
xss_payload='<svg onload=alert("XSS")></svg>'

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$xss_payload"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')

if [[ "$bio" != *"onload"* ]] && [[ "$bio" != *"<svg"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-003: SVG sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-003: SVG NOT sanitized"
fi

((++test_count))

# Test XSS-006: JavaScript protocol (from xss-tests.md)
xss_payload='<a href="javascript:alert('"'"'XSS'"'"')">Click</a>'

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$xss_payload"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')

if [[ "$bio" != *"javascript:"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-006: JavaScript protocol sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-006: JavaScript protocol NOT sanitized"
fi

((++test_count))

# Test XSS-005a: Body onload (from xss-tests.md)
xss_payload='<body onload=alert("XSS")>'

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$xss_payload"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')

if [[ "$bio" != *"onload"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-005a: Body onload sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-005a: Body onload NOT sanitized"
fi

((++test_count))

# Test XSS-005b: Input onfocus (from xss-tests.md)
xss_payload='<input onfocus=alert("XSS") autofocus>'

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$xss_payload"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')

if [[ "$bio" != *"onfocus"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-005b: Input onfocus sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-005b: Input onfocus NOT sanitized"
fi

((++test_count))

# Test XSS-005c: Marquee onstart (from xss-tests.md)
xss_payload='<marquee onstart=alert("XSS")>'

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$xss_payload"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')

if [[ "$bio" != *"onstart"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-005c: Marquee onstart sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-005c: Marquee onstart NOT sanitized"
fi

((++test_count))

# Test XSS-005d: Div onmouseover (from xss-tests.md)
xss_payload='<div onmouseover=alert("XSS")>Hover me</div>'

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$xss_payload"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $token")

bio=$(echo "$res" | jq -r '.bio')

if [[ "$bio" != *"onmouseover"* ]]; then
  echo "TEST $test_count: ✅ PASS - XSS-005d: Div onmouseover sanitized"
else
  echo "TEST $test_count: ❌ FAIL - XSS-005d: Div onmouseover NOT sanitized"
fi

((++test_count))

# Negative: changing other users' bio

adminId=1

attemptAdminBio="Wuzzah"

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$attemptAdminBio"'"}' \
    --location "$base_url/users/$adminId/profile" \
    -H "$auth_header $token")

err=$(echo "$res" | jq -r '.error')
bio=$(echo "$res" | jq -r '.bio')

if [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - (Regular) user can't change other's bio"
elif [[ "$bio" == "$attemptAdminBio" ]]; then
  echo "TEST $test_count: ❌ FAIL - (Regular) user can change other's bio"
  echo "$bio"
fi

((++test_count))

# Negative: accessing admin endpoint as regular user

res=$(curl -X GET -sS \
    -H "$content_type_header" \
    --location "$base_url/admin/users" \
    -H "Authorization: Bearer $token")

err=$(echo "$res" | jq -r '.error')
users=$(echo "$res" | jq -r '.users')

if [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - Regular user cannot access admin endpoint"
elif [[ -n "$users" && "$users" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - Regular user can access admin endpoint"
fi

((++test_count))

# Negative: accessing admin endpoint as regular user with admin key

adminKey=''

res=$(curl -X GET -sS \
    -H "$content_type_header" \
    --location "$base_url/admin/users" \
    -H "$auth_header $token" \
    -H "$x_admin_header $adminKey")

echo $res

err=$(echo "$res" | jq -r '.error')
users=$(echo "$res" | jq -r '.users')

if [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - Regular user cannot access admin endpoint even with admin key"
elif [[ -n "$users" && "$users" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - Regular user can access admin endpoint with admin key"
fi

((++test_count))

# Positive: login as admin

res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$admin_user_body" \
    --location "$base_url/auth/login")

adminToken=$(echo "$res" | jq -r '.token')

if [ -n "$adminToken" ] && [ "$adminToken" != "null" ]; then
  echo "TEST $test_count: ✅ PASS - Login as admin successful"
else
  echo "TEST $test_count: ❌ FAIL - Login as admin failed"
fi

((++test_count))

# Positive: checking current admin

res=$(curl -X GET -sS \
    -H "$content_type_header" \
    --location "$base_url/users/me" \
    -H "$auth_header $adminToken")

adminId=$(echo "$res" | jq -r '.id')
name=$(echo "$res" | jq -r '.name')
role=$(echo "$res" | jq -r '.role')
bio=$(echo "$res" | jq -r '.bio')
err=$(echo "$res" | jq -r '.error')


if [[ -n "$adminId" && "$adminId" != "null" &&
      -n "$name" && "$name" != "null" &&
      -n "$role" && "$role" != "null" &&
      "$role" == "admin" &&
      -n "$bio" && "$bio" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - /users/me as admin success"
elif [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - /users/me as admin failed"
fi

# Positive: changing admin's own bio

newBio="This is an admin user new bio with no XSS $rand"

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$newBio"'"}' \
    --location "$base_url/users/$adminId/profile" \
    -H "$auth_header $adminToken")

bio=$(echo "$res" | jq -r '.bio')
err=$(echo "$res" | jq -r '.error')

if [[ "$bio" == "$newBio" ]]; then
  echo "TEST $test_count: ✅ PASS - Bio of admin user changed successfully!"
elif [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - Bio of admin user not changed"
fi

((++test_count))

# Negative: changing other users' bio as admin

res=$(curl -X PUT -sS \
    -H "$content_type_header" \
    -d '{"bio":"'"$attemptAdminBio"'"}' \
    --location "$base_url/users/$id/profile" \
    -H "$auth_header $adminToken")

err=$(echo "$res" | jq -r '.error')
bio=$(echo "$res" | jq -r '.bio')

if [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - (Admin) user can't change other's bio"
elif [[ "$bio" == "$attemptAdminBio" ]]; then
  echo "TEST $test_count: ❌ FAIL - (Admin) user can change other's bio"
fi

((++test_count))

# Negative: getting all users as admin without admin key

res=$(curl -X GET -sS \
    -H "$content_type_header" \
    --location "$base_url/admin/users" \
    -H "$auth_header $adminToken")

err=$(echo "$res" | jq -r '.error')
users=$(echo "$res" | jq -r '.users')

if [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - Admin cannot access all users without admin key"
elif [[ -n "$users" && "$users" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - Admin can access all users without admin key"
fi

((++test_count))

# Positive: getting all users as admin with admin key

res=$(curl -X GET -sS \
    -H "$content_type_header" \
    --location "$base_url/admin/users" \
    -H "$auth_header $adminToken" \
    -H "$x_admin_header $adminKey")

err=$(echo "$res" | jq -r '.error')
users=$(echo "$res" | jq -r '.users')

if [[ -n "$users" && "$users" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - Admin can only access all users with admin key"
elif  [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - Admin cannot access all users even with admin key"
fi

((++test_count))

# Positive: register a new account with valid email and strong password

valid_email_username=$(LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 10)

valid_password=$(
  {
    LC_ALL=C tr -dc 'a-z' </dev/urandom | head -c 1
    LC_ALL=C tr -dc 'A-Z' </dev/urandom | head -c 1
    LC_ALL=C tr -dc '!@#$%^&*()-_=+' </dev/urandom | head -c 1
    LC_ALL=C tr -dc 'a-zA-Z0-9!@#$%^&*()-_=+' </dev/urandom | head -c 14
  } | fold -w17 | head -n1
)


new_user_body='{"email":"'$valid_email_username'@example.com", "password":"'$valid_password'", "name":"just a name"}'

res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$new_user_body" \
    --location "$base_url/auth/register")

err=$(echo "$res" | jq -r '.error')
message=$(echo "$res" | jq -r '.message')

if [[ -n "$message" && "$message" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - New user can register with valid email and password"
elif  [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - New user cannot register with valid email and password"
fi

((++test_count))

# Negative: register a new account with valid email and weak password

new_user_body='{"email":"new-user2@example.com", "password":"weak", "name":"just a name"}'

res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$new_user_body" \
    --location "$base_url/auth/register")

err=$(echo "$res" | jq -r '.error')
message=$(echo "$res" | jq -r '.message')

if [[ -n "$message" && "$message" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - New user can register with valid email and weak password"
elif  [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - New user cannot register with valid email and weak password"
fi

((++test_count))

# Negative: register a new account with invalid email and strong password

new_user_body='{"email":"news", "password":"cM.gGA?)okM_]JWY@aq-", "name":"just a name"}'

res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$new_user_body" \
    --location "$base_url/auth/register")

err=$(echo "$res" | jq -r '.error')
message=$(echo "$res" | jq -r '.message')

if [[ -n "$message" && "$message" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - New user can register with invalid email and strong password"
elif  [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - New user cannot register with invalid email and strong password"
fi

((++test_count))

# Negative: register a new account with invalid email and weak password

new_user_body='{"email":"news", "password":"weak", "name":"just a name"}'

res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$new_user_body" \
    --location "$base_url/auth/register")

err=$(echo "$res" | jq -r '.error')
message=$(echo "$res" | jq -r '.message')

if [[ -n "$message" && "$message" != "null" ]]; then
  echo "TEST $test_count: ❌ FAIL - New user can register with invalid email and weak password"
elif  [[ -n "$err" && "$err" != "null" ]]; then
  echo "TEST $test_count: ✅ PASS - New user cannot register with invalid email and weak password"
fi

((++test_count))

# Negative: try to login with invalid credentials more than 5 times

non_existing_user_body='{"email": "test1234@example.com","password": "arbitrary password"}'

for i in {1..5}; do
  res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$non_existing_user_body" \
    -w '\n%{http_code}' \
    --location "$base_url/auth/login")
done

res=$(curl -X POST -sS \
    -H "$content_type_header" \
    -d "$non_existing_user_body" \
    -w '\n%{http_code}' \
    --location "$base_url/auth/login")

http_code=$(echo "$res" | tail -n1)

if [[ "$http_code" == "429" ]]; then
  echo "TEST $test_count: ✅ PASS - Rate limiting works"
else 
  echo "TEST $test_count: ❌ FAIL - Rate limiting does not work"
fi