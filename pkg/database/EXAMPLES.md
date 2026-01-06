# Query Builder - Examples & Usage Guide

## 🔥 **New Advanced Features**

### 1. WHERE Operators

```go
// Basic WHERE (AND)
users := database.NewQueryBuilder("users").
    Where("age > $1", 18).
    Where("status = $1", "active").
    Execute(db)

// WHERE IN
database.NewQueryBuilder("users").
    WhereIn("role", []interface{}{"admin", "moderator", "user"}).
    Execute(db)
// SQL: WHERE role IN ($1, $2, $3)

// WHERE NOT IN
database.NewQueryBuilder("posts").
    WhereNotIn("status", []interface{}{"draft", "archived"}).
    Execute(db)

// WHERE LIKE
database.NewQueryBuilder("users").
    WhereLike("email", "%@gmail.com").
    Execute(db)

// WHERE BETWEEN
database.NewQueryBuilder("orders").
    WhereBetween("total", 100, 1000).
    Execute(db)

// WHERE NULL / NOT NULL
database.NewQueryBuilder("users").
    WhereNull("deleted_at").
    WhereNotNull("email_verified_at").
    Execute(db)
```

### 2. OR Conditions

```go
// Complex OR conditions
database.NewQueryBuilder("users").
    Where("status = $1", "active").
    OrWhere("role = 'admin'", "role = 'moderator'").
    Execute(db)
// SQL: WHERE (status = 'active') AND (role = 'admin' OR role = 'moderator')
```

### 3. JOINS

```go
// LEFT JOIN
database.NewQueryBuilder("users").
    Select("users.*, addresses.city, addresses.country").
    LeftJoin("addresses", "addresses.user_id = users.id").
    Where("users.deleted_at IS NULL").
    Execute(db)

// INNER JOIN dengan multiple kondisi
database.NewQueryBuilder("orders").
    Select("orders.*, users.name as user_name, products.name as product_name").
    InnerJoin("users", "users.id = orders.user_id").
    InnerJoin("products", "products.id = orders.product_id").
    Where("orders.status = $1", "completed").
    Execute(db)

// RIGHT JOIN
database.NewQueryBuilder("categories").
    RightJoin("products", "products.category_id = categories.id").
    Execute(db)
```

### 4. GROUP BY & HAVING

```go
// Aggregate dengan GROUP BY
database.NewQueryBuilder("orders").
    Select("user_id", "COUNT(*) as total_orders", "SUM(amount) as total_amount").
    GroupBy("user_id").
    Having("COUNT(*) > $1", 5).
    OrderBy("total_amount DESC").
    Execute(db)

// Multiple GROUP BY columns
database.NewQueryBuilder("sales").
    Select("year", "month", "COUNT(*) as count", "SUM(revenue) as revenue").
    GroupBy("year", "month").
    Having("SUM(revenue) > $1", 10000).
    Execute(db)
```

### 5. DISTINCT

```go
// Get unique cities
database.NewQueryBuilder("users").
    Distinct().
    Select("city").
    WhereNotNull("city").
    OrderBy("city ASC").
    Execute(db)
```

### 6. Pagination

```go
// Helper function untuk pagination
page := 2
perPage := 20

qb := database.NewQueryBuilder("posts").
    Where("published = $1", true).
    OrderBy("created_at DESC")

// Apply pagination
database.Paginate(qb, page, perPage)

rows, err := qb.Execute(db)
```

### 7. Aggregate Functions

```go
// Count total users
total, err := database.Count(db, "users", "status = $1", "active")
fmt.Printf("Total active users: %d\n", total)

// Check if user exists
exists, err := database.Exists(db, "users", "email = $1", "test@example.com")
if exists {
    fmt.Println("Email already registered")
}
```

### 8. Bulk Insert

```go
// Insert multiple rows efficiently
bulk := database.NewBulkInsertBuilder("products",
    []string{"name", "price", "category_id"}).
    AddRow("Product 1", 100.00, "cat-1").
    AddRow("Product 2", 200.00, "cat-1").
    AddRow("Product 3", 150.00, "cat-2")

rowsAffected, err := bulk.Execute(db)
fmt.Printf("Inserted %d products\n", rowsAffected)
```

---

## 📖 **Real-World Examples**

### Example 1: User Search with Filters

