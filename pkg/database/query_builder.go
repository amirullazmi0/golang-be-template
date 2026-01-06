package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/amirullazmi0/kratify-backend/pkg/logger"
	"go.uber.org/zap"
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
	orWhere    [][]string // For OR conditions, each element is a group of AND conditions
	whereArgs  []interface{}
	orderBy    []string
	groupBy    []string
	having     []string
	havingArgs []interface{}
	limit      int
	offset     int
	join       []string
	distinct   bool
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{
		table:      table,
		columns:    []string{"*"},
		where:      []string{},
		orWhere:    [][]string{},
		whereArgs:  []interface{}{},
		orderBy:    []string{},
		groupBy:    []string{},
		having:     []string{},
		havingArgs: []interface{}{},
		join:       []string{},
		distinct:   false,
	}
}

// Select specifies columns to select
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.columns = columns
	return qb
}

// Distinct adds DISTINCT to SELECT
func (qb *QueryBuilder) Distinct() *QueryBuilder {
	qb.distinct = true
	return qb
}

// Where adds a WHERE condition (AND)
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, args...)
	return qb
}

// WhereIn adds WHERE column IN (values)
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", len(qb.whereArgs)+i+1)
	}
	condition := fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", "))
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, values...)
	return qb
}

// WhereNotIn adds WHERE column NOT IN (values)
func (qb *QueryBuilder) WhereNotIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", len(qb.whereArgs)+i+1)
	}
	condition := fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ", "))
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, values...)
	return qb
}

// WhereLike adds WHERE column LIKE pattern
func (qb *QueryBuilder) WhereLike(column string, pattern string) *QueryBuilder {
	condition := fmt.Sprintf("%s LIKE $%d", column, len(qb.whereArgs)+1)
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, pattern)
	return qb
}

// WhereBetween adds WHERE column BETWEEN start AND end
func (qb *QueryBuilder) WhereBetween(column string, start, end interface{}) *QueryBuilder {
	condition := fmt.Sprintf("%s BETWEEN $%d AND $%d", column, len(qb.whereArgs)+1, len(qb.whereArgs)+2)
	qb.where = append(qb.where, condition)
	qb.whereArgs = append(qb.whereArgs, start, end)
	return qb
}

// WhereNull adds WHERE column IS NULL
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s IS NULL", column))
	return qb
}

// WhereNotNull adds WHERE column IS NOT NULL
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s IS NOT NULL", column))
	return qb
}

// OrWhere adds OR WHERE conditions
func (qb *QueryBuilder) OrWhere(conditions ...string) *QueryBuilder {
	if len(conditions) > 0 {
		qb.orWhere = append(qb.orWhere, conditions)
	}
	return qb
}

// Join adds a JOIN clause
func (qb *QueryBuilder) Join(joinClause string) *QueryBuilder {
	qb.join = append(qb.join, joinClause)
	return qb
}

// LeftJoin adds a LEFT JOIN
func (qb *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	qb.join = append(qb.join, fmt.Sprintf("LEFT JOIN %s ON %s", table, condition))
	return qb
}

// RightJoin adds a RIGHT JOIN
func (qb *QueryBuilder) RightJoin(table, condition string) *QueryBuilder {
	qb.join = append(qb.join, fmt.Sprintf("RIGHT JOIN %s ON %s", table, condition))
	return qb
}

// InnerJoin adds an INNER JOIN
func (qb *QueryBuilder) InnerJoin(table, condition string) *QueryBuilder {
	qb.join = append(qb.join, fmt.Sprintf("INNER JOIN %s ON %s", table, condition))
	return qb
}

// GroupBy adds GROUP BY clause
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, columns...)
	return qb
}

// Having adds HAVING clause
func (qb *QueryBuilder) Having(condition string, args ...interface{}) *QueryBuilder {
	qb.having = append(qb.having, condition)
	qb.havingArgs = append(qb.havingArgs, args...)
	return qb
}

