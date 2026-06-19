# Insecure Data Storage Test Cases

## Test Case 1: JWT Token Storage Method
**Test ID**: STORAGE-001  
**Priority**: High  
**Category**: Insecure Data Storage  

**Objective**: Verify that JWT tokens are NOT stored in localStorage (vulnerable to XSS theft)

**Preconditions**:
- Application is running in browser
- User is NOT logged in initially

**Test Steps**:
1. Open browser DevTools → Application → Local Storage
2. Login to the application
3. Check Local Storage for any token-related keys
4. Check Session Storage for any token-related keys
5. Check Cookies for token storage

**Verification Commands (Browser Console)**:
```javascript
// After login, check localStorage
console.log('localStorage token:', localStorage.getItem('token'));
console.log('localStorage keys:', Object.keys(localStorage));

// Check sessionStorage
console.log('sessionStorage keys:', Object.keys(sessionStorage));

// Check cookies
console.log('cookies:', document.cookie);
```

**Expected Result (Vulnerable)**:
- `localStorage.getItem('token')` returns the JWT string
- Token is accessible via JavaScript (XSS can steal it)

**Expected Result (Secure)**:
- `localStorage.getItem('token')` returns null
- Token stored in secure cookie (httpOnly preferred) or encrypted cookie
- Token NOT accessible via `document.cookie` if httpOnly

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: User Data in localStorage
**Test ID**: STORAGE-002  
**Priority**: High  

**Objective**: Verify that sensitive user data (especially passwords) is NOT stored in localStorage

**Preconditions**:
- Logged in to the application

**Test Steps**:
1. Login to application
2. Open DevTools → Application → Local Storage
3. Check for 'user' key
4. Inspect stored data for sensitive fields

**Verification Commands (Browser Console)**:
```javascript
// Check for user data in localStorage
const userData = localStorage.getItem('user');
console.log('User data in localStorage:', userData);

// If it exists, check for password field
if (userData) {
  const parsed = JSON.parse(userData);
  console.log('Has password field:', 'password' in parsed);
  console.log('Stored fields:', Object.keys(parsed));
}
```

**Expected Result (Vulnerable)**:
- Full user object stored in localStorage including password
- `JSON.parse(localStorage.getItem('user')).password` reveals the password

**Expected Result (Secure)**:
- No user object in localStorage
- User state managed in memory (e.g., Zustand/Redux store)
- No password field exposed anywhere in frontend storage

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: Debug Information in localStorage
**Test ID**: STORAGE-003  
**Priority**: Medium  

**Objective**: Verify that no debug information is stored in localStorage

**Preconditions**:
- Logged in to the application

**Test Steps**:
1. Login to application
2. Perform various actions (create tasks, search, navigate)
3. Check localStorage for debug-related keys

**Verification Commands (Browser Console)**:
```javascript
// Check for debug info
console.log('debugInfo:', localStorage.getItem('debugInfo'));
console.log('appSettings:', localStorage.getItem('appSettings'));

// List all localStorage keys
for (let i = 0; i < localStorage.length; i++) {
  const key = localStorage.key(i);
  console.log(`Key: ${key}, Value: ${localStorage.getItem(key)}`);
}
```

**Expected Result (Vulnerable)**:
- `debugInfo` key exists with user agent, timestamps, actions
- `appSettings` key exists with unnecessary configuration data
- Multiple sensitive keys present

**Expected Result (Secure)**:
- No `debugInfo` key in localStorage
- No `appSettings` key in localStorage
- Minimal or no data in localStorage

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: Token Encryption/Protection
**Test ID**: STORAGE-004  
**Priority**: High  

**Objective**: Verify that if token must be stored client-side, it is encrypted or protected

**Preconditions**:
- Application is running
- User logs in

**Test Steps**:
1. Login to the application
2. Check all storage mechanisms (localStorage, sessionStorage, cookies)
3. If token is found, verify it is NOT stored as plain JWT
4. Verify encryption is used

**Verification Commands (Browser Console)**:
```javascript
// Check cookies for auth token
const cookies = document.cookie.split(';').map(c => c.trim());
console.log('All cookies:', cookies);

// If cookie exists, check if it looks like raw JWT (three base64 parts separated by dots)
cookies.forEach(cookie => {
  const [name, value] = cookie.split('=');
  const jwtPattern = /^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$/;
  if (jwtPattern.test(value)) {
    console.log(`WARNING: Cookie '${name}' contains unencrypted JWT!`);
  } else {
    console.log(`Cookie '${name}' appears encrypted/protected`);
  }
});
```

