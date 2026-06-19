# SQL Injection Test Cases

## Test Case 1: Basic SQL Injection
**Test ID**: SQL-001  
**Priority**: Critical  
**Category**: SQL Injection  

**Objective**: Verify that search functionality is vulnerable to basic SQL injection

**Preconditions**:
- Application is running
- User is logged in
- At least one task exists

**Test Steps**:
1. Navigate to Dashboard
2. Enter in search box: `' OR '1'='1`
3. Click Search button

**Expected Result (Vulnerable)**:
- All tasks are returned
- SQL query is manipulated

**Expected Result (Secure)**:
- No results or only legitimate matches
- Special characters handled safely

**Test Data**: `' OR '1'='1`

**Status**: [x] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: Union-Based SQL Injection
**Test ID**: SQL-002  
**Priority**: Critical  

**Objective**: Test for union-based SQL injection to extract data from other tables

**Test Steps**:
1. Navigate to Dashboard
2. Enter in search: `' UNION SELECT id, email, password, 'x', 'x' FROM users--`
3. Submit search

**Expected Result (Vulnerable)**:
- Returns user credentials
- Exposes database structure

**Expected Result (Secure)**:
- Query is safely handled
- No data leakage

**Test Data**: `' UNION SELECT * FROM users--`

**Status**: [x] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: Comment-Based SQL Injection
**Test ID**: SQL-003  
**Priority**: High  

**Objective**: Test SQL injection using comment syntax

**Test Steps**:
1. Use search with payload: `admin'--`
2. Observe results

**Expected Result (Vulnerable)**:
- Comment bypasses query logic
- Returns unintended results

**Expected Result (Secure)**:
- Comments are escaped
- Query executes safely

**Test Data**: 
- `admin'--`
- `' OR 1=1--`
- `'; --`

**Status**: [x] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: Time-Based Blind SQL Injection
**Test ID**: SQL-004  
**Priority**: High  

**Objective**: Test for time-based blind SQL injection

**Test Steps**:
1. Enter payload: `' OR SLEEP(5)--`
2. Measure response time

**Expected Result (Vulnerable)**:
- Response delayed by 5 seconds
- Database command executed

**Expected Result (Secure)**:
- Normal response time
- Payload not executed

**Test Data**: `' OR SLEEP(5)--`

**Status**: [x] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Error-Based SQL Injection
**Test ID**: SQL-005  
**Priority**: Medium  

**Objective**: Test if SQL errors expose database information

**Test Steps**:
1. Enter invalid SQL: `'`
2. Check error message

**Expected Result (Vulnerable)**:
- SQL error message displayed
- Database info leaked

**Expected Result (Secure)**:
- Generic error message
- No technical details exposed

**Test Data**: `'`, `''`, `')'`

**Status**: [x] Pass [ ] Fail [ ] Blocked


---

## Test Case 6: Boolean Injection with Comment
**Test ID**: SQL-006  
**Priority**: Critical  
**Category**: SQL Injection  

**Objective**: Verify that boolean injection with SQL comment is handled safely

**Preconditions**:
- Application is running
- User is logged in

**Test Steps**:
1. Navigate to Dashboard
2. Enter in search box: `' OR 1=1--`
3. Click Search button

**Expected Result (Vulnerable)**:
- All tasks are returned
- Comment truncates rest of query

**Expected Result (Secure)**:
- No results or only legitimate matches
- Payload treated as literal search term

**Test Data**: `' OR 1=1--`

**Status**: [x] Pass [ ] Fail [ ] Blocked

---

## Test Case 7: Destructive SQL Injection (DROP TABLE)
**Test ID**: SQL-007  
**Priority**: Critical  
**Category**: SQL Injection  

**Objective**: Verify that destructive SQL commands cannot be injected

**Preconditions**:
- Application is running
- User is logged in
- Tasks table exists with data

**Test Steps**:
1. Navigate to Dashboard
2. Enter in search box: `'; DROP TABLE tasks; --`
3. Click Search button
4. Verify tasks table still exists (try loading tasks)

**Expected Result (Vulnerable)**:
- Tasks table is deleted
- Subsequent requests fail

**Expected Result (Secure)**:
- No destructive action occurs
- Tasks table intact
- Normal response (empty results or error)

**Test Data**: `'; DROP TABLE tasks; --`

**Status**: [x] Pass [ ] Fail [ ] Blocked

---

## Test Case 8: Boolean Injection Variant
**Test ID**: SQL-008  
**Priority**: High  
**Category**: SQL Injection  

**Objective**: Verify that string-based boolean injection is handled safely

**Preconditions**:
- Application is running
- User is logged in

**Test Steps**:
1. Navigate to Dashboard
2. Enter in search box: `' OR 'a'='a`
3. Click Search button

**Expected Result (Vulnerable)**:
- All tasks are returned
- Boolean condition always evaluates to true

**Expected Result (Secure)**:
- No results or only legitimate matches
- Special characters handled safely

**Test Data**: `' OR 'a'='a`

**Status**: [x] Pass [ ] Fail [ ] Blocked
