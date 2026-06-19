# Profile Page Security Test Cases

## Test Case 1: Debug Info Exposure
**Test ID**: PROFILE-001  
**Priority**: High  
**Category**: Data Exposure  

**Objective**: Verify that sensitive user data (debug info) is NOT displayed on the profile page

**Preconditions**:
- Application is running in browser
- User is logged in

**Test Steps**:
1. Login to application
2. Navigate to `/profile`
3. Scroll to bottom of page
4. Look for any "Debug Info" section or raw JSON output

**Expected Result (Vulnerable)**:
- "Debug Info" section visible
- Full JSON user object displayed including password, internal fields
- `JSON.stringify(user, null, 2)` rendered in `<pre>` tag

**Expected Result (Secure)**:
- No debug section on the page
- No raw JSON displayed
- Only form fields (name, bio) visible with proper labels

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: XSS via Bio Field (dangerouslySetInnerHTML)
**Test ID**: PROFILE-002  
**Priority**: High  

**Objective**: Verify that bio content is NOT rendered as raw HTML (XSS prevention)

**Preconditions**:
- User is logged in
- User can update profile

**Test Steps**:
1. Navigate to `/profile`
2. Update bio with: `<script>alert('XSS')</script>`
3. Save profile
4. Observe the "Current Bio" section

**Test Payloads**:
- `<script>alert('XSS')</script>`
- `<img src=x onerror=alert(document.cookie)>`
- `<svg onload=alert('XSS')>`
- `<iframe src="javascript:alert('XSS')"></iframe>`

**Expected Result (Vulnerable)**:
- Alert box pops up
- JavaScript executes when bio is rendered
- Using `dangerouslySetInnerHTML={{ __html: user.bio }}`

**Expected Result (Secure)**:
- Bio displayed as plain text
- HTML tags shown as text or stripped
- Using React's default text rendering: `{user.bio}`

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: No Input Validation on Name Field
**Test ID**: PROFILE-003  
**Priority**: Medium  

**Objective**: Verify that name field has proper validation

**Preconditions**:
- User is logged in
- On profile page

**Test Steps**:
1. Clear name field entirely
2. Click "Update Profile"
3. Observe response

**Test Data**:
- Empty string: ``
- Single character: `a`
- Very long name: `A` repeated 200 times

**Expected Result (Vulnerable)**:
- Empty name accepted and saved
- No character limit enforced
- No minimum length check

**Expected Result (Secure)**:
- Validation error: "Name must be at least 2 characters"
- Maximum length enforced (100 chars)
- Field-level error message displayed

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: No Input Validation on Bio Field
**Test ID**: PROFILE-004  
**Priority**: Medium  

**Objective**: Verify that bio field has proper length validation

**Preconditions**:
- User is logged in
- On profile page

**Test Steps**:
1. Enter bio text exceeding 500 characters
2. Click "Update Profile"
3. Observe response

**Test Data**:
- 501 characters of text
- 10000 characters of text

**Expected Result (Vulnerable)**:
- Unlimited length accepted
- No character counter shown
- Possible database/display issues with very long text

**Expected Result (Secure)**:
- Validation error: "Bio must be 500 characters or fewer"
- Character counter visible (e.g., "450/500")
- Form does not submit if limit exceeded

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Password Field Not Exposed in Profile Data
**Test ID**: PROFILE-005  
**Priority**: High  

**Objective**: Verify that the user object rendered on profile page does not contain a password field

**Preconditions**:
- User is logged in

**Test Steps**:
1. Navigate to `/profile`
2. Open DevTools → Elements tab
3. Search for "password" in the DOM
4. Check if any password value is visible in the page source

**Verification (Browser Console)**:
```javascript
// Check if password is in any visible text on the page
document.body.innerText.includes('password123');
// Should be false
```

**Expected Result (Vulnerable)**:
- Password visible in debug section or DOM
- API response includes password field stored in state

**Expected Result (Secure)**:
- No password anywhere in rendered DOM
- API `/users/me` response excludes password (json:"-")

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 6: Profile Form Field-Level Error Display
**Test ID**: PROFILE-006  
**Priority**: Low  

**Objective**: Verify that validation errors are shown per-field with clear messages

**Preconditions**:
- User is logged in
- On profile page

**Test Steps**:
1. Clear name field (leave empty)
2. Enter bio with 501+ characters
3. Submit form
4. Observe error messages

**Expected Result (Vulnerable)**:
- No error messages shown
- Form submits silently
- Errors only appear as server-side response

**Expected Result (Secure)**:
- Field-level error under name: "Name must be at least 2 characters"
- Field-level error under bio: "Bio must be 500 characters or fewer"
- Red border on invalid fields
- Form does NOT submit until valid

**Status**: [ ] Pass [ ] Fail [ ] Blocked

