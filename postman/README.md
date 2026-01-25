# Postman Testing Guide - 5 Core Implementation Tasks

## 📦 Setup Postman

### 1. Import Collection

1. Buka Postman
2. Klik **Import** (tombol di kiri atas)
3. Pilih file `postman/Testing-5-Tasks.postman_collection.json`
4. Klik **Import**

### 2. Import Environment

1. Klik ⚙️ **Settings** (kanan atas)
2. Pilih tab **Environments**
3. Klik **Import**
4. Pilih file `postman/Development.postman_environment.json`
5. Pilih environment **Development** di dropdown (kanan atas)

### 3. Jalankan Server

```bash
# Pastikan Redis berjalan
docker-compose -f docker-compose.redis.yml up -d

# Jalankan aplikasi
go run cmd/app/main.go
```

Server akan berjalan di `http://localhost:8080`

---

## 🧪 Testing Workflow

### Task 1: JWT Multi-Role Authentication

**Folder**: Task 1 - JWT Multi-Role Auth

#### 1.1 Register User

```
POST /api/auth/register
```

- Membuat user baru untuk testing
- User otomatis mendapat role "user"
- **Expected**: Status 201, user created

#### 1.2 Login and Get JWT Token

```
POST /api/auth/login
```

- Login dengan credentials
- Mendapatkan JWT token
- Token disimpan otomatis ke environment variable `{{auth_token}}`
- **Expected**: Status 200, token dengan multi-role

**Verifikasi**:

- ✅ Token contains UserID
- ✅ Token contains roles array
- ✅ Token is valid JWT format

#### 1.3 Verify Token Contains Multiple Roles

```
GET /api/auth/profile
Authorization: Bearer {{auth_token}}
```

- Verifikasi JWT token berisi role array
- **Expected**: Status 200, roles array exists

---

### Task 2: RBAC Token Payload

**Folder**: Task 2 - RBAC Token Payload

#### 2.1 Verify Token Payload Structure

```
GET /api/auth/profile
Authorization: Bearer {{auth_token}}
```

**Tests Performed**:

- ✅ Token contains `user_id`
- ✅ Token contains `username`
- ✅ Token contains `email`
- ✅ Token contains `roles` (RBAC)

**Expected Response**:

```json
{
  "data": {
    "id": "uuid-here",
    "username": "testuser",
    "email": "test@example.com",
    "roles": ["user"]
  }
}
```

#### 2.2 Check Token Type (User vs CMS)

```
GET /api/auth/profile
```

- Verifikasi middleware dapat membedakan user vs CMS token
- Token type ada di JWT payload
- **Expected**: Token type identified correctly

---

### Task 3: Middleware Token Identification

**Folder**: Task 3 - Middleware Token Identification

#### 3.1 Valid Bearer Token ✅

```
GET /api/auth/profile
Authorization: Bearer {{auth_token}}
```

- **Expected**: Status 200, request accepted

#### 3.2 Missing Authorization Header ❌

```
GET /api/auth/profile
(no Authorization header)
```

- **Expected**: Status 401, request rejected

#### 3.3 Invalid Bearer Format ❌

```
GET /api/auth/profile
Authorization: InvalidFormat token123
```

- **Expected**: Status 401, invalid format rejected

#### 3.4 Expired Token ❌

```
GET /api/auth/profile
Authorization: Bearer <expired-token>
```

- **Expected**: Status 401, expired token rejected

**Middleware Checks**:

- ✅ Extracts token from "Bearer {token}" format
- ✅ Validates token signature
- ✅ Checks token expiration
- ✅ Rejects invalid formats
- ✅ Distinguishes user vs CMS tokens

---

### Task 4: Database Transactions

**Folder**: Task 4 - Database Transactions

#### 4.1 Create User (Transaction Success) ✅

```
POST /api/auth/register
{
  "name": "Transaction Test User",
  "username": "transactiontest",
  "email": "transaction@example.com",
  "password": "SecurePass123!",
  "phone": "+628987654321"
}
```

- **Expected**: Status 201, transaction committed
- Database record created successfully

#### 4.2 Create Duplicate User (Transaction Rollback) ❌

```
POST /api/auth/register
(same data as above)
```

- **Expected**: Status 400/409, error message
- Transaction rolled back
- Database remains consistent (no partial writes)

**Verification**:

```sql
-- Check user tidak terbuat
SELECT * FROM users WHERE username = 'transactiontest';
-- Should return only 1 record (from first request)
```

#### 4.3 Update User (Transaction with Context)

```
PUT /api/users/{{user_id}}
Authorization: Bearer {{auth_token}}
{
  "name": "Updated Name",
  "phone": "+628111111111"
}
```

- **Expected**: Status 200, update successful
- Transaction supports context timeout
- Cache invalidated after update

**Transaction Features Tested**:

