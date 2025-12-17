package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// QueryResult holds query string and arguments
type QueryResult struct {
	Query string
	Args  []interface{}
}

// QueryBuilder is a simple SQL query builder
type QueryBuilder struct {
	table      string
	columns    []string
	where      []string
	whereArgs  []interface{}
	orderBy    string
	limit      int
	offset     int
	join       []string
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table:     table,
		columns:   []string{"*"},
		where:     []string{},
		whereArgs: []interface{}{},
		join:      []string{},
	}
}

// Select specifies columns to select
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.columns = columns
	return qb
}

// Where adds a WHERE condition
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, args...)
	return qb
}

// Join adds a JOIN clause
func (qb *QueryBuilder) Join(joinClause string) *QueryBuilder {
	qb.join = append(qb.join, joinClause)
	return qb
}

// OrderBy sets ORDER BY clause
func (qb *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	qb.orderBy = orderBy
	return qb
}

// Limit sets LIMIT
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset sets OFFSET
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Build builds the SELECT query
func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(qb.columns, ", "), qb.table)

	if len(qb.join) > 0 {
		query += " " + strings.Join(qb.join, " ")
	}

	if len(qb.where) > 0 {
		query += " WHERE " + strings.Join(qb.where, " AND ")
	}

	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query, qb.whereArgs
}

// BuildResult builds query and returns QueryResult
func (qb *QueryBuilder) BuildResult() QueryResult {
	query, args := qb.Build()
	return QueryResult{Query: query, Args: args}
}

// Execute executes the query and returns rows
func (qb *QueryBuilder) Execute(db *sql.DB) (*sql.Rows, error) {
	query, args := qb.Build()
	return db.Query(query, args...)
}

// InsertBuilder builds INSERT queries
type InsertBuilder struct {
	table     string
	columns   []string
	values    []interface{}
	createdBy *string // User UUID who created
}

// NewInsertBuilder creates a new insert builder
func NewInsertBuilder(table string) *InsertBuilder {
	return &InsertBuilder{
		table:   table,
		columns: []string{},
		values:  []interface{}{},
	}
}

// Set adds a column and value
func (ib *InsertBuilder) Set(column string, value interface{}) *InsertBuilder {
	ib.columns = append(ib.columns, column)
	ib.values = append(ib.values, value)
	return ib
}

// SetCreatedBy sets created_by and created_at automatically
func (ib *InsertBuilder) SetCreatedBy(userID string) *InsertBuilder {
	ib.createdBy = &userID
	return ib
}

