# Authentication & Authorization Test Cases

## Test Case 1: Access Protected Endpoint Without Authentication
**Test ID**: AUTH-001  
**Priority**: Critical  
**Category**: Authentication  

**Objective**: Verify that protected endpoints reject unauthenticated requests

**Preconditions**:
- Application is running
- User is NOT logged in (no token)

**Test Steps**:
1. Send GET request to `/api/tasks` without Authorization header
2. Send DELETE request to `/api/tasks/1` without Authorization header
3. Send GET request to `/api/tasks/search?q=test` without Authorization header
4. Send PUT request to `/api/users/1/profile` without Authorization header
5. Send GET request to `/api/admin/users` without Authorization header

**Test Commands**:
```bash
# Test 1: Access tasks without auth
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/tasks

# Test 2: Delete task without auth
curl -s -o /dev/null -w "%{http_code}" -X DELETE http://localhost:8080/api/tasks/1

# Test 3: Search tasks without auth
curl -s -o /dev/null -w "%{http_code}" "http://localhost:8080/api/tasks/search?q=test"

# Test 4: Update profile without auth
curl -s -o /dev/null -w "%{http_code}" -X PUT http://localhost:8080/api/users/1/profile \
  -H "Content-Type: application/json" \
  -d '{"name":"Hacked"}'

# Test 5: Access admin endpoint without auth
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/admin/users
```

**Expected Result (Vulnerable)**:
- Returns 200 with data for all endpoints
- Actions performed without authentication

**Expected Result (Secure)**:
- All return 401 Unauthorized
- No data exposed

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: Delete Task Without Authentication
**Test ID**: AUTH-002  
**Priority**: Critical  

**Objective**: Verify that task deletion requires authentication

**Preconditions**:
- Application is running
- At least one task exists in the database

**Test Steps**:
1. Without logging in, send DELETE request to a task endpoint
2. Observe response

**Test Commands**:
```bash
# Try to delete task ID 1 without auth token
curl -s -w "\nHTTP Status: %{http_code}" -X DELETE http://localhost:8080/api/tasks/1

# Vulnerable: returns 200 and deletes the task
# Secure: returns 401 Unauthorized
```

**Expected Result (Vulnerable)**:
- Task is deleted
- Returns 200 OK

**Expected Result (Secure)**:
- Returns 401 Unauthorized
- Task remains intact

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: Authorization Bypass — Update Another User's Profile
**Test ID**: AUTH-003  
**Priority**: Critical  

**Objective**: Verify that users cannot update other users' profiles (horizontal privilege escalation)

**Preconditions**:
- Application is running
- Two users exist (user ID 1 and user ID 2)
- Logged in as user ID 2

**Test Steps**:
1. Login as regular user (user@example.com)
2. Get token
3. Try to update user ID 1's profile using user ID 2's token

**Test Commands**:
```bash
# Login as regular user (ID 2)
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq -r '.token')

# Try to update admin's profile (ID 1) using regular user's token
curl -s -w "\nHTTP Status: %{http_code}" -X PUT http://localhost:8080/api/users/1/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Hacked Admin","bio":"Pwned"}'

# Vulnerable: returns 200, admin profile updated
# Secure: returns 403 Forbidden
```

**Expected Result (Vulnerable)**:
- Admin profile is modified by regular user
- Returns 200 OK

**Expected Result (Secure)**:
- Returns 403 Forbidden "Cannot update other users' profiles"
- Admin profile remains unchanged

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: Admin Endpoint Without Admin Role
**Test ID**: AUTH-004  
**Priority**: Critical  

**Objective**: Verify that admin endpoints check user role server-side

**Preconditions**:
- Application is running
- Logged in as regular user (role: "user")

**Test Steps**:
1. Login as regular user
2. Access admin endpoint with regular user's token

**Test Commands**:
```bash
# Login as regular user
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq -r '.token')

# Try to access admin users list
curl -s -w "\nHTTP Status: %{http_code}" http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Admin-Key: some-random-key"

# Vulnerable: returns 200 with all users
# Secure: returns 403 Forbidden
```

**Expected Result (Vulnerable)**:
- Returns list of all users including passwords
- No role check performed