// OrderBy sets ORDER BY clause
func (qb *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	qb.orderBy = append(qb.orderBy, orderBy)
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
	selectClause := "SELECT"
	if qb.distinct {
		selectClause = "SELECT DISTINCT"
	}

	query := fmt.Sprintf("%s %s FROM %s", selectClause, strings.Join(qb.columns, ", "), qb.table)

	if len(qb.join) > 0 {
		query += " " + strings.Join(qb.join, " ")
	}

	allArgs := qb.whereArgs

	if len(qb.where) > 0 || len(qb.orWhere) > 0 {
		query += " WHERE "
		conditions := []string{}

		// Add AND conditions
		if len(qb.where) > 0 {
			conditions = append(conditions, "("+strings.Join(qb.where, " AND ")+")")
		}

		// Add OR groups
		for _, orGroup := range qb.orWhere {
			if len(orGroup) > 0 {
				conditions = append(conditions, "("+strings.Join(orGroup, " OR ")+")")
			}
		}

		query += strings.Join(conditions, " AND ")
	}

	if len(qb.groupBy) > 0 {
		query += " GROUP BY " + strings.Join(qb.groupBy, ", ")
	}

	if len(qb.having) > 0 {
		query += " HAVING " + strings.Join(qb.having, " AND ")
		allArgs = append(allArgs, qb.havingArgs...)
	}

	if len(qb.orderBy) > 0 {
		query += " ORDER BY " + strings.Join(qb.orderBy, ", ")
	}

	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return query, allArgs
}

// BuildResult builds query and returns QueryResult
func (qb *QueryBuilder) BuildResult() QueryResult {
	query, args := qb.Build()
	return QueryResult{Query: query, Args: args}
}

// Execute executes the query and returns rows
func (qb *QueryBuilder) Execute(db *sql.DB) (*sql.Rows, error) {
	query, args := qb.Build()
	start := time.Now()
	rows, err := db.Query(query, args...)
	duration := time.Since(start)

	// Log query execution
	logFields := []zap.Field{
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.String("operation", "SELECT"),
	}

	if err != nil {
		logger.Error("Database Query Failed", append(logFields, zap.Error(err))...)
	} else {
		logger.Info("Database Query", logFields...)
	}

	return rows, err
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
	start := time.Now()
	var id string
	err := db.QueryRow(query, args...).Scan(&id)
	duration := time.Since(start)

	// Log query execution
	logFields := []zap.Field{
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.String("operation", "INSERT"),
	}

	if err != nil {
		logger.Error("Database Insert Failed", append(logFields, zap.Error(err))...)
	} else {
		logger.Info("Database Insert", append(logFields, zap.String("id", id))...)
	}

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
	start := time.Now()
	result, err := db.Exec(query, args...)
	duration := time.Since(start)

	var rowsAffected int64
	if err == nil {
		rowsAffected, _ = result.RowsAffected()
	}

	// Log query execution
	logFields := []zap.Field{
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.String("operation", "UPDATE"),
		zap.Int64("rows_affected", rowsAffected),
	}

	if err != nil {
		logger.Error("Database Update Failed", append(logFields, zap.Error(err))...)
		return 0, err
	}

	logger.Info("Database Update", logFields...)
	return rowsAffected, nil
}

// DeleteBuilder builds soft DELETE queries (UPDATE deleted_at)
type DeleteBuilder struct {
	table      string
	where      []string
	whereArgs  []interface{}
	deletedBy  *string // User UUID who deleted
	hardDelete bool    // True = hard delete, False = soft delete
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
	start := time.Now()
	result, err := sqlDB.Exec(query, args...)
	duration := time.Since(start)

	var rowsAffected int64
	if err == nil {
		rowsAffected, _ = result.RowsAffected()
	}

	// Log query execution
	logFields := []zap.Field{
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.String("operation", "DELETE"),
		zap.Int64("rows_affected", rowsAffected),
		zap.Bool("hard_delete", db.hardDelete),
	}

	if err != nil {
		logger.Error("Database Delete Failed", append(logFields, zap.Error(err))...)
		return 0, err
	}

	logger.Info("Database Delete", logFields...)
	return rowsAffected, nil
}