- ✅ ACID compliance
- ✅ Automatic rollback on error
- ✅ Context timeout support
- ✅ Constraint validation

---

### Task 5: Redis Caching

**Folder**: Task 5 - Redis Caching

#### 5.1 First Request (Cache Miss) 🔴

```
GET /api/auth/profile
Authorization: Bearer {{auth_token}}
```

- **First call**: Data dari database
- Response time disimpan untuk comparison
- Data di-cache untuk request berikutnya
- **Expected**: Status 200, slower response

#### 5.2 Second Request (Cache Hit) 🟢

```
GET /api/auth/profile
Authorization: Bearer {{auth_token}}
```

- **Second call**: Data dari Redis cache
- Response time lebih cepat
- **Test**: Cache hit faster than cache miss
- **Expected**: Status 200, faster response

**Performance Comparison**:

```
First request (DB):     15-30ms
Second request (Cache): 2-5ms
Improvement:            ~80-90%
```

#### 5.3 Create Shortlink (Cache Invalidation Test)

```
POST /api/shortenlinks
Authorization: Bearer {{auth_token}}
{
  "url": "https://example.com/test-caching",
  "custom_code": "testcache"
}
```

- **Expected**: Status 201, shortlink created
- Code cached dengan TTL 48 jam
- Code disimpan ke `{{shortlink_code}}`

#### 5.4 Get Shortlink by Code (Cache Hit)

```
GET /api/shortenlinks/{{shortlink_code}}
```

- **Expected**: Status 200, very fast response
- Data dari cache (48-hour TTL)
- No database query

#### 5.5 Update Shortlink (Cache Invalidation)

```
PUT /api/shortenlinks/{{shortlink_code}}
Authorization: Bearer {{auth_token}}
{
  "url": "https://example.com/updated-url"
}
```