// Build builds the INSERT query
func (ib *InsertBuilder) Build() (string, []interface{}) {
	// Add audit fields if createdBy is set
	if ib.createdBy != nil {
		ib.columns = append(ib.columns, "created_by", "created_at")
		ib.values = append(ib.values, *ib.createdBy, time.Now())
	}

	placeholders := make([]string, len(ib.columns))
	for i := range ib.columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		ib.table,
		strings.Join(ib.columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return query, ib.values
}

// BuildResult builds query and returns QueryResult
func (ib *InsertBuilder) BuildResult() QueryResult {
	query, args := ib.Build()
	return QueryResult{Query: query, Args: args}
}

// Execute executes the insert query and returns UUID
func (ib *InsertBuilder) Execute(db *sql.DB) (string, error) {
	query, args := ib.Build()
	var id string
	err := db.QueryRow(query, args...).Scan(&id)
	return id, err
}

// UpdateBuilder builds UPDATE queries
type UpdateBuilder struct {
	table     string
	sets      []string
	setArgs   []interface{}
	where     []string
	whereArgs []interface{}
	updatedBy *string // User UUID who updated
}

// NewUpdateBuilder creates a new update builder
func NewUpdateBuilder(table string) *UpdateBuilder {
	return &UpdateBuilder{
		table:     table,
		sets:      []string{},
		setArgs:   []interface{}{},
		where:     []string{},
		whereArgs: []interface{}{},
	}
}

// Set adds a column to update
func (ub *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	ub.sets = append(ub.sets, fmt.Sprintf("%s = $%d", column, len(ub.setArgs)+1))
	ub.setArgs = append(ub.setArgs, value)
	return ub
}

// SetUpdatedBy sets updated_by and updated_at automatically
func (ub *UpdateBuilder) SetUpdatedBy(userID string) *UpdateBuilder {
	ub.updatedBy = &userID
	return ub
}

// Where adds a WHERE condition
func (ub *UpdateBuilder) Where(condition string, args ...interface{}) *UpdateBuilder {
	// Adjust placeholder numbering
	paramOffset := len(ub.setArgs)
	for i := range args {
		condition = strings.Replace(condition, fmt.Sprintf("$%d", i+1), fmt.Sprintf("$%d", paramOffset+i+1), 1)
	}
	ub.where = append(ub.where, condition)
	ub.whereArgs = append(ub.whereArgs, args...)
	return ub
}

// Build builds the UPDATE query
func (ub *UpdateBuilder) Build() (string, []interface{}) {
	// Add audit fields if updatedBy is set
	if ub.updatedBy != nil {
		ub.sets = append(ub.sets, fmt.Sprintf("updated_by = $%d", len(ub.setArgs)+1))
		ub.setArgs = append(ub.setArgs, *ub.updatedBy)
		ub.sets = append(ub.sets, fmt.Sprintf("updated_at = $%d", len(ub.setArgs)+1))
		ub.setArgs = append(ub.setArgs, time.Now())
	}

	query := fmt.Sprintf("UPDATE %s SET %s", ub.table, strings.Join(ub.sets, ", "))

	if len(ub.where) > 0 {
		query += " WHERE " + strings.Join(ub.where, " AND ")
	}

	allArgs := append(ub.setArgs, ub.whereArgs...)
	return query, allArgs
}

// BuildResult builds query and returns QueryResult
func (ub *UpdateBuilder) BuildResult() QueryResult {
	query, args := ub.Build()
	return QueryResult{Query: query, Args: args}
}

// Execute executes the update query
func (ub *UpdateBuilder) Execute(db *sql.DB) (int64, error) {
	query, args := ub.Build()
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// DeleteBuilder builds soft DELETE queries (UPDATE deleted_at)
type DeleteBuilder struct {
	table     string
	where     []string
	whereArgs []interface{}
	deletedBy *string // User UUID who deleted
	hardDelete bool   // True = hard delete, False = soft delete
}

// NewDeleteBuilder creates a new delete builder (soft delete by default)
func NewDeleteBuilder(table string) *DeleteBuilder {
	return &DeleteBuilder{
		table:      table,
		where:      []string{},
		whereArgs:  []interface{}{},
		hardDelete: false,
	}
}

// Where adds a WHERE condition
func (db *DeleteBuilder) Where(condition string, args ...interface{}) *DeleteBuilder {
	db.where = append(db.where, condition)
	db.whereArgs = append(db.whereArgs, args...)
	return db
}

// SetDeletedBy sets deleted_by and deleted_at for soft delete
func (db *DeleteBuilder) SetDeletedBy(userID string) *DeleteBuilder {
	db.deletedBy = &userID
	return db
}

// HardDelete sets mode to hard delete (actual DELETE FROM)
func (db *DeleteBuilder) HardDelete() *DeleteBuilder {
	db.hardDelete = true
	return db
}

// Build builds the DELETE or UPDATE query
func (db *DeleteBuilder) Build() (string, []interface{}) {
	if db.hardDelete {
		// Hard delete - actual DELETE FROM
		query := fmt.Sprintf("DELETE FROM %s", db.table)
		if len(db.where) > 0 {
			query += " WHERE " + strings.Join(db.where, " AND ")
		}
		return query, db.whereArgs
	}

	// Soft delete - UPDATE deleted_at
	sets := []string{
		"deleted_at = $1",
	}
	args := []interface{}{time.Now()}

	if db.deletedBy != nil {
		sets = append(sets, "deleted_by = $2")
		args = append(args, *db.deletedBy)
		sets = append(sets, "updated_at = $3")
		args = append(args, time.Now())
	} else {
		sets = append(sets, "updated_at = $2")
		args = append(args, time.Now())
	}

	query := fmt.Sprintf("UPDATE %s SET %s", db.table, strings.Join(sets, ", "))

	if len(db.where) > 0 {
		// Adjust placeholders for WHERE clause
		offset := len(args)
		adjustedWhere := make([]string, len(db.where))
		for i, w := range db.where {
			adjustedWhere[i] = w
			// Replace $1, $2, etc with offset values
			for j := 1; j <= len(db.whereArgs); j++ {
				oldPlaceholder := fmt.Sprintf("$%d", j)
				newPlaceholder := fmt.Sprintf("$%d", offset+j)
				adjustedWhere[i] = strings.Replace(adjustedWhere[i], oldPlaceholder, newPlaceholder, 1)
			}
		}
		query += " WHERE " + strings.Join(adjustedWhere, " AND ")
		args = append(args, db.whereArgs...)
	}

	return query, args
}

// BuildResult builds query and returns QueryResult
func (db *DeleteBuilder) BuildResult() QueryResult {
	query, args := db.Build()
	return QueryResult{Query: query, Args: args}
}

// Execute executes the delete query
func (db *DeleteBuilder) Execute(sqlDB *sql.DB) (int64, error) {
	query, args := db.Build()
	result, err := sqlDB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// RawQuery executes a raw SQL query
func RawQuery(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

// RawExec executes a raw SQL command
func RawExec(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

// RawQueryRow executes a raw SQL query for single row
func RawQueryRow(db *sql.DB, query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}
