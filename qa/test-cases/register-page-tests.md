# Register Page Security Test Cases

## Test Case 1: Error Data Exposure on Registration Failure
**Test ID**: REG-001  
**Priority**: Medium  
**Category**: Information Disclosure  

**Objective**: Verify that detailed internal error data is NOT exposed to the user on registration failure

**Preconditions**:
- Application is running in browser
- Navigate to `/register`

**Test Steps**:
1. Try to register with an existing email (e.g., `user@example.com`)
2. Observe the error message displayed in the UI
3. Open DevTools → Console
4. Check if full error object is logged

**Expected Result (Vulnerable)**:
- Detailed error message like stack trace, internal field names, or SQL errors shown
- Console shows full `err` object with response details, headers, config
- Error reveals system internals (e.g., "duplicate key violation on email column")

**Expected Result (Secure)**:
- Generic, user-friendly error: "email already exists" or "Registration failed"
- Console only shows minimal log like `'Registration error:'` without sensitive data
- No internal server details exposed

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: No Email Validation
**Test ID**: REG-002  
**Priority**: Medium  

**Objective**: Verify that email field has proper format validation before submission

**Preconditions**:
- On register page

**Test Steps**:
1. Enter invalid email: `notanemail`
2. Enter invalid email: `user@`
3. Enter invalid email: `@domain.com`
4. Click Register

**Test Data**:
- `notanemail`
- `user@`
- `@domain.com`
- `user@domain` (missing TLD)
- `` (empty)

**Expected Result (Vulnerable)**:
- Form submits without client-side email validation
- Only server returns error (if it even validates)
- No feedback to user about format requirements

**Expected Result (Secure)**:
- Client-side validation: "Invalid email address"
- Form does NOT submit with invalid email
- Field turns red with error message below

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: No Password Complexity Validation
**Test ID**: REG-003  
**Priority**: High  

**Objective**: Verify that password field enforces complexity requirements on the client side

**Preconditions**:
- On register page

**Test Steps**:
1. Enter password: `abc` (too short, no uppercase, no symbol)
2. Enter password: `abcdefghijklmnop` (16 chars but no uppercase or symbol)
3. Enter password: `ABCDEFGHIJKLMNOP` (no lowercase or symbol)
4. Enter password: `Abcdefghijklmno!` (16 chars, has all requirements — valid)
5. Attempt to submit each

**Test Data & Expected Errors**:
| Password | Expected Error |
|----------|---------------|
| `abc` | "Password must be at least 16 characters" |
| `abcdefghijklmnop` | "Password must contain at least one uppercase letter" |
| `ABCDEFGHIJKLMNOP` | "Password must contain at least one lowercase letter" |
| `ABCDabcd12345678` | "Password must contain at least one symbol" |
| `Abcdefghijklmno!` | No error (valid) |

**Expected Result (Vulnerable)**:
- Any password accepted by the form
- No complexity feedback shown
- No minimum length enforced client-side

**Expected Result (Secure)**:
- Each rule validated separately with specific error message
- Password requirements listed below the field
- Form does not submit until all requirements met

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: No Name Validation
**Test ID**: REG-004  
**Priority**: Low  

**Objective**: Verify that name field has minimum length validation

**Preconditions**:
- On register page

**Test Steps**:
1. Leave name field empty
2. Enter single character: `a`
3. Try to submit

**Expected Result (Vulnerable)**:
- Empty or single-char name accepted
- Only `required` HTML attribute (easily bypassed)

**Expected Result (Secure)**:
- Validation error: "Name must be at least 2 characters"
- Maximum length enforced (100 chars)
- Field-level error displayed

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Password Requirements Visibility
**Test ID**: REG-005  
**Priority**: Low  

**Objective**: Verify that password requirements are clearly communicated to the user

**Preconditions**:
- On register page

**Test Steps**:
1. Navigate to `/register`
2. Look for password requirements text near the password field
3. Verify all rules are listed

**Expected Result (Vulnerable)**:
- No indication of password requirements
- User submits weak password, gets confusing server error
- Trial and error required to find correct format

**Expected Result (Secure)**:
- Text below password field: "Min 16 chars, one uppercase letter, one lowercase letter, one symbol"
- Requirements visible BEFORE submission attempt
- Clear guidance for user

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 6: Console Error Logging
**Test ID**: REG-006  
**Priority**: Medium  

**Objective**: Verify that registration errors are not logged with full details to console

**Preconditions**:
- On register page
- DevTools Console open

**Test Steps**:
1. Open DevTools → Console
2. Clear console
3. Attempt registration with duplicate email
4. Check console output

**Expected Result (Vulnerable)**:
- `console.error('Registration error:', err)` logs full Axios error object
- Response data, headers, request config all visible
- Stack trace exposed

**Expected Result (Secure)**:
- Only generic message logged: `'Registration error:'` with error object (acceptable for dev)
- No sensitive data (passwords, tokens) in console
- In production build: no console output at all

**Status**: [ ] Pass [ ] Fail [ ] Blocked