**Expected Result (Vulnerable)**:
- Raw JWT visible in localStorage: `eyJhbGci...`
- Easily decoded with jwt.io
- XSS can steal and reuse directly

**Expected Result (Secure)**:
- Token encrypted before storage (e.g., AES-GCM)
- Stored in secure cookie with proper flags
- Not directly usable if intercepted without decryption key

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Sensitive Data in API Responses
**Test ID**: STORAGE-005  
**Priority**: High  

**Objective**: Verify that API responses do not include password fields

**Preconditions**:
- Logged in as any user

**Test Steps**:
1. Login and capture the response
2. Call GET /api/users/me
3. If admin, call GET /api/admin/users
4. Check all responses for password fields

**Test Commands**:
```bash
# Login and check response for password
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq '.user | has("password")'

# Check /users/me response
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq -r '.token')

curl -s http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN" | jq 'has("password")'

# Check admin users list (login as admin)
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' | jq -r '.token')

curl -s http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "X-Admin-Key: $(grep ADMIN_API_KEY .env | cut -d= -f2)" | jq '.users[0] | has("password")'

# Vulnerable: returns true (password field present)
# Secure: returns false (password excluded via json:"-")
```

**Expected Result (Vulnerable)**:
- Password field present in login response user object
- Password field present in /users/me response
- Passwords visible in admin user list

**Expected Result (Secure)**:
- No `password` field in any API response (json:"-" tag)
- User object only contains id, email, name, role, bio, timestamps

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 6: Console Logging of Sensitive Data
**Test ID**: STORAGE-006  
**Priority**: Medium  

**Objective**: Verify that API requests/responses are NOT logged to browser console

**Preconditions**:
- Application is running in browser

**Test Steps**:
1. Open browser DevTools → Console
2. Clear console
3. Login to application
4. Create a task
5. Search for tasks
6. Navigate between pages
7. Check console for sensitive data output

**Expected Result (Vulnerable)**:
- Console shows "API Request: POST /auth/login {email, password}"
- Console shows "API Response: {token: ..., user: {...}}"
- All API traffic logged with sensitive data

**Expected Result (Secure)**:
- No API request/response logging in console
- Only error messages without sensitive details
- Production builds have no console.log statements

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 7: Session Cleanup on Logout
**Test ID**: STORAGE-007  
**Priority**: Medium  

**Objective**: Verify that all sensitive data is cleared on logout

**Preconditions**:
- Logged in to the application

**Test Steps**:
1. Login to application
2. Note all storage keys before logout
3. Click Logout
4. Check all storage mechanisms are cleared

**Verification Commands (Browser Console)**:
```javascript
// Before logout - note what's stored
console.log('Before logout:');
console.log('localStorage keys:', Object.keys(localStorage));
console.log('sessionStorage keys:', Object.keys(sessionStorage));
console.log('cookies:', document.cookie);

// After logout - verify cleanup
// (run after clicking Logout)
console.log('After logout:');
console.log('localStorage keys:', Object.keys(localStorage));
console.log('sessionStorage keys:', Object.keys(sessionStorage));
console.log('cookies:', document.cookie);
```

**Expected Result (Vulnerable)**:
- Token remains in localStorage after logout
- User data remains in storage
- Session not properly invalidated

**Expected Result (Secure)**:
- All auth-related cookies removed
- No user data in localStorage/sessionStorage
- Memory state cleared (Zustand store reset)

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 8: Password Visibility in Admin Panel
**Test ID**: STORAGE-008  
**Priority**: High  

**Objective**: Verify that user passwords are NOT displayed in the admin panel table

**Preconditions**:
- Logged in as admin user
- Navigate to Admin Panel

**Test Steps**:
1. Login as admin@example.com
2. Navigate to /admin
3. Observe the users table
4. Check if password column exists

**Expected Result (Vulnerable)**:
- "Password" column visible in the table
- Plain text passwords displayed for all users

**Expected Result (Secure)**:
- No password column in the admin table
- API response does not include password field
- Only ID, Name, Email, Role columns shown

**Status**: [ ] Pass [ ] Fail [ ] Blocked

