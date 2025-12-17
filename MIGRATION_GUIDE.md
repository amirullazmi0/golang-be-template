# Kratify Backend - Prisma Migration + Custom Query Builder

## 🎯 Architecture

-    **Migration**: Prisma (type-safe, foreign keys, rollback support)
-    **Query**: Custom Query Builder (`database/sql`)

## 🚀 Cara Migration

### 1. Create Database (One Time)

Buat database `dbkratify` di DBeaver atau psql:

```sql
CREATE DATABASE dbkratify;
```

### 2. Edit Schema

Edit file `prisma/schema.prisma`:

```prisma
model YourModel {
  id   Int    @id @default(autoincrement())
  name String
}
```

### 3. Generate Migration

```bash
npm run migrate:dev
# Masukkan nama migration misal: "add_users_table"
```

### 4. Apply Migration (Production)

```bash
npm run migrate:deploy
```

### 5. View Data

```bash
npm run studio
```

## 📝 Cara Pakai Query Builder

### SELECT

```go
import "github.com/amirullazmi0/kratify-backend/pkg/database"

// Simple select
query := database.NewQueryBuilder("users").
    Select("id", "email", "name").
    Where("email = $1", "user@example.com").
    OrderBy("created_at DESC").
    Limit(10)

rows, err := query.Execute(db.DB)

// With JOIN
query := database.NewQueryBuilder("products").
    Select("products.*, users.name as user_name").
    Join("INNER JOIN users ON products.user_id = users.id").
    Where("products.price > $1", 100).
    Build()
```

### INSERT

```go
id, err := database.NewInsertBuilder("users").
    Set("email", "user@example.com").
    Set("password", hashedPassword).
    Set("name", "John Doe").
    Execute(db.DB)
```

### UPDATE

```go
rowsAffected, err := database.NewUpdateBuilder("users").
    Set("name", "Jane Doe").
    Set("updated_at", time.Now()).
    Where("id = $1", userID).
    Execute(db.DB)
```

### DELETE

```go
rowsAffected, err := database.NewDeleteBuilder("users").
    Where("id = $1", userID).
    Execute(db.DB)
```

## 🔥 Example: User Repository dengan Query Builder

```go
package repository

import (
    "database/sql"
    "github.com/amirullazmi0/kratify-backend/pkg/database"
    "github.com/amirullazmi0/kratify-backend/internal/model"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    query, args := database.NewQueryBuilder("users").
        Where("email = $1", email).
        Where("deleted_at IS NULL").
        Limit(1).
        Build()

    var user model.User
    err := r.db.QueryRow(query, args...).Scan(
        &user.ID,
        &user.Email,
        &user.Password,
        &user.Name,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    return &user, err
}

func (r *UserRepository) Create(user *model.User) (int64, error) {
    return database.NewInsertBuilder("users").
        Set("email", user.Email).
        Set("password", user.Password).
        Set("name", user.Name).
        Execute(r.db)
}

func (r *UserRepository) Update(user *model.User) error {
    _, err := database.NewUpdateBuilder("users").
        Set("name", user.Name).
        Set("updated_at", time.Now()).
        Where("id = $1", user.ID).
        Execute(r.db)
    return err
}
```

## 📋 Prisma Commands

| Command                  | Deskripsi                               |
| ------------------------ | --------------------------------------- |
| `npm run migrate:dev`    | Create & apply migration (dev)          |
| `npm run migrate:deploy` | Apply pending migrations (prod)         |
| `npm run migrate:reset`  | Reset database & rerun all migrations   |
| `npm run db:push`        | Push schema tanpa create migration file |
| `npm run db:pull`        | Pull schema dari database existing      |
| `npm run studio`         | Open Prisma Studio (GUI database)       |

## 🎨 Benefits

✅ **Prisma Migration**:

-    Type-safe schema
-    Auto foreign keys & indexes
-    Migration history
-    Easy rollback
-    Team collaboration

✅ **Custom Query Builder**:

-    Full control
-    No ORM overhead
-    Raw SQL capability
-    Lightweight
-    Fast performance

## 🔄 Workflow

1. **Schema Changes**: Edit `prisma/schema.prisma`
2. **Generate Migration**: `npm run migrate:dev`
3. **Update Go Structs**: Manual sync dengan schema
4. **Write Queries**: Pakai query builder di repository
5. **Deploy**: `npm run migrate:deploy`