- **Expected**: Status 200, update successful
- Cache invalidated (code cache + user's list cache)
- Next request will be cache miss

#### 5.6 Verify Cache TTL - User (30 min)

```
GET /api/auth/profile
```

- User cache: **30 minutes**
- Auth cache: **5 minutes** (security-sensitive)

#### 5.7 Verify Cache TTL - Shortlink (48 hours)

```
GET /api/shortenlinks/{{shortlink_code}}
```

- Shortlink code cache: **48 hours** (read-heavy)
- Shortlink data cache: **24 hours**

**Cache Strategy Verification**:

```
User Cache:          30 min  (moderate access)
User Auth Cache:     5 min   (security-sensitive)
Shortlink Data:      24 hrs  (rarely changes)
Shortlink Code:      48 hrs  (read-heavy, redirect URLs)
Permissions:         15 min  (security-sensitive)
Roles:               15 min  (security-sensitive)
```

**Cache Invalidation Patterns**:

- ✅ User update → Invalidate user + permissions + roles
- ✅ Shortlink create → Cache code for fast redirect
- ✅ Shortlink update → Invalidate code + user's list
- ✅ Shortlink delete → Invalidate code + user's list

---

### Integration Test - All Tasks Combined

**Folder**: Integration Test - All Tasks Combined

#### Full Flow Test

```
GET /api/auth/profile
Authorization: Bearer {{auth_token}}
```

**This single request validates**:

- ✅ Task 1: JWT Multi-Role Auth (token is valid)
- ✅ Task 2: RBAC Token Payload (credentials in token)
- ✅ Task 3: Middleware Identification (token extracted and validated)
- ✅ Task 4: Database Transaction (user fetched from DB if needed)
- ✅ Task 5: Redis Caching (data cached for performance)

---

## 📊 Running Tests

### Option 1: Manual Testing

1. Pilih folder/request yang ingin di-test
2. Klik **Send**
3. Lihat hasil di **Test Results** tab

### Option 2: Collection Runner

1. Klik kanan pada collection **Testing 5 Core Implementation Tasks**
2. Pilih **Run collection**
3. Klik **Run Testing 5 Core Implementation Tasks**
4. Lihat hasil semua tests

**Expected Results**:

```
Task 1 - JWT Multi-Role Auth:              3/3 tests passed ✅
Task 2 - RBAC Token Payload:               2/2 tests passed ✅
Task 3 - Middleware Token Identification:  4/4 tests passed ✅
Task 4 - Database Transactions:            3/3 tests passed ✅
Task 5 - Redis Caching:                    7/7 tests passed ✅
Integration Test:                          1/1 tests passed ✅
-----------------------------------------------------------
Total:                                    20/20 tests passed ✅
```

### Option 3: Newman (CLI)

```bash
# Install Newman
npm install -g newman

# Run collection
newman run postman/Testing-5-Tasks.postman_collection.json \
  -e postman/Development.postman_environment.json \
  --reporters cli,json

# Run with HTML report
newman run postman/Testing-5-Tasks.postman_collection.json \
  -e postman/Development.postman_environment.json \
  --reporters cli,htmlextra \
  --reporter-htmlextra-export newman-report.html
```

---

## 🔍 Test Scenarios

### Scenario 1: Happy Path (All Success)

1. Register → Login → Get Profile → Create Shortlink → Get Shortlink
2. **Expected**: All requests succeed, caching works

### Scenario 2: Authentication Failures

1. Login with wrong password → 401
2. Access protected endpoint without token → 401
3. Use invalid token format → 401
4. Use expired token → 401

### Scenario 3: Transaction Rollback

1. Create user → Success
2. Create duplicate user → Rollback, error returned
3. Verify database consistency

### Scenario 4: Cache Performance

1. First request → Cache miss (slower)
2. Second request → Cache hit (faster)
3. Update data → Cache invalidated
4. Next request → Cache miss again

### Scenario 5: Multi-Role Testing

1. Create admin user (manual via DB/seeder)
2. Login as admin
3. Verify token contains ["admin", "user"] roles
4. Access admin-only endpoints

---

## 📈 Performance Benchmarks

### Expected Response Times

| Endpoint                     | Cache Miss | Cache Hit | Improvement |
| ---------------------------- | ---------- | --------- | ----------- |
| GET /api/auth/profile        | 15-30ms    | 2-5ms     | ~80-90%     |
| GET /api/shortenlinks/{code} | 10-20ms    | 1-3ms     | ~85-95%     |
| GET /api/users               | 20-40ms    | 5-10ms    | ~75-80%     |

### Cache Hit Rates (After Warmup)

| Cache Type     | Expected Hit Rate   |
| -------------- | ------------------- |
| User Profile   | 85-95%              |
| Shortlink Code | 95-99% (read-heavy) |
| Permissions    | 90-95%              |

---

## 🐛 Troubleshooting

### Issue: Tests Failing with 401

**Solution**:

- Pastikan sudah login dan token tersimpan di `{{auth_token}}`
- Check token belum expired
- Pastikan environment **Development** aktif

### Issue: Cache Tests Showing No Improvement

**Solution**:

- Pastikan Redis running: `docker ps | grep redis`
- Check Redis connection: `redis-cli ping`
- Restart aplikasi untuk reconnect ke Redis

### Issue: Transaction Tests Failing

**Solution**:

- Check database connection
- Pastikan migrations sudah dijalankan
- Check constraint violations di database

### Issue: Response Time Inconsistent

**Solution**:

- Run request 2-3x untuk warmup
- Close other heavy applications
- Check system resources (CPU, RAM)

---

## 📝 Test Documentation

### Automated Assertions

Setiap request memiliki automated tests:

```javascript
// Example test script
pm.test('Status code is 200', function () {
  pm.response.to.have.status(200);
});

pm.test('Response has required fields', function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData.data).to.exist;
});

pm.test('Cache performance improved', function () {
  var time1 = pm.environment.get('response_time_1');
  var time2 = pm.response.responseTime;
  pm.expect(time2).to.be.below(time1);
});
```

### Console Logging

Tests menampilkan informasi di Console:

```javascript
console.log('User ID:', jsonData.data.id);
console.log('RBAC Roles:', jsonData.data.roles);
console.log('Response time:', pm.response.responseTime, 'ms');
console.log('Performance improvement:', improvement, '%');
```

---

## ✅ Success Criteria

### Task 1: JWT Multi-Role Auth

- [x] User dapat register
- [x] User dapat login dan mendapat token
- [x] Token berisi multiple roles
- [x] Token valid untuk authentication

### Task 2: RBAC Token Payload

- [x] Token payload berisi user_id
- [x] Token payload berisi username, email
- [x] Token payload berisi roles array
- [x] Token dapat distinguish user vs CMS

### Task 3: Middleware Identification

- [x] Middleware extract Bearer token
- [x] Middleware validate token signature
- [x] Middleware reject invalid format
- [x] Middleware reject expired token

### Task 4: Database Transactions

- [x] Transaction commit on success
- [x] Transaction rollback on error
- [x] Context timeout support
- [x] ACID compliance

### Task 5: Redis Caching

- [x] User data cached (30 min TTL)
- [x] Shortlink code cached (48 hrs TTL)
- [x] Cache hit faster than cache miss
- [x] Cache invalidation on update
- [x] Multi-key invalidation works

---

## 🎯 Next Steps

1. ✅ Import Postman collection
2. ✅ Import environment
3. ✅ Start server & Redis
4. ✅ Run **Task 1** tests
5. ✅ Run **Task 2** tests
6. ✅ Run **Task 3** tests
7. ✅ Run **Task 4** tests
8. ✅ Run **Task 5** tests
9. ✅ Run **Integration Test**
10. ✅ Review test results

**Happy Testing! 🚀**