```go
func SearchUsers(db *sql.DB, filters UserFilters) ([]User, error) {
    qb := database.NewQueryBuilder("users").
        Select("id", "name", "email", "role", "created_at").
        WhereNull("deleted_at")

    // Optional filters
    if filters.Role != "" {
        qb.Where("role = $1", filters.Role)
    }

    if filters.SearchTerm != "" {
        qb.WhereLike("name", "%"+filters.SearchTerm+"%")
    }

    if len(filters.Statuses) > 0 {
        qb.WhereIn("status", filters.Statuses)
    }

    if filters.CreatedAfter != nil {
        qb.Where("created_at > $1", filters.CreatedAfter)
    }

    // Pagination
    database.Paginate(qb, filters.Page, filters.PerPage)

    rows, err := qb.Execute(db)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // ... scan results
}
```

### Example 2: Dashboard Statistics

```go
func GetDashboardStats(db *sql.DB) (*DashboardStats, error) {
    stats := &DashboardStats{}

    // Total active users
    stats.TotalUsers, _ = database.Count(db, "users", "deleted_at IS NULL AND status = $1", "active")

    // Total orders today
    today := time.Now().Format("2006-01-02")
    stats.TodayOrders, _ = database.Count(db, "orders", "DATE(created_at) = $1", today)

    // Revenue by month
    rows, err := database.NewQueryBuilder("orders").
        Select("DATE_TRUNC('month', created_at) as month", "SUM(amount) as revenue").
        Where("status = $1", "completed").
        GroupBy("month").
        OrderBy("month DESC").
        Limit(12).
        Execute(db)

    // ... process results

    return stats, nil
}
```

### Example 3: Complex Report with Joins

```go
func GetUserOrderReport(db *sql.DB, userID string) ([]OrderReport, error) {
    rows, err := database.NewQueryBuilder("orders").
        Select(`
            orders.id,
            orders.order_number,
            orders.total,
            orders.status,
            users.name as customer_name,
            COUNT(order_items.id) as total_items,
            STRING_AGG(products.name, ', ') as product_names
        `).
        InnerJoin("users", "users.id = orders.user_id").
        LeftJoin("order_items", "order_items.order_id = orders.id").
        LeftJoin("products", "products.id = order_items.product_id").
        Where("orders.user_id = $1", userID).
        WhereNull("orders.deleted_at").
        GroupBy("orders.id", "users.name").
        OrderBy("orders.created_at DESC").
        Execute(db)

    // ... process results
}
```

### Example 4: Bulk Import

```go
func ImportUsers(db *sql.DB, users []ImportUser) error {
    bulk := database.NewBulkInsertBuilder("users",
        []string{"name", "email", "phone", "role"})

    for _, user := range users {
        bulk.AddRow(user.Name, user.Email, user.Phone, "user")
    }

    _, err := bulk.Execute(db)
    return err
}
```

---

## 🎯 **Best Practices**

### ✅ DO:

```go
// Use parameterized queries
qb.Where("email = $1", email)

// Use soft delete
qb.WhereNull("deleted_at")

// Add audit trail
database.NewInsertBuilder("users").
    Set("name", name).
    SetCreatedBy(currentUserID)

// Use bulk insert for multiple rows
database.NewBulkInsertBuilder(...)
```

### ❌ DON'T:

```go
// Never concatenate user input
qb.Where("email = '" + email + "'") // SQL INJECTION!

// Don't use COUNT for existence check
count, _ := database.Count(...)
if count > 0 { ... }

// Use Exists instead
exists, _ := database.Exists(...)
if exists { ... }
```

---

## 📊 **Monitoring in Grafana**

All queries are automatically logged. Query examples in Grafana:

```logql
# All database operations
{job="kratify-backend"} | json | operation=~"SELECT|INSERT|UPDATE|DELETE"

# Slow queries (> 50ms)
{job="kratify-backend"} | json | duration > 0.05

# Failed queries
{job="kratify-backend"} | json | level="error" | message=~".*Database.*Failed"

# Queries to specific table
{job="kratify-backend"} | json | query=~".*users.*"

# Bulk operations
{job="kratify-backend"} | json | operation="BULK_INSERT"
```

---

Sekarang query builder-mu **jauh lebih powerful** bro! 🚀
