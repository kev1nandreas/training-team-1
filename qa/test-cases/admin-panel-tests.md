# Admin Panel Security Test Cases

## Test Case 1: Password Column Visibility
**Test ID**: ADMIN-001  
**Priority**: High  
**Category**: Data Exposure  

**Objective**: Verify that users' passwords are NOT shown in the admin panel data table

**Preconditions**:
- Application is running
- Logged in as admin (admin@example.com)
- Navigate to `/admin`

**Test Steps**:
1. Login as admin
2. Navigate to `/admin`
3. Observe the users table
4. Check table headers for "Password" column
5. Check table rows for any password data

**Expected Result (Vulnerable)**:
- "Password" column header visible in table
- Plain text passwords displayed for each user (e.g., "password123", "admin123")
- Password data rendered in red monospace font

**Expected Result (Secure)**:
- No "Password" column in table headers
- Table only shows: ID, Name, Email, Role
- No password data anywhere on the page

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: API Response Excludes Password
**Test ID**: ADMIN-002  
**Priority**: High  

**Objective**: Verify that the admin users API endpoint does not return password fields

**Preconditions**:
- Logged in as admin with valid token and API key

**Test Steps**:
1. Login as admin
2. Call GET /api/admin/users with proper authorization
3. Inspect response JSON for password field

**Test Commands**:
```bash
# Login as admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' | jq -r '.token')

# Get all users and check for password field
curl -s http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "X-Admin-Key: <admin-api-key>" | jq '.users[0] | keys'

# Check specifically for password
curl -s http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "X-Admin-Key: <admin-api-key>" | jq '.users[0].password'

# Vulnerable: returns actual password string
# Secure: returns null (field not in response)
```

**Expected Result (Vulnerable)**:
- Response JSON includes `"password": "password123"` or hash
- `jq '.users[0].password'` returns a value

**Expected Result (Secure)**:
- No `password` key in response
- `jq '.users[0].password'` returns `null`
- User model uses `json:"-"` tag

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: Admin Panel Access Control (Frontend)
**Test ID**: ADMIN-003  
**Priority**: High  

**Objective**: Verify that non-admin users cannot access the admin panel via frontend routing

**Preconditions**:
- Logged in as regular user (role: "user")

**Test Steps**:
1. Login as regular user (user@example.com)
2. Manually navigate to `/admin` in browser URL bar
3. Observe behavior

**Expected Result (Vulnerable)**:
- Admin panel renders
- Only a console warning about non-admin access
- No redirect

**Expected Result (Secure)**:
- User redirected to `/dashboard`
- Admin panel does not render
- withAuth HOC blocks access based on role

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: Admin Panel Access Control (Backend)
**Test ID**: ADMIN-004  
**Priority**: Critical  

**Objective**: Verify that the admin API rejects non-admin users even if frontend is bypassed

**Preconditions**:
- Logged in as regular user

**Test Steps**:
1. Login as regular user
2. Directly call admin API endpoint with regular user's token

**Test Commands**:
```bash
# Login as regular user
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq -r '.token')

# Try to access admin endpoint
curl -s -w "\nHTTP Status: %{http_code}" http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Admin-Key: any-key-here"

# Vulnerable: 200 with all user data
# Secure: 403 "Admin access required"
```

**Expected Result (Vulnerable)**:
- Returns 200 with all users data
- No server-side role verification

**Expected Result (Secure)**:
- Returns 403 Forbidden
- Error: "Admin access required"
- Server checks JWT role claim

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Client-Side Role Manipulation
**Test ID**: ADMIN-005  
**Priority**: High  

**Objective**: Verify that modifying role in client state/storage does not bypass server-side checks

**Preconditions**:
- Logged in as regular user
- Browser DevTools available

**Test Steps**:
1. Login as regular user
2. Open DevTools → Application
3. If possible, modify user state to set role to "admin"
4. Navigate to `/admin`
5. Observe if data loads (even if page renders, API should reject)

**Expected Result (Vulnerable)**:
- Changing localStorage role to "admin" grants full access
- Admin panel loads with user data
- Server does not validate role from JWT

**Expected Result (Secure)**:
- Even if frontend renders admin page, API call fails with 403
- Server role check based on JWT token claims (not modifiable client-side)
- No user data loaded

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 6: Admin Panel Table Fields
**Test ID**: ADMIN-006  
**Priority**: Medium  

**Objective**: Verify that only appropriate, non-sensitive fields are displayed in the admin table

**Preconditions**:
- Logged in as admin
- Admin panel loaded with users

**Test Steps**:
1. Navigate to `/admin`
2. Inspect table columns
3. Verify only safe fields are shown

**Expected Fields (Secure)**:
- ID
- Name
- Email
- Role

**Fields that should NOT be present**:
- Password
- Hashed password
- API tokens
- Internal timestamps (unless needed)
- Bio (unless intentional)

**Expected Result (Vulnerable)**:
- Sensitive columns visible (password, tokens)
- Full user object rendered without field selection

**Expected Result (Secure)**:
- Only ID, Name, Email, Role columns visible
- No sensitive data in table

**Status**: [ ] Pass [ ] Fail [ ] Blocked

