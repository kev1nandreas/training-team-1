# Hardcoded Credentials Test Cases

## Test Case 1: Backend Source Code — JWT Secret
**Test ID**: SECRET-001  
**Priority**: High  
**Category**: Hardcoded Secrets  

**Objective**: Verify that JWT secret is not hardcoded in source code

**Preconditions**:
- Access to backend source code

**Test Steps**:
1. Search for hardcoded JWT secret in backend source
2. Verify environment variable is used instead

**Test Commands**:
```bash
# Search for hardcoded secret
grep -r "supersecret123" backend/
grep -r "jwtSecret.*=.*\"" backend/handlers/auth.go

# Check if os.Getenv is used
grep -r "os.Getenv.*JWT" backend/
```

**Expected Result (Vulnerable)**:
- `var jwtSecret = []byte("supersecret123")` found in auth.go
- Hardcoded string in source code

**Expected Result (Secure)**:
- No hardcoded secret found
- `os.Getenv("JWT_SECRET")` used to retrieve from environment

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: Backend Source Code — Database Credentials
**Test ID**: SECRET-002  
**Priority**: High  

**Objective**: Verify that database credentials are not hardcoded in source code

**Test Steps**:
1. Search for hardcoded DB password in backend source
2. Verify environment variables used in database connection

**Test Commands**:
```bash
# Search for hardcoded DB credentials
grep -r "taskpass" backend/
grep -r "taskuser" backend/
grep -r "host=localhost.*password" backend/

# Check if env vars are used
grep -r "os.Getenv.*DB" backend/database/
```

**Expected Result (Vulnerable)**:
- `DB_CONNECTION = "host=localhost user=taskuser password=taskpass123..."` in main.go
- Hardcoded DSN in database.go

**Expected Result (Secure)**:
- No hardcoded credentials
- `os.Getenv("DB_HOST")`, `os.Getenv("DB_PASSWORD")`, etc. used

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: Backend Source Code — Admin API Key
**Test ID**: SECRET-003  
**Priority**: High  

**Objective**: Verify that admin API key is not hardcoded in source code

**Test Steps**:
1. Search for hardcoded admin key in backend source
2. Verify environment variable is used

**Test Commands**:
```bash
# Search for hardcoded admin key
grep -r "admin-key-12345" backend/
grep -r "ADMIN_KEY.*=.*\"" backend/

# Check if env var is used
grep -r "os.Getenv.*ADMIN" backend/
```

**Expected Result (Vulnerable)**:
- `ADMIN_KEY = "admin-key-12345"` in main.go constants
- Hardcoded key value in source

**Expected Result (Secure)**:
- No hardcoded key
- `os.Getenv("ADMIN_API_KEY")` used

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: Frontend Source Code — config.js Secrets
**Test ID**: SECRET-004  
**Priority**: High  

**Objective**: Verify that sensitive config.js with hardcoded secrets has been removed

**Test Steps**:
1. Check if config.js exists in frontend source
2. Search for hardcoded API keys, credentials, AWS keys

**Test Commands**:
```bash
# Check if config.js still exists
ls frontend/src/config.js 2>/dev/null && echo "EXISTS — VULNERABLE" || echo "REMOVED — SECURE"

# Search for hardcoded secrets in frontend
grep -r "admin-key-12345" frontend/src/
grep -r "AKIAIOSFODNN" frontend/src/
grep -r "wJalrXUtnFEMI" frontend/src/
grep -r "DEFAULT_CREDENTIALS" frontend/src/
```

**Expected Result (Vulnerable)**:
- `config.js` exists with ADMIN_API_KEY, DEFAULT_CREDENTIALS, AWS keys
- Secrets visible in browser source

**Expected Result (Secure)**:
- `config.js` removed or contains no secrets
- Environment variables used via `import.meta.env`
- No AWS keys or default credentials in source

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: .env File in .gitignore
**Test ID**: SECRET-005  
**Priority**: High  

**Objective**: Verify that .env file is excluded from version control

**Test Steps**:
1. Check if .env is listed in .gitignore
2. Verify .env is not tracked by Git

**Test Commands**:
```bash
# Check .gitignore
grep "\.env" .gitignore

# Check if .env is tracked
git ls-files --error-unmatch .env 2>&1
# Should return error "not in index" if properly ignored
```

**Expected Result (Vulnerable)**:
- `.env` not in `.gitignore`
- `.env` tracked and committed to repository

**Expected Result (Secure)**:
- `.env` listed in `.gitignore`
- `.env` not tracked by Git
- `.env.example` provided with empty placeholder values

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 6: Hardcoded Constants in main.go
**Test ID**: SECRET-006  
**Priority**: High  

**Objective**: Verify that hardcoded credential constants have been removed from main.go

**Test Steps**:
1. Search for const block with credentials in main.go
2. Verify no sensitive constants exist

**Test Commands**:
```bash
# Search for the vulnerable constants
grep -A5 "^const" backend/main.go
grep "DB_CONNECTION\|JWT_SECRET\|ADMIN_KEY" backend/main.go
```

**Expected Result (Vulnerable)**:
- `const` block with DB_CONNECTION, JWT_SECRET, ADMIN_KEY

**Expected Result (Secure)**:
- No credential constants in source
- Environment loaded via `godotenv.Load()` or similar

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 7: .env.example Provided
**Test ID**: SECRET-007  
**Priority**: Medium  

**Objective**: Verify that a .env.example file exists with empty placeholder values

**Test Steps**:
1. Check for .env.example file existence
2. Verify it contains required variable names without actual values

**Test Commands**:
```bash
# Check if .env.example exists
cat .env.example

# Verify it has placeholder values, not real secrets
grep -v "^#" .env.example | grep "="
```

**Expected Result (Vulnerable)**:
- No .env.example file
- Or .env.example contains actual secret values

**Expected Result (Secure)**:
- `.env.example` exists with variable names
- Values are empty or clearly marked as placeholders
- Instructions provided for generating secrets

**Status**: [ ] Pass [ ] Fail [ ] Blocked