// RawQuery executes a raw SQL query
func RawQuery(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.Query(query, args...)
	duration := time.Since(start)

	// Log query execution
	logFields := []zap.Field{
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.String("operation", "RAW_QUERY"),
	}

	if err != nil {
		logger.Error("Database Raw Query Failed", append(logFields, zap.Error(err))...)
	} else {
		logger.Info("Database Raw Query", logFields...)
	}

	return rows, err
}

// RawExec executes a raw SQL command
func RawExec(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := db.Exec(query, args...)
	duration := time.Since(start)

	var rowsAffected int64
	if err == nil {
		rowsAffected, _ = result.RowsAffected()
	}

	// Log query execution
	logFields := []zap.Field{
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.String("operation", "RAW_EXEC"),
		zap.Int64("rows_affected", rowsAffected),
	}

	if err != nil {
		logger.Error("Database Raw Exec Failed", append(logFields, zap.Error(err))...)
	} else {
		logger.Info("Database Raw Exec", logFields...)
	}

	return result, err
}

// RawQueryRow executes a raw SQL query for single row
func RawQueryRow(db *sql.DB, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := db.QueryRow(query, args...)
	duration := time.Since(start)

	// Log query execution
	logger.Info("Database Raw QueryRow",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.String("operation", "RAW_QUERY_ROW"),
	)

	return row
}

// Helper functions for common aggregate queries

// Count returns count of rows
func Count(db *sql.DB, table string, where string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if where != "" {
		query += " WHERE " + where
	}
	
	start := time.Now()
	var count int64
	err := db.QueryRow(query, args...).Scan(&count)
	duration := time.Since(start)
	
	logger.Info("Database Count",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
		zap.Int64("count", count),
	)
	
	return count, err
}

// Exists checks if rows exist
func Exists(db *sql.DB, table string, where string, args ...interface{}) (bool, error) {
	count, err := Count(db, table, where, args...)
	return count > 0, err
}

// Paginate helper for pagination
type PaginationResult struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func Paginate(qb *QueryBuilder, page, perPage int) *QueryBuilder {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	
	offset := (page - 1) * perPage
	return qb.Limit(perPage).Offset(offset)
}

// BulkInsert for inserting multiple rows
type BulkInsertBuilder struct {
	table   string
	columns []string
	rows    [][]interface{}
}

func NewBulkInsertBuilder(table string, columns []string) *BulkInsertBuilder {
	return &BulkInsertBuilder{
		table:   table,
		columns: columns,
		rows:    [][]interface{}{},
	}
}

func (bib *BulkInsertBuilder) AddRow(values ...interface{}) *BulkInsertBuilder {
	if len(values) != len(bib.columns) {
		logger.Error("Bulk insert: column count mismatch",
			zap.Int("expected", len(bib.columns)),
			zap.Int("got", len(values)),
		)
		return bib
	}
	bib.rows = append(bib.rows, values)
	return bib
}

func (bib *BulkInsertBuilder) Execute(db *sql.DB) (int64, error) {
	if len(bib.rows) == 0 {
		return 0, fmt.Errorf("no rows to insert")
	}
	
	// Build placeholders
	valuePlaceholders := []string{}
	allValues := []interface{}{}
	placeholder := 1
	
	for _, row := range bib.rows {
		rowPlaceholders := []string{}
		for range row {
			rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", placeholder))
			placeholder++
		}
		valuePlaceholders = append(valuePlaceholders, "("+strings.Join(rowPlaceholders, ", ")+")")
		allValues = append(allValues, row...)
	}
	
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		bib.table,
		strings.Join(bib.columns, ", "),
		strings.Join(valuePlaceholders, ", "),
	)
	
	start := time.Now()
	result, err := db.Exec(query, allValues...)
	duration := time.Since(start)
	
	var rowsAffected int64
	if err == nil {
		rowsAffected, _ = result.RowsAffected()
	}
	
	logger.Info("Database Bulk Insert",
		zap.String("query", query),
		zap.Int("rows", len(bib.rows)),
		zap.Duration("duration", duration),
		zap.Int64("rows_affected", rowsAffected),
	)
	
	return rowsAffected, err
}

