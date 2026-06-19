# Brute Force / Rate Limiting Test Cases

## Test Case 1: Rate Limiting After 5 Failed Attempts
**Test ID**: RATE-001  
**Priority**: Medium  
**Category**: Brute Force / Authentication  

**Objective**: Verify that login endpoint blocks after 5 failed login attempts

**Preconditions**:
- Application is running
- A valid user account exists (e.g., user@example.com)
- No recent failed login attempts for the test email

**Test Steps**:
1. Send 5 failed login attempts with wrong password
2. Send 6th login attempt
3. Verify 6th attempt returns 429 Too Many Requests

**Test Commands**:
```bash
# Send 6 login attempts with wrong password
for i in {1..6}; do
  echo -n "Attempt $i: "
  curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"wrongpassword"}'
  echo
done

# Expected output:
# Attempt 1: 401
# Attempt 2: 401
# Attempt 3: 401
# Attempt 4: 401
# Attempt 5: 401
# Attempt 6: 429
```

**Expected Result (Vulnerable)**:
- All 6 (and more) attempts return 401
- No rate limiting or lockout
- Unlimited brute-force possible

**Expected Result (Secure)**:
- First 5 attempts return 401 (invalid credentials)
- 6th attempt returns 429 Too Many Requests
- Error message includes lockout duration

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: Lockout Duration (15 Minutes)
**Test ID**: RATE-002  
**Priority**: Medium  

**Objective**: Verify that lockout expires after 15 minutes

**Preconditions**:
- Application is running
- Account is locked out (5+ failed attempts made)

**Test Steps**:
1. Trigger lockout with 6 failed attempts
2. Confirm locked (429 response)
3. Wait 15 minutes
4. Try again — should be allowed

**Test Commands**:
```bash
# Trigger lockout
for i in {1..6}; do
  curl -s -o /dev/null -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"lockout-test@example.com","password":"wrong"}'
done

# Confirm lockout
echo -n "Immediately after lockout: "
curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"lockout-test@example.com","password":"wrong"}'
echo
# Expected: 429

# Wait 15 minutes...
sleep 900

# Try again after lockout period
echo -n "After 15 minutes: "
curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"lockout-test@example.com","password":"wrong"}'
echo
# Expected: 401 (not 429 — lockout expired, counter reset)
```

**Expected Result (Vulnerable)**:
- No lockout mechanism exists
- All attempts always processed

**Expected Result (Secure)**:
- Returns 429 during lockout period
- Returns 401 after lockout expires (15 minutes)
- Counter is reset after expiry

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: Successful Login Resets Counter
**Test ID**: RATE-003  
**Priority**: Medium  

**Objective**: Verify that a successful login resets the failed attempt counter

**Preconditions**:
- Application is running
- User account exists with known password

**Test Steps**:
1. Make 3 failed login attempts
2. Login successfully with correct password
3. Make 5 more failed attempts
4. Verify that the 6th (total) does NOT lock out (counter was reset)

**Test Commands**:
```bash
# 3 failed attempts
for i in {1..3}; do
  echo -n "Failed attempt $i: "
  curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"wrong"}'
  echo
done
# Expected: 401, 401, 401

# Successful login (resets counter)
echo -n "Successful login: "
curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
echo
# Expected: 200

# 5 more failed attempts + 1 to trigger lockout
for i in {1..6}; do
  echo -n "After reset, attempt $i: "
  curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"wrong"}'
  echo
done
# Expected: 401, 401, 401, 401, 401, 429
# Counter was reset by successful login — 5 new attempts allowed before block
```

**Expected Result (Vulnerable)**:
- No counter exists
- Unlimited attempts always possible

**Expected Result (Secure)**:
- Successful login resets the failed counter
- After reset, 5 new failed attempts allowed before lockout

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: Rate Limit Response Body
**Test ID**: RATE-004  
**Priority**: Low  

**Objective**: Verify that rate limit response provides helpful information

**Preconditions**:
- Account is locked out

**Test Steps**:
1. Trigger lockout (6 failed attempts)
2. Check response body of 429 response

**Test Commands**:
```bash
# Trigger lockout and check response
for i in {1..5}; do
  curl -s -o /dev/null -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"wrong"}'
done

# Check the lockout response body
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"wrong"}'

# Expected response: {"error":"too many login attempts. please try again after 15 minutes"}
```

**Expected Result (Vulnerable)**:
- No 429 response ever returned
- Or unhelpful error message

**Expected Result (Secure)**:
- Clear error message indicating lockout
- Mentions wait duration (15 minutes)
- Does NOT reveal whether the email exists

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Mass Brute Force Script
**Test ID**: RATE-005  
**Priority**: High  

**Objective**: Verify that rapid automated login attempts are blocked

**Preconditions**:
- Application is running

**Test Steps**:
1. Run automated script sending 100 login attempts in rapid succession
2. Verify that most are blocked after the 5th attempt

**Test Commands**:
```bash
#!/bin/bash
# Rapid-fire 20 requests
BLOCKED=0
ALLOWED=0

for i in {1..20}; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"bruteforce@example.com","password":"attempt'$i'"}')
  
  if [ "$STATUS" = "429" ]; then
    BLOCKED=$((BLOCKED + 1))
  else
    ALLOWED=$((ALLOWED + 1))
  fi
done

echo "Allowed: $ALLOWED"
echo "Blocked: $BLOCKED"
# Expected: Allowed ~5, Blocked ~15
```

**Expected Result (Vulnerable)**:
- All 20 requests return 401 (none blocked)
- Blocked count = 0

**Expected Result (Secure)**:
- First 5 return 401 (allowed to attempt)
- Remaining 15 return 429 (blocked by rate limiter)

**Status**: [ ] Pass [ ] Fail [ ] Blocked