**Expected Result (Secure)**:
- Returns 403 "Admin access required"
- No user data exposed

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Admin Endpoint Without Valid API Key
**Test ID**: AUTH-005  
**Priority**: High  

**Objective**: Verify that admin endpoints require valid X-Admin-Key header

**Preconditions**:
- Application is running
- Logged in as admin user

**Test Steps**:
1. Login as admin user
2. Access admin endpoint without X-Admin-Key header
3. Access admin endpoint with invalid X-Admin-Key

**Test Commands**:
```bash
# Login as admin
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' | jq -r '.token')

# Try without X-Admin-Key
curl -s -w "\nHTTP Status: %{http_code}" http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $TOKEN"

# Try with wrong X-Admin-Key
curl -s -w "\nHTTP Status: %{http_code}" http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Admin-Key: wrong-key"

# Vulnerable: returns 200 without key check
# Secure: returns 401 "Invalid admin API key"
```

**Expected Result (Vulnerable)**:
- Returns all users without requiring valid key
- No key verification

**Expected Result (Secure)**:
- Returns 401 Unauthorized without valid key
- Only grants access with correct key + admin role

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 6: Client-Side Authorization Bypass
**Test ID**: AUTH-006  
**Priority**: High  

**Objective**: Verify that modifying client-side role does not grant server-side access

**Preconditions**:
- Application is running in browser
- Logged in as regular user

**Test Steps**:
1. Login as regular user (user@example.com)
2. Open DevTools → Application → Cookies/Storage
3. If user data is in state, note that role is "user"
4. Navigate to `/admin` in the browser
5. Observe if admin panel loads and if API returns data

**Expected Result (Vulnerable)**:
- Admin panel renders
- API returns all users because no server-side check

**Expected Result (Secure)**:
- Navigation redirected to dashboard (frontend guard)
- Even if bypassed, API returns 403 (server-side guard)

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 7: Delete Task Owned by Another User
**Test ID**: AUTH-007  
**Priority**: High  

**Objective**: Verify that users can only delete their own tasks

**Preconditions**:
- Application is running
- User A and User B both have tasks
- Logged in as User B

**Test Steps**:
1. Login as user@example.com
2. Try to delete a task that belongs to admin@example.com

**Test Commands**:
```bash
# Login as regular user
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq -r '.token')

# Create a task as admin first to know the ID
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' | jq -r '.token')

TASK_ID=$(curl -s -X POST http://localhost:8080/api/tasks \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Admin Task","description":"owned by admin","priority":"high"}' | jq -r '.id')

# Try to delete admin's task as regular user
curl -s -w "\nHTTP Status: %{http_code}" -X DELETE "http://localhost:8080/api/tasks/$TASK_ID" \
  -H "Authorization: Bearer $TOKEN"

# Vulnerable: task deleted, returns 200
# Secure: returns 404 "Task not found" (scoped to user_id)
```

**Expected Result (Vulnerable)**:
- Task is deleted regardless of ownership
- Returns 200 OK

**Expected Result (Secure)**:
- Returns 404 Not Found (query scoped to user's tasks)
- Admin's task remains intact

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 8: Search Tasks Scoped to Authenticated User
**Test ID**: AUTH-008  
**Priority**: High  

**Objective**: Verify search results are scoped to the authenticated user's tasks only

**Preconditions**:
- Both admin and regular user have tasks
- Logged in as regular user

**Test Steps**:
1. Login as regular user
2. Search for a term that matches admin's tasks
3. Verify only user's own tasks are returned

**Test Commands**:
```bash
# Login as regular user
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq -r '.token')

# Search — should only return tasks owned by this user
curl -s "http://localhost:8080/api/tasks/search?q=task" \
  -H "Authorization: Bearer $TOKEN" | jq '.[] | .user_id'

# Vulnerable: returns tasks from all users
# Secure: returns only tasks where user_id matches authenticated user
```

**Expected Result (Vulnerable)**:
- Returns tasks from all users
- Leaks other users' data

**Expected Result (Secure)**:
- Only returns tasks owned by the authenticated user
- No data leakage from other accounts

**Status**: [ ] Pass [ ] Fail [ ] Blocked

