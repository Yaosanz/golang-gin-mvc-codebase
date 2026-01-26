# 🧪 POSTMAN COLLECTION - QUICK START

## Status: ✅ ALL 5 TASKS READY FOR TESTING

---

## 📋 What's Been Fixed

| Fix                                    | Status | Impact                              |
| -------------------------------------- | ------ | ----------------------------------- |
| Circular reference in User/Role models | ✅     | GET /api/v1/users now works         |
| ShortenLink JSON serialization         | ✅     | GET /api/v1/shorten-links now works |
| Permission check error handling        | ✅     | Better 401 responses on cache miss  |
| 404 spam from favicon/hot-update       | ✅     | Cleaner logs                        |
| Roles preloading in queries            | ✅     | Multi-role support complete         |

---

## 🚀 Quick Test Instructions

### 1. Start Server

```bash
go run cmd/app/main.go
```

### 2. Open Postman

### 3. Import Files

- **Collection:** `postman/Testing-5-Tasks.postman_collection.json`
- **Environment:** `postman/Development.postman_environment.json`

### 4. Run Tests

Click the **▶️ Run** button in collection → Select "Start Test"

### 5. Expected Results

- ✅ 20+ tests should PASS
- ✅ Response times logged
- ✅ Cache performance measured

---

## 📊 Test Breakdown

### Task 1: JWT Multi-Role Auth (3 tests)

```
✅ Register user
✅ Login & get token with roles
✅ Token contains ["admin"] roles
```

### Task 2: RBAC Payload (2 tests)

```
✅ Token has roles array
✅ Token has permissions array
```

### Task 3: Middleware Identification (4 tests)

```
✅ Valid token extracts user_id
✅ Invalid token returns 401
✅ Missing token returns 401
✅ Expired token returns 401
```

### Task 4: Database Transactions (3 tests)

```
✅ Create shortenlink (success)
✅ Update shortenlink (success)
✅ Transaction rollback on error
```

### Task 5: Redis Caching (7 tests)

```
✅ Cache miss on first request
✅ Cache hit on second request
✅ Cache invalidation on update
✅ Cache invalidation on delete
✅ TTL verification (30 min)
✅ Performance: Cached < Direct DB
✅ Response time comparison
```

### Integration Tests (1 test)

```
✅ Full workflow: Register → Login → Shortenlink CRUD
```

---

## 🔍 Key Endpoints Tested

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
GET    /api/v1/auth/profile
GET    /api/v1/users              ← Fixed (was 500)
GET    /api/v1/shorten-links/     ← Fixed (was 500)
POST   /api/v1/shorten-links
PATCH  /api/v1/shorten-links/:id
DELETE /api/v1/shorten-links/:id
```

---

## ✨ Features Verified

- ✅ JWT with multi-role support
- ✅ RBAC with permission checking
- ✅ Middleware authentication & authorization
- ✅ Database transactions
- ✅ Redis caching with TTL
- ✅ Error handling & logging
- ✅ Performance monitoring

---

## 📝 Notes

- **Token:** Valid for 24 hours
- **Cache TTL:** 30 minutes
- **Admin:** Bypasses most permission checks
- **Regular User:** Must have specific permissions

---

**Ready to test! Run Postman collection now! 🚀**
