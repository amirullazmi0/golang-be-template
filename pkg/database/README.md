# Query Builder Documentation

Query builder untuk membuat SQL query dengan mudah dan aman di Kratify Backend.

## 📋 Table of Contents

-    [QueryBuilder (SELECT)](#querybuilder-select)
-    [InsertBuilder (INSERT)](#insertbuilder-insert)
-    [UpdateBuilder (UPDATE)](#updatebuilder-update)
-    [DeleteBuilder (DELETE/SOFT DELETE)](#deletebuilder-deletesoft-delete)
-    [Raw Query](#raw-query)

---

## QueryBuilder (SELECT)

Untuk membaca data dari database.

### Basic Usage

```go
// SELECT semua kolom
query, args := database.NewQueryBuilder("users").Build()
// Result: SELECT * FROM users

// SELECT kolom tertentu
query, args := database.NewQueryBuilder("users").
    Select("id", "email", "name").
    Build()
// Result: SELECT id, email, name FROM users
```

### Where Conditions

```go
// Single WHERE
query, args := database.NewQueryBuilder("users").
    Where("id = $1", userID).
    Build()
// Result: SELECT * FROM users WHERE id = $1

// Multiple WHERE (otomatis pakai AND)
query, args := database.NewQueryBuilder("users").
    Where("email = $1", email).
    Where("deleted_at IS NULL").
    Build()
// Result: SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL

// Exclude soft deleted
query, args := database.NewQueryBuilder("users").
    Where("deleted_at IS NULL").
    Build()
```

### Order By, Limit, Offset

```go
// ORDER BY
query, args := database.NewQueryBuilder("users").
    OrderBy("created_at DESC").
    Build()

// LIMIT & OFFSET (untuk pagination)
query, args := database.NewQueryBuilder("users").
    Limit(10).
    Offset(20).
    Build()
// Result: SELECT * FROM users LIMIT 10 OFFSET 20
```

### Join Tables

```go
query, args := database.NewQueryBuilder("addresses").
    Select("addresses.*", "users.name as user_name").
    Join("INNER JOIN users ON addresses.user_id = users.id").
    Where("addresses.user_id = $1", userID).
    Build()
```

### Execute Query

```go
// Method 1: Build lalu execute manual
query, args := database.NewQueryBuilder("users").
    Where("deleted_at IS NULL").
    Build()
rows, err := db.Query(query, args...)

// Method 2: Execute langsung
rows, err := database.NewQueryBuilder("users").
    Where("deleted_at IS NULL").
    Execute(db)
```

### Complete Example

```go
func (r *userRepository) FindAll() ([]model.User, error) {
    query, args := database.NewQueryBuilder("users").
        Select("id", "email", "name", "role").
        Where("deleted_at IS NULL").
        OrderBy("created_at DESC").
        Limit(100).
        Build()

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []model.User
    for rows.Next() {
        var user model.User
        err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}
```

---

## InsertBuilder (INSERT)

Untuk menambah data baru ke database.

### Basic Usage

```go
// INSERT data
id, err := database.NewInsertBuilder("users").
    Set("email", "user@example.com").
    Set("password", hashedPassword).
    Set("name", "John Doe").
    Execute(db)
// Result: INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id
```

### Dengan Audit Trail

```go
// Otomatis set created_by & created_at
id, err := database.NewInsertBuilder("addresses").
    Set("user_id", userID).
    Set("street", "Jl. Sudirman").
    Set("city", "Jakarta").
    SetCreatedBy(currentUserID). // Tambahkan ini
    Execute(db)
// Otomatis add: created_by = currentUserID, created_at = NOW()
```

### Complete Example

```go
func (r *userRepository) Create(user *model.User) (string, error) {
    id, err := database.NewInsertBuilder("users").
        Set("email", user.Email).
        Set("password", user.Password).
        Set("name", user.Name).
        Set("role", user.Role).
        Execute(r.db)

    return id, err
}
```

---

## UpdateBuilder (UPDATE)

Untuk mengubah data yang sudah ada.

### Basic Usage

```go
// UPDATE data
rowsAffected, err := database.NewUpdateBuilder("users").
    Set("name", "New Name").
    Set("email", "newemail@example.com").
    Where("id = $1", userID).
    Execute(db)
```

### Dengan Audit Trail

```go
// Otomatis set updated_by & updated_at
rowsAffected, err := database.NewUpdateBuilder("users").
    Set("name", "Updated Name").
    SetUpdatedBy(currentUserID). // Tambahkan ini
    Where("id = $1", userID).
    Execute(db)
// Otomatis add: updated_by = currentUserID, updated_at = NOW()
```

### Multiple WHERE

```go
rowsAffected, err := database.NewUpdateBuilder("addresses").
    Set("city", "Bandung").
    Where("id = $1", addressID).
    Where("user_id = $2", userID). // Security: pastikan milik user yang sama
    Execute(db)
```

### Complete Example

```go
func (r *userRepository) Update(user *model.User) error {
    _, err := database.NewUpdateBuilder("users").
        Set("name", user.Name).
        Set("email", user.Email).
        Set("updated_at", time.Now()).
        Where("id = $1", user.ID).
        Execute(r.db)

    return err
}
```

---

## DeleteBuilder (DELETE/SOFT DELETE)

Untuk menghapus data. Default adalah **soft delete** (set `deleted_at`).

### Soft Delete (Default)

```go
// Soft delete - data masih ada di DB tapi deleted_at terisi
rowsAffected, err := database.NewDeleteBuilder("users").
    Where("id = $1", userID).
    Execute(db)
// Result: UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1
```

### Soft Delete dengan Audit Trail

```go
// Otomatis set deleted_by, deleted_at, updated_at
rowsAffected, err := database.NewDeleteBuilder("users").
    Where("id = $1", userID).
    SetDeletedBy(currentUserID). // Tambahkan ini
    Execute(db)
```

### Hard Delete

```go
// Hard delete - data benar-benar dihapus dari DB
rowsAffected, err := database.NewDeleteBuilder("users").
    Where("id = $1", userID).
    HardDelete(). // Tambahkan ini untuk DELETE FROM
    Execute(db)
// Result: DELETE FROM users WHERE id = $1
```

### Complete Example

```go
func (r *userRepository) Delete(id string) error {
    // Soft delete
    _, err := database.NewDeleteBuilder("users").
        Where("id = $1", id).
        Execute(r.db)

    return err
}

func (r *addressRepository) Delete(id, userID string) error {
    // Soft delete dengan validasi ownership
    rowsAffected, err := database.NewDeleteBuilder("addresses").
        Where("id = $1", id).
        Where("user_id = $2", userID). // Security check
        SetDeletedBy(userID).
        Execute(r.db)

    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return errors.New("address not found or unauthorized")
    }

    return nil
}
```

---

## Raw Query

Untuk query kompleks yang tidak bisa dibuat dengan query builder.

### Raw SELECT

```go
rows, err := database.RawQuery(db,
    "SELECT * FROM users WHERE email LIKE $1 AND role = $2",
    "%@example.com", "ADMIN")
```

### Raw INSERT/UPDATE/DELETE

```go
result, err := database.RawExec(db,
    "UPDATE users SET is_verified = true WHERE email = $1",
    email)

rowsAffected, _ := result.RowsAffected()
```

### Raw Single Row

```go
var count int
row := database.RawQueryRow(db, "SELECT COUNT(*) FROM users WHERE role = $1", "ADMIN")
err := row.Scan(&count)
```

---

## 🔥 Best Practices

### 1. Selalu Exclude Soft Deleted

```go
// ✅ GOOD
database.NewQueryBuilder("users").
    Where("deleted_at IS NULL").
    Build()

// ❌ BAD - akan dapat data yang sudah dihapus
database.NewQueryBuilder("users").Build()
```

### 2. Gunakan Prepared Statements (Placeholder)

```go
// ✅ GOOD - aman dari SQL injection
Where("id = $1", userID)

// ❌ BAD - vulnerable ke SQL injection
Where(fmt.Sprintf("id = '%s'", userID))
```

### 3. Audit Trail untuk Data Sensitif

```go
// ✅ GOOD - track siapa yang ubah
database.NewUpdateBuilder("users").
    Set("role", "ADMIN").
    SetUpdatedBy(currentUserID).
    Where("id = $1", targetUserID).
    Execute(db)
```

### 4. Validation Ownership pada Update/Delete

```go
// ✅ GOOD - pastikan user cuma bisa hapus data sendiri
database.NewDeleteBuilder("addresses").
    Where("id = $1", addressID).
    Where("user_id = $2", currentUserID). // Security check
    Execute(db)

// ❌ BAD - user bisa hapus data user lain
database.NewDeleteBuilder("addresses").
    Where("id = $1", addressID).
    Execute(db)
```

### 5. Check RowsAffected

```go
rowsAffected, err := database.NewUpdateBuilder("users").
    Set("name", newName).
    Where("id = $1", userID).
    Execute(db)

if err != nil {
    return err
}
if rowsAffected == 0 {
    return errors.New("user not found")
}
```

---

## 📌 Common Patterns

### Pagination

```go
func (r *userRepository) FindAll(page, pageSize int) ([]model.User, error) {
    offset := (page - 1) * pageSize

    query, args := database.NewQueryBuilder("users").
        Where("deleted_at IS NULL").
        OrderBy("created_at DESC").
        Limit(pageSize).
        Offset(offset).
        Build()

    // ... execute & scan
}
```

### Search

```go
func (r *userRepository) Search(keyword string) ([]model.User, error) {
    searchPattern := "%" + keyword + "%"

    rows, err := database.RawQuery(r.db,
        `SELECT * FROM users
         WHERE deleted_at IS NULL
         AND (name LIKE $1 OR email LIKE $1)
         ORDER BY created_at DESC`,
        searchPattern)

    // ... scan
}
```

### Count

```go
func (r *userRepository) Count() (int, error) {
    var count int
    row := database.RawQueryRow(r.db,
        "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL")
    err := row.Scan(&count)
    return count, err
}
```

---

## 🎯 Tips

1. **Query Builder** cocok untuk CRUD sederhana
2. **Raw Query** untuk query kompleks (JOIN banyak, subquery, aggregation)
3. Selalu gunakan `Where("deleted_at IS NULL")` untuk exclude soft deleted data
4. Gunakan `SetCreatedBy`, `SetUpdatedBy`, `SetDeletedBy` untuk audit trail
5. Jangan lupa `defer rows.Close()` setelah query
6. Check `rowsAffected` untuk validasi update/delete berhasil

---

**Happy coding! 🚀**
