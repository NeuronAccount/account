package user_db

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/sql/wrap"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

var _ = sql.ErrNoRows
var _ = mysql.ErrOldProtocol

type BaseQuery struct {
	forUpdate     bool
	forShare      bool
	where         string
	limit         string
	order         string
	groupByFields []string
}

func (q *BaseQuery) buildQueryString() string {
	buf := bytes.NewBufferString("")

	if q.where != "" {
		buf.WriteString(" WHERE ")
		buf.WriteString(q.where)
	}

	if q.groupByFields != nil && len(q.groupByFields) > 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(strings.Join(q.groupByFields, ","))
	}

	if q.order != "" {
		buf.WriteString(" order by ")
		buf.WriteString(q.order)
	}

	if q.limit != "" {
		buf.WriteString(q.limit)
	}

	if q.forUpdate {
		buf.WriteString(" FOR UPDATE ")
	}

	if q.forShare {
		buf.WriteString(" LOCK IN SHARE MODE ")
	}

	return buf.String()
}

const ACCESS_TOKEN_TABLE_NAME = "access_token"

type ACCESS_TOKEN_FIELD string

const ACCESS_TOKEN_FIELD_ID = ACCESS_TOKEN_FIELD("id")
const ACCESS_TOKEN_FIELD_USER_ID = ACCESS_TOKEN_FIELD("user_id")
const ACCESS_TOKEN_FIELD_ACCESS_TOKEN = ACCESS_TOKEN_FIELD("access_token")
const ACCESS_TOKEN_FIELD_CREATE_TIME = ACCESS_TOKEN_FIELD("create_time")
const ACCESS_TOKEN_FIELD_UPDATE_TIME = ACCESS_TOKEN_FIELD("update_time")

const ACCESS_TOKEN_ALL_FIELDS_STRING = "id,user_id,access_token,create_time,update_time"

var ACCESS_TOKEN_ALL_FIELDS = []string{
	"id",
	"user_id",
	"access_token",
	"create_time",
	"update_time",
}

type AccessToken struct {
	Id          uint64 //size=20
	UserId      string //size=32
	AccessToken string //size=1024
	CreateTime  time.Time
	UpdateTime  time.Time
}

type AccessTokenQuery struct {
	BaseQuery
	dao *AccessTokenDao
}

func NewAccessTokenQuery(dao *AccessTokenDao) *AccessTokenQuery {
	q := &AccessTokenQuery{}
	q.dao = dao

	return q
}

func (q *AccessTokenQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*AccessToken, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *AccessTokenQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*AccessToken, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *AccessTokenQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *AccessTokenQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *AccessTokenQuery) ForUpdate() *AccessTokenQuery {
	q.forUpdate = true
	return q
}

func (q *AccessTokenQuery) ForShare() *AccessTokenQuery {
	q.forShare = true
	return q
}

func (q *AccessTokenQuery) GroupBy(fields ...ACCESS_TOKEN_FIELD) *AccessTokenQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *AccessTokenQuery) Limit(startIncluded int64, count int64) *AccessTokenQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *AccessTokenQuery) OrderBy(fieldName ACCESS_TOKEN_FIELD, asc bool) *AccessTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *AccessTokenQuery) OrderByGroupCount(asc bool) *AccessTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *AccessTokenQuery) w(format string, a ...interface{}) *AccessTokenQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *AccessTokenQuery) Left() *AccessTokenQuery  { return q.w(" ( ") }
func (q *AccessTokenQuery) Right() *AccessTokenQuery { return q.w(" ) ") }
func (q *AccessTokenQuery) And() *AccessTokenQuery   { return q.w(" AND ") }
func (q *AccessTokenQuery) Or() *AccessTokenQuery    { return q.w(" OR ") }
func (q *AccessTokenQuery) Not() *AccessTokenQuery   { return q.w(" NOT ") }

func (q *AccessTokenQuery) Id_Equal(v uint64) *AccessTokenQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_NotEqual(v uint64) *AccessTokenQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_Less(v uint64) *AccessTokenQuery { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *AccessTokenQuery) Id_LessEqual(v uint64) *AccessTokenQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_Greater(v uint64) *AccessTokenQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_GreaterEqual(v uint64) *AccessTokenQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UserId_Equal(v string) *AccessTokenQuery {
	return q.w("user_id='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UserId_NotEqual(v string) *AccessTokenQuery {
	return q.w("user_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UserId_Less(v string) *AccessTokenQuery {
	return q.w("user_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UserId_LessEqual(v string) *AccessTokenQuery {
	return q.w("user_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UserId_Greater(v string) *AccessTokenQuery {
	return q.w("user_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UserId_GreaterEqual(v string) *AccessTokenQuery {
	return q.w("user_id>='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccessToken_Equal(v string) *AccessTokenQuery {
	return q.w("access_token='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccessToken_NotEqual(v string) *AccessTokenQuery {
	return q.w("access_token<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccessToken_Less(v string) *AccessTokenQuery {
	return q.w("access_token<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccessToken_LessEqual(v string) *AccessTokenQuery {
	return q.w("access_token<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccessToken_Greater(v string) *AccessTokenQuery {
	return q.w("access_token>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccessToken_GreaterEqual(v string) *AccessTokenQuery {
	return q.w("access_token>='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) CreateTime_Equal(v time.Time) *AccessTokenQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) CreateTime_NotEqual(v time.Time) *AccessTokenQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) CreateTime_Less(v time.Time) *AccessTokenQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) CreateTime_LessEqual(v time.Time) *AccessTokenQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) CreateTime_Greater(v time.Time) *AccessTokenQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) CreateTime_GreaterEqual(v time.Time) *AccessTokenQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UpdateTime_Equal(v time.Time) *AccessTokenQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UpdateTime_NotEqual(v time.Time) *AccessTokenQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UpdateTime_Less(v time.Time) *AccessTokenQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UpdateTime_LessEqual(v time.Time) *AccessTokenQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UpdateTime_Greater(v time.Time) *AccessTokenQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) UpdateTime_GreaterEqual(v time.Time) *AccessTokenQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type AccessTokenDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewAccessTokenDao(db *DB) (t *AccessTokenDao, err error) {
	t = &AccessTokenDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *AccessTokenDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccessTokenDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO access_token (user_id,access_token) VALUES (?,?)")
	return err
}

func (dao *AccessTokenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE access_token SET user_id=?,access_token=? WHERE id=?")
	return err
}

func (dao *AccessTokenDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM access_token WHERE id=?")
	return err
}

func (dao *AccessTokenDao) Insert(ctx context.Context, tx *wrap.Tx, e *AccessToken) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.UserId, e.AccessToken)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *AccessTokenDao) Update(ctx context.Context, tx *wrap.Tx, e *AccessToken) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.UserId, e.AccessToken, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccessTokenDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccessTokenDao) scanRow(row *wrap.Row) (*AccessToken, error) {
	e := &AccessToken{}
	err := row.Scan(&e.Id, &e.UserId, &e.AccessToken, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *AccessTokenDao) scanRows(rows *wrap.Rows) (list []*AccessToken, err error) {
	list = make([]*AccessToken, 0)
	for rows.Next() {
		e := AccessToken{}
		err = rows.Scan(&e.Id, &e.UserId, &e.AccessToken, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *AccessTokenDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*AccessToken, error) {
	querySql := "SELECT " + ACCESS_TOKEN_ALL_FIELDS_STRING + " FROM access_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *AccessTokenDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*AccessToken, err error) {
	querySql := "SELECT " + ACCESS_TOKEN_ALL_FIELDS_STRING + " FROM access_token " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *AccessTokenDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM access_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *AccessTokenDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM access_token " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *AccessTokenDao) GetQuery() *AccessTokenQuery {
	return NewAccessTokenQuery(dao)
}

const ACCOUNT_OPERATION_TABLE_NAME = "account_operation"

type ACCOUNT_OPERATION_FIELD string

const ACCOUNT_OPERATION_FIELD_ID = ACCOUNT_OPERATION_FIELD("id")
const ACCOUNT_OPERATION_FIELD_USER_ID = ACCOUNT_OPERATION_FIELD("user_id")
const ACCOUNT_OPERATION_FIELD_OPERATIONTYPE = ACCOUNT_OPERATION_FIELD("operationType")
const ACCOUNT_OPERATION_FIELD_USER_AGENT = ACCOUNT_OPERATION_FIELD("user_agent")
const ACCOUNT_OPERATION_FIELD_PHONE_NUMBER = ACCOUNT_OPERATION_FIELD("phone_number")
const ACCOUNT_OPERATION_FIELD_CREATE_TIME = ACCOUNT_OPERATION_FIELD("create_time")

const ACCOUNT_OPERATION_ALL_FIELDS_STRING = "id,user_id,operationType,user_agent,phone_number,create_time"

var ACCOUNT_OPERATION_ALL_FIELDS = []string{
	"id",
	"user_id",
	"operationType",
	"user_agent",
	"phone_number",
	"create_time",
}

type AccountOperation struct {
	Id            uint64 //size=20
	UserId        string //size=32
	OperationType string //size=32
	UserAgent     string //size=256
	PhoneNumber   string //size=32
	CreateTime    time.Time
}

type AccountOperationQuery struct {
	BaseQuery
	dao *AccountOperationDao
}

func NewAccountOperationQuery(dao *AccountOperationDao) *AccountOperationQuery {
	q := &AccountOperationQuery{}
	q.dao = dao

	return q
}

func (q *AccountOperationQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*AccountOperation, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *AccountOperationQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*AccountOperation, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *AccountOperationQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *AccountOperationQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *AccountOperationQuery) ForUpdate() *AccountOperationQuery {
	q.forUpdate = true
	return q
}

func (q *AccountOperationQuery) ForShare() *AccountOperationQuery {
	q.forShare = true
	return q
}

func (q *AccountOperationQuery) GroupBy(fields ...ACCOUNT_OPERATION_FIELD) *AccountOperationQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *AccountOperationQuery) Limit(startIncluded int64, count int64) *AccountOperationQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *AccountOperationQuery) OrderBy(fieldName ACCOUNT_OPERATION_FIELD, asc bool) *AccountOperationQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *AccountOperationQuery) OrderByGroupCount(asc bool) *AccountOperationQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *AccountOperationQuery) w(format string, a ...interface{}) *AccountOperationQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *AccountOperationQuery) Left() *AccountOperationQuery  { return q.w(" ( ") }
func (q *AccountOperationQuery) Right() *AccountOperationQuery { return q.w(" ) ") }
func (q *AccountOperationQuery) And() *AccountOperationQuery   { return q.w(" AND ") }
func (q *AccountOperationQuery) Or() *AccountOperationQuery    { return q.w(" OR ") }
func (q *AccountOperationQuery) Not() *AccountOperationQuery   { return q.w(" NOT ") }

func (q *AccountOperationQuery) Id_Equal(v uint64) *AccountOperationQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_NotEqual(v uint64) *AccountOperationQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_Less(v uint64) *AccountOperationQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_LessEqual(v uint64) *AccountOperationQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_Greater(v uint64) *AccountOperationQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_GreaterEqual(v uint64) *AccountOperationQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserId_Equal(v string) *AccountOperationQuery {
	return q.w("user_id='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserId_NotEqual(v string) *AccountOperationQuery {
	return q.w("user_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserId_Less(v string) *AccountOperationQuery {
	return q.w("user_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserId_LessEqual(v string) *AccountOperationQuery {
	return q.w("user_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserId_Greater(v string) *AccountOperationQuery {
	return q.w("user_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserId_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("user_id>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_Equal(v string) *AccountOperationQuery {
	return q.w("operationType='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_NotEqual(v string) *AccountOperationQuery {
	return q.w("operationType<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_Less(v string) *AccountOperationQuery {
	return q.w("operationType<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_LessEqual(v string) *AccountOperationQuery {
	return q.w("operationType<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_Greater(v string) *AccountOperationQuery {
	return q.w("operationType>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("operationType>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_Equal(v string) *AccountOperationQuery {
	return q.w("user_agent='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_NotEqual(v string) *AccountOperationQuery {
	return q.w("user_agent<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_Less(v string) *AccountOperationQuery {
	return q.w("user_agent<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_LessEqual(v string) *AccountOperationQuery {
	return q.w("user_agent<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_Greater(v string) *AccountOperationQuery {
	return q.w("user_agent>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("user_agent>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneNumber_Equal(v string) *AccountOperationQuery {
	return q.w("phone_number='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneNumber_NotEqual(v string) *AccountOperationQuery {
	return q.w("phone_number<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneNumber_Less(v string) *AccountOperationQuery {
	return q.w("phone_number<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneNumber_LessEqual(v string) *AccountOperationQuery {
	return q.w("phone_number<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneNumber_Greater(v string) *AccountOperationQuery {
	return q.w("phone_number>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneNumber_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("phone_number>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) CreateTime_Equal(v time.Time) *AccountOperationQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) CreateTime_NotEqual(v time.Time) *AccountOperationQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) CreateTime_Less(v time.Time) *AccountOperationQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) CreateTime_LessEqual(v time.Time) *AccountOperationQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) CreateTime_Greater(v time.Time) *AccountOperationQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) CreateTime_GreaterEqual(v time.Time) *AccountOperationQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}

type AccountOperationDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewAccountOperationDao(db *DB) (t *AccountOperationDao, err error) {
	t = &AccountOperationDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *AccountOperationDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountOperationDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO account_operation (user_id,operationType,user_agent,phone_number) VALUES (?,?,?,?)")
	return err
}

func (dao *AccountOperationDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE account_operation SET user_id=?,operationType=?,user_agent=?,phone_number=? WHERE id=?")
	return err
}

func (dao *AccountOperationDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM account_operation WHERE id=?")
	return err
}

func (dao *AccountOperationDao) Insert(ctx context.Context, tx *wrap.Tx, e *AccountOperation) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.UserId, e.OperationType, e.UserAgent, e.PhoneNumber)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *AccountOperationDao) Update(ctx context.Context, tx *wrap.Tx, e *AccountOperation) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.UserId, e.OperationType, e.UserAgent, e.PhoneNumber, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountOperationDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountOperationDao) scanRow(row *wrap.Row) (*AccountOperation, error) {
	e := &AccountOperation{}
	err := row.Scan(&e.Id, &e.UserId, &e.OperationType, &e.UserAgent, &e.PhoneNumber, &e.CreateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *AccountOperationDao) scanRows(rows *wrap.Rows) (list []*AccountOperation, err error) {
	list = make([]*AccountOperation, 0)
	for rows.Next() {
		e := AccountOperation{}
		err = rows.Scan(&e.Id, &e.UserId, &e.OperationType, &e.UserAgent, &e.PhoneNumber, &e.CreateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *AccountOperationDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*AccountOperation, error) {
	querySql := "SELECT " + ACCOUNT_OPERATION_ALL_FIELDS_STRING + " FROM account_operation " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *AccountOperationDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*AccountOperation, err error) {
	querySql := "SELECT " + ACCOUNT_OPERATION_ALL_FIELDS_STRING + " FROM account_operation " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *AccountOperationDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM account_operation " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *AccountOperationDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM account_operation " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *AccountOperationDao) GetQuery() *AccountOperationQuery {
	return NewAccountOperationQuery(dao)
}

const LOGIN_SMS_CODE_TABLE_NAME = "login_sms_code"

type LOGIN_SMS_CODE_FIELD string

const LOGIN_SMS_CODE_FIELD_ID = LOGIN_SMS_CODE_FIELD("id")
const LOGIN_SMS_CODE_FIELD_PHONE_NUMBER = LOGIN_SMS_CODE_FIELD("phone_number")
const LOGIN_SMS_CODE_FIELD_SMS_CODE = LOGIN_SMS_CODE_FIELD("sms_code")
const LOGIN_SMS_CODE_FIELD_CREATE_TIME = LOGIN_SMS_CODE_FIELD("create_time")
const LOGIN_SMS_CODE_FIELD_UPDATE_TIME = LOGIN_SMS_CODE_FIELD("update_time")

const LOGIN_SMS_CODE_ALL_FIELDS_STRING = "id,phone_number,sms_code,create_time,update_time"

var LOGIN_SMS_CODE_ALL_FIELDS = []string{
	"id",
	"phone_number",
	"sms_code",
	"create_time",
	"update_time",
}

type LoginSmsCode struct {
	Id          uint64 //size=20
	PhoneNumber string //size=32
	SmsCode     string //size=8
	CreateTime  time.Time
	UpdateTime  time.Time
}

type LoginSmsCodeQuery struct {
	BaseQuery
	dao *LoginSmsCodeDao
}

func NewLoginSmsCodeQuery(dao *LoginSmsCodeDao) *LoginSmsCodeQuery {
	q := &LoginSmsCodeQuery{}
	q.dao = dao

	return q
}

func (q *LoginSmsCodeQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*LoginSmsCode, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *LoginSmsCodeQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*LoginSmsCode, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *LoginSmsCodeQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *LoginSmsCodeQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *LoginSmsCodeQuery) ForUpdate() *LoginSmsCodeQuery {
	q.forUpdate = true
	return q
}

func (q *LoginSmsCodeQuery) ForShare() *LoginSmsCodeQuery {
	q.forShare = true
	return q
}

func (q *LoginSmsCodeQuery) GroupBy(fields ...LOGIN_SMS_CODE_FIELD) *LoginSmsCodeQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *LoginSmsCodeQuery) Limit(startIncluded int64, count int64) *LoginSmsCodeQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *LoginSmsCodeQuery) OrderBy(fieldName LOGIN_SMS_CODE_FIELD, asc bool) *LoginSmsCodeQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *LoginSmsCodeQuery) OrderByGroupCount(asc bool) *LoginSmsCodeQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *LoginSmsCodeQuery) w(format string, a ...interface{}) *LoginSmsCodeQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *LoginSmsCodeQuery) Left() *LoginSmsCodeQuery  { return q.w(" ( ") }
func (q *LoginSmsCodeQuery) Right() *LoginSmsCodeQuery { return q.w(" ) ") }
func (q *LoginSmsCodeQuery) And() *LoginSmsCodeQuery   { return q.w(" AND ") }
func (q *LoginSmsCodeQuery) Or() *LoginSmsCodeQuery    { return q.w(" OR ") }
func (q *LoginSmsCodeQuery) Not() *LoginSmsCodeQuery   { return q.w(" NOT ") }

func (q *LoginSmsCodeQuery) Id_Equal(v uint64) *LoginSmsCodeQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) Id_NotEqual(v uint64) *LoginSmsCodeQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) Id_Less(v uint64) *LoginSmsCodeQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) Id_LessEqual(v uint64) *LoginSmsCodeQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) Id_Greater(v uint64) *LoginSmsCodeQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) Id_GreaterEqual(v uint64) *LoginSmsCodeQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) PhoneNumber_Equal(v string) *LoginSmsCodeQuery {
	return q.w("phone_number='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) PhoneNumber_NotEqual(v string) *LoginSmsCodeQuery {
	return q.w("phone_number<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) PhoneNumber_Less(v string) *LoginSmsCodeQuery {
	return q.w("phone_number<'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) PhoneNumber_LessEqual(v string) *LoginSmsCodeQuery {
	return q.w("phone_number<='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) PhoneNumber_Greater(v string) *LoginSmsCodeQuery {
	return q.w("phone_number>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) PhoneNumber_GreaterEqual(v string) *LoginSmsCodeQuery {
	return q.w("phone_number>='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) SmsCode_Equal(v string) *LoginSmsCodeQuery {
	return q.w("sms_code='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) SmsCode_NotEqual(v string) *LoginSmsCodeQuery {
	return q.w("sms_code<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) SmsCode_Less(v string) *LoginSmsCodeQuery {
	return q.w("sms_code<'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) SmsCode_LessEqual(v string) *LoginSmsCodeQuery {
	return q.w("sms_code<='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) SmsCode_Greater(v string) *LoginSmsCodeQuery {
	return q.w("sms_code>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) SmsCode_GreaterEqual(v string) *LoginSmsCodeQuery {
	return q.w("sms_code>='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) CreateTime_Equal(v time.Time) *LoginSmsCodeQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) CreateTime_NotEqual(v time.Time) *LoginSmsCodeQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) CreateTime_Less(v time.Time) *LoginSmsCodeQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) CreateTime_LessEqual(v time.Time) *LoginSmsCodeQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) CreateTime_Greater(v time.Time) *LoginSmsCodeQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) CreateTime_GreaterEqual(v time.Time) *LoginSmsCodeQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) UpdateTime_Equal(v time.Time) *LoginSmsCodeQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) UpdateTime_NotEqual(v time.Time) *LoginSmsCodeQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) UpdateTime_Less(v time.Time) *LoginSmsCodeQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) UpdateTime_LessEqual(v time.Time) *LoginSmsCodeQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) UpdateTime_Greater(v time.Time) *LoginSmsCodeQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *LoginSmsCodeQuery) UpdateTime_GreaterEqual(v time.Time) *LoginSmsCodeQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type LoginSmsCodeDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewLoginSmsCodeDao(db *DB) (t *LoginSmsCodeDao, err error) {
	t = &LoginSmsCodeDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *LoginSmsCodeDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *LoginSmsCodeDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO login_sms_code (phone_number,sms_code) VALUES (?,?)")
	return err
}

func (dao *LoginSmsCodeDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE login_sms_code SET phone_number=?,sms_code=? WHERE id=?")
	return err
}

func (dao *LoginSmsCodeDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM login_sms_code WHERE id=?")
	return err
}

func (dao *LoginSmsCodeDao) Insert(ctx context.Context, tx *wrap.Tx, e *LoginSmsCode) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.PhoneNumber, e.SmsCode)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *LoginSmsCodeDao) Update(ctx context.Context, tx *wrap.Tx, e *LoginSmsCode) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.PhoneNumber, e.SmsCode, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *LoginSmsCodeDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *LoginSmsCodeDao) scanRow(row *wrap.Row) (*LoginSmsCode, error) {
	e := &LoginSmsCode{}
	err := row.Scan(&e.Id, &e.PhoneNumber, &e.SmsCode, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *LoginSmsCodeDao) scanRows(rows *wrap.Rows) (list []*LoginSmsCode, err error) {
	list = make([]*LoginSmsCode, 0)
	for rows.Next() {
		e := LoginSmsCode{}
		err = rows.Scan(&e.Id, &e.PhoneNumber, &e.SmsCode, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *LoginSmsCodeDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*LoginSmsCode, error) {
	querySql := "SELECT " + LOGIN_SMS_CODE_ALL_FIELDS_STRING + " FROM login_sms_code " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *LoginSmsCodeDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*LoginSmsCode, err error) {
	querySql := "SELECT " + LOGIN_SMS_CODE_ALL_FIELDS_STRING + " FROM login_sms_code " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *LoginSmsCodeDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM login_sms_code " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *LoginSmsCodeDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM login_sms_code " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *LoginSmsCodeDao) GetQuery() *LoginSmsCodeQuery {
	return NewLoginSmsCodeQuery(dao)
}

const OAUTH_ACCOUNT_TABLE_NAME = "oauth_account"

type OAUTH_ACCOUNT_FIELD string

const OAUTH_ACCOUNT_FIELD_ID = OAUTH_ACCOUNT_FIELD("id")
const OAUTH_ACCOUNT_FIELD_USER_ID = OAUTH_ACCOUNT_FIELD("user_id")
const OAUTH_ACCOUNT_FIELD_OAUTH_PROVIDER = OAUTH_ACCOUNT_FIELD("oauth_provider")
const OAUTH_ACCOUNT_FIELD_OAUTH_OPEN_ID = OAUTH_ACCOUNT_FIELD("oauth_open_id")
const OAUTH_ACCOUNT_FIELD_OAUTH_NAME = OAUTH_ACCOUNT_FIELD("oauth_name")
const OAUTH_ACCOUNT_FIELD_OAUTH_ICON = OAUTH_ACCOUNT_FIELD("oauth_icon")
const OAUTH_ACCOUNT_FIELD_CREATE_TIME = OAUTH_ACCOUNT_FIELD("create_time")
const OAUTH_ACCOUNT_FIELD_UPDATE_TIME = OAUTH_ACCOUNT_FIELD("update_time")

const OAUTH_ACCOUNT_ALL_FIELDS_STRING = "id,user_id,oauth_provider,oauth_open_id,oauth_name,oauth_icon,create_time,update_time"

var OAUTH_ACCOUNT_ALL_FIELDS = []string{
	"id",
	"user_id",
	"oauth_provider",
	"oauth_open_id",
	"oauth_name",
	"oauth_icon",
	"create_time",
	"update_time",
}

type OauthAccount struct {
	Id            uint64 //size=20
	UserId        string //size=32
	OauthProvider string //size=32
	OauthOpenId   string //size=128
	OauthName     string //size=32
	OauthIcon     string //size=256
	CreateTime    time.Time
	UpdateTime    mysql.NullTime
}

type OauthAccountQuery struct {
	BaseQuery
	dao *OauthAccountDao
}

func NewOauthAccountQuery(dao *OauthAccountDao) *OauthAccountQuery {
	q := &OauthAccountQuery{}
	q.dao = dao

	return q
}

func (q *OauthAccountQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*OauthAccount, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *OauthAccountQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*OauthAccount, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *OauthAccountQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *OauthAccountQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *OauthAccountQuery) ForUpdate() *OauthAccountQuery {
	q.forUpdate = true
	return q
}

func (q *OauthAccountQuery) ForShare() *OauthAccountQuery {
	q.forShare = true
	return q
}

func (q *OauthAccountQuery) GroupBy(fields ...OAUTH_ACCOUNT_FIELD) *OauthAccountQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *OauthAccountQuery) Limit(startIncluded int64, count int64) *OauthAccountQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *OauthAccountQuery) OrderBy(fieldName OAUTH_ACCOUNT_FIELD, asc bool) *OauthAccountQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthAccountQuery) OrderByGroupCount(asc bool) *OauthAccountQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthAccountQuery) w(format string, a ...interface{}) *OauthAccountQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *OauthAccountQuery) Left() *OauthAccountQuery  { return q.w(" ( ") }
func (q *OauthAccountQuery) Right() *OauthAccountQuery { return q.w(" ) ") }
func (q *OauthAccountQuery) And() *OauthAccountQuery   { return q.w(" AND ") }
func (q *OauthAccountQuery) Or() *OauthAccountQuery    { return q.w(" OR ") }
func (q *OauthAccountQuery) Not() *OauthAccountQuery   { return q.w(" NOT ") }

func (q *OauthAccountQuery) Id_Equal(v uint64) *OauthAccountQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) Id_NotEqual(v uint64) *OauthAccountQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) Id_Less(v uint64) *OauthAccountQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) Id_LessEqual(v uint64) *OauthAccountQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) Id_Greater(v uint64) *OauthAccountQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) Id_GreaterEqual(v uint64) *OauthAccountQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserId_Equal(v string) *OauthAccountQuery {
	return q.w("user_id='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserId_NotEqual(v string) *OauthAccountQuery {
	return q.w("user_id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserId_Less(v string) *OauthAccountQuery {
	return q.w("user_id<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserId_LessEqual(v string) *OauthAccountQuery {
	return q.w("user_id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserId_Greater(v string) *OauthAccountQuery {
	return q.w("user_id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserId_GreaterEqual(v string) *OauthAccountQuery {
	return q.w("user_id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthProvider_Equal(v string) *OauthAccountQuery {
	return q.w("oauth_provider='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthProvider_NotEqual(v string) *OauthAccountQuery {
	return q.w("oauth_provider<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthProvider_Less(v string) *OauthAccountQuery {
	return q.w("oauth_provider<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthProvider_LessEqual(v string) *OauthAccountQuery {
	return q.w("oauth_provider<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthProvider_Greater(v string) *OauthAccountQuery {
	return q.w("oauth_provider>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthProvider_GreaterEqual(v string) *OauthAccountQuery {
	return q.w("oauth_provider>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthOpenId_Equal(v string) *OauthAccountQuery {
	return q.w("oauth_open_id='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthOpenId_NotEqual(v string) *OauthAccountQuery {
	return q.w("oauth_open_id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthOpenId_Less(v string) *OauthAccountQuery {
	return q.w("oauth_open_id<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthOpenId_LessEqual(v string) *OauthAccountQuery {
	return q.w("oauth_open_id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthOpenId_Greater(v string) *OauthAccountQuery {
	return q.w("oauth_open_id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthOpenId_GreaterEqual(v string) *OauthAccountQuery {
	return q.w("oauth_open_id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthName_Equal(v string) *OauthAccountQuery {
	return q.w("oauth_name='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthName_NotEqual(v string) *OauthAccountQuery {
	return q.w("oauth_name<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthName_Less(v string) *OauthAccountQuery {
	return q.w("oauth_name<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthName_LessEqual(v string) *OauthAccountQuery {
	return q.w("oauth_name<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthName_Greater(v string) *OauthAccountQuery {
	return q.w("oauth_name>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthName_GreaterEqual(v string) *OauthAccountQuery {
	return q.w("oauth_name>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthIcon_Equal(v string) *OauthAccountQuery {
	return q.w("oauth_icon='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthIcon_NotEqual(v string) *OauthAccountQuery {
	return q.w("oauth_icon<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthIcon_Less(v string) *OauthAccountQuery {
	return q.w("oauth_icon<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthIcon_LessEqual(v string) *OauthAccountQuery {
	return q.w("oauth_icon<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthIcon_Greater(v string) *OauthAccountQuery {
	return q.w("oauth_icon>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OauthIcon_GreaterEqual(v string) *OauthAccountQuery {
	return q.w("oauth_icon>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) CreateTime_Equal(v time.Time) *OauthAccountQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) CreateTime_NotEqual(v time.Time) *OauthAccountQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) CreateTime_Less(v time.Time) *OauthAccountQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) CreateTime_LessEqual(v time.Time) *OauthAccountQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) CreateTime_Greater(v time.Time) *OauthAccountQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) CreateTime_GreaterEqual(v time.Time) *OauthAccountQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UpdateTime_Equal(v time.Time) *OauthAccountQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UpdateTime_NotEqual(v time.Time) *OauthAccountQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UpdateTime_Less(v time.Time) *OauthAccountQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UpdateTime_LessEqual(v time.Time) *OauthAccountQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UpdateTime_Greater(v time.Time) *OauthAccountQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UpdateTime_GreaterEqual(v time.Time) *OauthAccountQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UpdateTime_IsNull() *OauthAccountQuery { return q.w("update_time IS NULL") }
func (q *OauthAccountQuery) UpdateTime_NotNull() *OauthAccountQuery {
	return q.w("update_time IS NOT NULL")
}

type OauthAccountDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewOauthAccountDao(db *DB) (t *OauthAccountDao, err error) {
	t = &OauthAccountDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *OauthAccountDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthAccountDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO oauth_account (user_id,oauth_provider,oauth_open_id,oauth_name,oauth_icon) VALUES (?,?,?,?,?)")
	return err
}

func (dao *OauthAccountDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE oauth_account SET user_id=?,oauth_provider=?,oauth_open_id=?,oauth_name=?,oauth_icon=? WHERE id=?")
	return err
}

func (dao *OauthAccountDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM oauth_account WHERE id=?")
	return err
}

func (dao *OauthAccountDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthAccount) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.UserId, e.OauthProvider, e.OauthOpenId, e.OauthName, e.OauthIcon)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *OauthAccountDao) Update(ctx context.Context, tx *wrap.Tx, e *OauthAccount) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.UserId, e.OauthProvider, e.OauthOpenId, e.OauthName, e.OauthIcon, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthAccountDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthAccountDao) scanRow(row *wrap.Row) (*OauthAccount, error) {
	e := &OauthAccount{}
	err := row.Scan(&e.Id, &e.UserId, &e.OauthProvider, &e.OauthOpenId, &e.OauthName, &e.OauthIcon, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *OauthAccountDao) scanRows(rows *wrap.Rows) (list []*OauthAccount, err error) {
	list = make([]*OauthAccount, 0)
	for rows.Next() {
		e := OauthAccount{}
		err = rows.Scan(&e.Id, &e.UserId, &e.OauthProvider, &e.OauthOpenId, &e.OauthName, &e.OauthIcon, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *OauthAccountDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*OauthAccount, error) {
	querySql := "SELECT " + OAUTH_ACCOUNT_ALL_FIELDS_STRING + " FROM oauth_account " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *OauthAccountDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*OauthAccount, err error) {
	querySql := "SELECT " + OAUTH_ACCOUNT_ALL_FIELDS_STRING + " FROM oauth_account " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *OauthAccountDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM oauth_account " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *OauthAccountDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM oauth_account " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *OauthAccountDao) GetQuery() *OauthAccountQuery {
	return NewOauthAccountQuery(dao)
}

const OAUTH_STATE_TABLE_NAME = "oauth_state"

type OAUTH_STATE_FIELD string

const OAUTH_STATE_FIELD_ID = OAUTH_STATE_FIELD("id")
const OAUTH_STATE_FIELD_OAUTH_STATE = OAUTH_STATE_FIELD("oauth_state")
const OAUTH_STATE_FIELD_IS_USED = OAUTH_STATE_FIELD("is_used")
const OAUTH_STATE_FIELD_USER_AGENT = OAUTH_STATE_FIELD("user_agent")
const OAUTH_STATE_FIELD_CREATE_TIME = OAUTH_STATE_FIELD("create_time")
const OAUTH_STATE_FIELD_UPDATE_TIME = OAUTH_STATE_FIELD("update_time")

const OAUTH_STATE_ALL_FIELDS_STRING = "id,oauth_state,is_used,user_agent,create_time,update_time"

var OAUTH_STATE_ALL_FIELDS = []string{
	"id",
	"oauth_state",
	"is_used",
	"user_agent",
	"create_time",
	"update_time",
}

type OauthState struct {
	Id         uint64 //size=20
	OauthState string //size=128
	IsUsed     int32  //size=1
	UserAgent  string //size=256
	CreateTime time.Time
	UpdateTime time.Time
}

type OauthStateQuery struct {
	BaseQuery
	dao *OauthStateDao
}

func NewOauthStateQuery(dao *OauthStateDao) *OauthStateQuery {
	q := &OauthStateQuery{}
	q.dao = dao

	return q
}

func (q *OauthStateQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*OauthState, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *OauthStateQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*OauthState, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *OauthStateQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *OauthStateQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *OauthStateQuery) ForUpdate() *OauthStateQuery {
	q.forUpdate = true
	return q
}

func (q *OauthStateQuery) ForShare() *OauthStateQuery {
	q.forShare = true
	return q
}

func (q *OauthStateQuery) GroupBy(fields ...OAUTH_STATE_FIELD) *OauthStateQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *OauthStateQuery) Limit(startIncluded int64, count int64) *OauthStateQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *OauthStateQuery) OrderBy(fieldName OAUTH_STATE_FIELD, asc bool) *OauthStateQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthStateQuery) OrderByGroupCount(asc bool) *OauthStateQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthStateQuery) w(format string, a ...interface{}) *OauthStateQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *OauthStateQuery) Left() *OauthStateQuery  { return q.w(" ( ") }
func (q *OauthStateQuery) Right() *OauthStateQuery { return q.w(" ) ") }
func (q *OauthStateQuery) And() *OauthStateQuery   { return q.w(" AND ") }
func (q *OauthStateQuery) Or() *OauthStateQuery    { return q.w(" OR ") }
func (q *OauthStateQuery) Not() *OauthStateQuery   { return q.w(" NOT ") }

func (q *OauthStateQuery) Id_Equal(v uint64) *OauthStateQuery { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *OauthStateQuery) Id_NotEqual(v uint64) *OauthStateQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) Id_Less(v uint64) *OauthStateQuery { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *OauthStateQuery) Id_LessEqual(v uint64) *OauthStateQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) Id_Greater(v uint64) *OauthStateQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) Id_GreaterEqual(v uint64) *OauthStateQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_Equal(v string) *OauthStateQuery {
	return q.w("oauth_state='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_NotEqual(v string) *OauthStateQuery {
	return q.w("oauth_state<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_Less(v string) *OauthStateQuery {
	return q.w("oauth_state<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_LessEqual(v string) *OauthStateQuery {
	return q.w("oauth_state<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_Greater(v string) *OauthStateQuery {
	return q.w("oauth_state>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_GreaterEqual(v string) *OauthStateQuery {
	return q.w("oauth_state>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) IsUsed_Equal(v int32) *OauthStateQuery {
	return q.w("is_used='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) IsUsed_NotEqual(v int32) *OauthStateQuery {
	return q.w("is_used<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) IsUsed_Less(v int32) *OauthStateQuery {
	return q.w("is_used<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) IsUsed_LessEqual(v int32) *OauthStateQuery {
	return q.w("is_used<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) IsUsed_Greater(v int32) *OauthStateQuery {
	return q.w("is_used>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) IsUsed_GreaterEqual(v int32) *OauthStateQuery {
	return q.w("is_used>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UserAgent_Equal(v string) *OauthStateQuery {
	return q.w("user_agent='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UserAgent_NotEqual(v string) *OauthStateQuery {
	return q.w("user_agent<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UserAgent_Less(v string) *OauthStateQuery {
	return q.w("user_agent<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UserAgent_LessEqual(v string) *OauthStateQuery {
	return q.w("user_agent<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UserAgent_Greater(v string) *OauthStateQuery {
	return q.w("user_agent>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UserAgent_GreaterEqual(v string) *OauthStateQuery {
	return q.w("user_agent>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_Equal(v time.Time) *OauthStateQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_NotEqual(v time.Time) *OauthStateQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_Less(v time.Time) *OauthStateQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_LessEqual(v time.Time) *OauthStateQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_Greater(v time.Time) *OauthStateQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_GreaterEqual(v time.Time) *OauthStateQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_Equal(v time.Time) *OauthStateQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_NotEqual(v time.Time) *OauthStateQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_Less(v time.Time) *OauthStateQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_LessEqual(v time.Time) *OauthStateQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_Greater(v time.Time) *OauthStateQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_GreaterEqual(v time.Time) *OauthStateQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type OauthStateDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewOauthStateDao(db *DB) (t *OauthStateDao, err error) {
	t = &OauthStateDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *OauthStateDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthStateDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO oauth_state (oauth_state,is_used,user_agent) VALUES (?,?,?)")
	return err
}

func (dao *OauthStateDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE oauth_state SET oauth_state=?,is_used=?,user_agent=? WHERE id=?")
	return err
}

func (dao *OauthStateDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM oauth_state WHERE id=?")
	return err
}

func (dao *OauthStateDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthState) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.OauthState, e.IsUsed, e.UserAgent)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *OauthStateDao) Update(ctx context.Context, tx *wrap.Tx, e *OauthState) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.OauthState, e.IsUsed, e.UserAgent, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthStateDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthStateDao) scanRow(row *wrap.Row) (*OauthState, error) {
	e := &OauthState{}
	err := row.Scan(&e.Id, &e.OauthState, &e.IsUsed, &e.UserAgent, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *OauthStateDao) scanRows(rows *wrap.Rows) (list []*OauthState, err error) {
	list = make([]*OauthState, 0)
	for rows.Next() {
		e := OauthState{}
		err = rows.Scan(&e.Id, &e.OauthState, &e.IsUsed, &e.UserAgent, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *OauthStateDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*OauthState, error) {
	querySql := "SELECT " + OAUTH_STATE_ALL_FIELDS_STRING + " FROM oauth_state " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *OauthStateDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*OauthState, err error) {
	querySql := "SELECT " + OAUTH_STATE_ALL_FIELDS_STRING + " FROM oauth_state " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *OauthStateDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM oauth_state " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *OauthStateDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM oauth_state " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *OauthStateDao) GetQuery() *OauthStateQuery {
	return NewOauthStateQuery(dao)
}

const PHONE_ACCOUNT_TABLE_NAME = "phone_account"

type PHONE_ACCOUNT_FIELD string

const PHONE_ACCOUNT_FIELD_ID = PHONE_ACCOUNT_FIELD("id")
const PHONE_ACCOUNT_FIELD_USER_ID = PHONE_ACCOUNT_FIELD("user_id")
const PHONE_ACCOUNT_FIELD_PHONE_NUMBER = PHONE_ACCOUNT_FIELD("phone_number")
const PHONE_ACCOUNT_FIELD_CREATE_TIME = PHONE_ACCOUNT_FIELD("create_time")
const PHONE_ACCOUNT_FIELD_UPDATE_TIME = PHONE_ACCOUNT_FIELD("update_time")

const PHONE_ACCOUNT_ALL_FIELDS_STRING = "id,user_id,phone_number,create_time,update_time"

var PHONE_ACCOUNT_ALL_FIELDS = []string{
	"id",
	"user_id",
	"phone_number",
	"create_time",
	"update_time",
}

type PhoneAccount struct {
	Id          uint64 //size=20
	UserId      string //size=32
	PhoneNumber string //size=32
	CreateTime  time.Time
	UpdateTime  time.Time
}

type PhoneAccountQuery struct {
	BaseQuery
	dao *PhoneAccountDao
}

func NewPhoneAccountQuery(dao *PhoneAccountDao) *PhoneAccountQuery {
	q := &PhoneAccountQuery{}
	q.dao = dao

	return q
}

func (q *PhoneAccountQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*PhoneAccount, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *PhoneAccountQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*PhoneAccount, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *PhoneAccountQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *PhoneAccountQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *PhoneAccountQuery) ForUpdate() *PhoneAccountQuery {
	q.forUpdate = true
	return q
}

func (q *PhoneAccountQuery) ForShare() *PhoneAccountQuery {
	q.forShare = true
	return q
}

func (q *PhoneAccountQuery) GroupBy(fields ...PHONE_ACCOUNT_FIELD) *PhoneAccountQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *PhoneAccountQuery) Limit(startIncluded int64, count int64) *PhoneAccountQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *PhoneAccountQuery) OrderBy(fieldName PHONE_ACCOUNT_FIELD, asc bool) *PhoneAccountQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *PhoneAccountQuery) OrderByGroupCount(asc bool) *PhoneAccountQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *PhoneAccountQuery) w(format string, a ...interface{}) *PhoneAccountQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *PhoneAccountQuery) Left() *PhoneAccountQuery  { return q.w(" ( ") }
func (q *PhoneAccountQuery) Right() *PhoneAccountQuery { return q.w(" ) ") }
func (q *PhoneAccountQuery) And() *PhoneAccountQuery   { return q.w(" AND ") }
func (q *PhoneAccountQuery) Or() *PhoneAccountQuery    { return q.w(" OR ") }
func (q *PhoneAccountQuery) Not() *PhoneAccountQuery   { return q.w(" NOT ") }

func (q *PhoneAccountQuery) Id_Equal(v uint64) *PhoneAccountQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) Id_NotEqual(v uint64) *PhoneAccountQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) Id_Less(v uint64) *PhoneAccountQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) Id_LessEqual(v uint64) *PhoneAccountQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) Id_Greater(v uint64) *PhoneAccountQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) Id_GreaterEqual(v uint64) *PhoneAccountQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UserId_Equal(v string) *PhoneAccountQuery {
	return q.w("user_id='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UserId_NotEqual(v string) *PhoneAccountQuery {
	return q.w("user_id<>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UserId_Less(v string) *PhoneAccountQuery {
	return q.w("user_id<'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UserId_LessEqual(v string) *PhoneAccountQuery {
	return q.w("user_id<='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UserId_Greater(v string) *PhoneAccountQuery {
	return q.w("user_id>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UserId_GreaterEqual(v string) *PhoneAccountQuery {
	return q.w("user_id>='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) PhoneNumber_Equal(v string) *PhoneAccountQuery {
	return q.w("phone_number='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) PhoneNumber_NotEqual(v string) *PhoneAccountQuery {
	return q.w("phone_number<>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) PhoneNumber_Less(v string) *PhoneAccountQuery {
	return q.w("phone_number<'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) PhoneNumber_LessEqual(v string) *PhoneAccountQuery {
	return q.w("phone_number<='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) PhoneNumber_Greater(v string) *PhoneAccountQuery {
	return q.w("phone_number>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) PhoneNumber_GreaterEqual(v string) *PhoneAccountQuery {
	return q.w("phone_number>='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) CreateTime_Equal(v time.Time) *PhoneAccountQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) CreateTime_NotEqual(v time.Time) *PhoneAccountQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) CreateTime_Less(v time.Time) *PhoneAccountQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) CreateTime_LessEqual(v time.Time) *PhoneAccountQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) CreateTime_Greater(v time.Time) *PhoneAccountQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) CreateTime_GreaterEqual(v time.Time) *PhoneAccountQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UpdateTime_Equal(v time.Time) *PhoneAccountQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UpdateTime_NotEqual(v time.Time) *PhoneAccountQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UpdateTime_Less(v time.Time) *PhoneAccountQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UpdateTime_LessEqual(v time.Time) *PhoneAccountQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UpdateTime_Greater(v time.Time) *PhoneAccountQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) UpdateTime_GreaterEqual(v time.Time) *PhoneAccountQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type PhoneAccountDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewPhoneAccountDao(db *DB) (t *PhoneAccountDao, err error) {
	t = &PhoneAccountDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *PhoneAccountDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *PhoneAccountDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO phone_account (user_id,phone_number) VALUES (?,?)")
	return err
}

func (dao *PhoneAccountDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE phone_account SET user_id=?,phone_number=? WHERE id=?")
	return err
}

func (dao *PhoneAccountDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM phone_account WHERE id=?")
	return err
}

func (dao *PhoneAccountDao) Insert(ctx context.Context, tx *wrap.Tx, e *PhoneAccount) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.UserId, e.PhoneNumber)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *PhoneAccountDao) Update(ctx context.Context, tx *wrap.Tx, e *PhoneAccount) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.UserId, e.PhoneNumber, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *PhoneAccountDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *PhoneAccountDao) scanRow(row *wrap.Row) (*PhoneAccount, error) {
	e := &PhoneAccount{}
	err := row.Scan(&e.Id, &e.UserId, &e.PhoneNumber, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *PhoneAccountDao) scanRows(rows *wrap.Rows) (list []*PhoneAccount, err error) {
	list = make([]*PhoneAccount, 0)
	for rows.Next() {
		e := PhoneAccount{}
		err = rows.Scan(&e.Id, &e.UserId, &e.PhoneNumber, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *PhoneAccountDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*PhoneAccount, error) {
	querySql := "SELECT " + PHONE_ACCOUNT_ALL_FIELDS_STRING + " FROM phone_account " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *PhoneAccountDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*PhoneAccount, err error) {
	querySql := "SELECT " + PHONE_ACCOUNT_ALL_FIELDS_STRING + " FROM phone_account " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *PhoneAccountDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM phone_account " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *PhoneAccountDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM phone_account " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *PhoneAccountDao) GetQuery() *PhoneAccountQuery {
	return NewPhoneAccountQuery(dao)
}

const REFRESH_TOKEN_TABLE_NAME = "refresh_token"

type REFRESH_TOKEN_FIELD string

const REFRESH_TOKEN_FIELD_ID = REFRESH_TOKEN_FIELD("id")
const REFRESH_TOKEN_FIELD_USER_ID = REFRESH_TOKEN_FIELD("user_id")
const REFRESH_TOKEN_FIELD_REFRESH_TOKEN = REFRESH_TOKEN_FIELD("refresh_token")
const REFRESH_TOKEN_FIELD_IS_LOGOUT = REFRESH_TOKEN_FIELD("is_logout")
const REFRESH_TOKEN_FIELD_LOGOUT_TIME = REFRESH_TOKEN_FIELD("logout_time")
const REFRESH_TOKEN_FIELD_CREATE_TIME = REFRESH_TOKEN_FIELD("create_time")
const REFRESH_TOKEN_FIELD_UPDATE_TIME = REFRESH_TOKEN_FIELD("update_time")

const REFRESH_TOKEN_ALL_FIELDS_STRING = "id,user_id,refresh_token,is_logout,logout_time,create_time,update_time"

var REFRESH_TOKEN_ALL_FIELDS = []string{
	"id",
	"user_id",
	"refresh_token",
	"is_logout",
	"logout_time",
	"create_time",
	"update_time",
}

type RefreshToken struct {
	Id           uint64 //size=20
	UserId       string //size=32
	RefreshToken string //size=128
	IsLogout     int32  //size=1
	LogoutTime   time.Time
	CreateTime   time.Time
	UpdateTime   time.Time
}

type RefreshTokenQuery struct {
	BaseQuery
	dao *RefreshTokenDao
}

func NewRefreshTokenQuery(dao *RefreshTokenDao) *RefreshTokenQuery {
	q := &RefreshTokenQuery{}
	q.dao = dao

	return q
}

func (q *RefreshTokenQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*RefreshToken, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *RefreshTokenQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*RefreshToken, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *RefreshTokenQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *RefreshTokenQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *RefreshTokenQuery) ForUpdate() *RefreshTokenQuery {
	q.forUpdate = true
	return q
}

func (q *RefreshTokenQuery) ForShare() *RefreshTokenQuery {
	q.forShare = true
	return q
}

func (q *RefreshTokenQuery) GroupBy(fields ...REFRESH_TOKEN_FIELD) *RefreshTokenQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *RefreshTokenQuery) Limit(startIncluded int64, count int64) *RefreshTokenQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *RefreshTokenQuery) OrderBy(fieldName REFRESH_TOKEN_FIELD, asc bool) *RefreshTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *RefreshTokenQuery) OrderByGroupCount(asc bool) *RefreshTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *RefreshTokenQuery) w(format string, a ...interface{}) *RefreshTokenQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *RefreshTokenQuery) Left() *RefreshTokenQuery  { return q.w(" ( ") }
func (q *RefreshTokenQuery) Right() *RefreshTokenQuery { return q.w(" ) ") }
func (q *RefreshTokenQuery) And() *RefreshTokenQuery   { return q.w(" AND ") }
func (q *RefreshTokenQuery) Or() *RefreshTokenQuery    { return q.w(" OR ") }
func (q *RefreshTokenQuery) Not() *RefreshTokenQuery   { return q.w(" NOT ") }

func (q *RefreshTokenQuery) Id_Equal(v uint64) *RefreshTokenQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_NotEqual(v uint64) *RefreshTokenQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_Less(v uint64) *RefreshTokenQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_LessEqual(v uint64) *RefreshTokenQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_Greater(v uint64) *RefreshTokenQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_GreaterEqual(v uint64) *RefreshTokenQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserId_Equal(v string) *RefreshTokenQuery {
	return q.w("user_id='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserId_NotEqual(v string) *RefreshTokenQuery {
	return q.w("user_id<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserId_Less(v string) *RefreshTokenQuery {
	return q.w("user_id<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserId_LessEqual(v string) *RefreshTokenQuery {
	return q.w("user_id<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserId_Greater(v string) *RefreshTokenQuery {
	return q.w("user_id>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserId_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("user_id>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_Equal(v string) *RefreshTokenQuery {
	return q.w("refresh_token='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_NotEqual(v string) *RefreshTokenQuery {
	return q.w("refresh_token<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_Less(v string) *RefreshTokenQuery {
	return q.w("refresh_token<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_LessEqual(v string) *RefreshTokenQuery {
	return q.w("refresh_token<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_Greater(v string) *RefreshTokenQuery {
	return q.w("refresh_token>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("refresh_token>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_Equal(v int32) *RefreshTokenQuery {
	return q.w("is_logout='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_NotEqual(v int32) *RefreshTokenQuery {
	return q.w("is_logout<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_Less(v int32) *RefreshTokenQuery {
	return q.w("is_logout<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_LessEqual(v int32) *RefreshTokenQuery {
	return q.w("is_logout<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_Greater(v int32) *RefreshTokenQuery {
	return q.w("is_logout>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_GreaterEqual(v int32) *RefreshTokenQuery {
	return q.w("is_logout>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) LogoutTime_Equal(v time.Time) *RefreshTokenQuery {
	return q.w("logout_time='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) LogoutTime_NotEqual(v time.Time) *RefreshTokenQuery {
	return q.w("logout_time<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) LogoutTime_Less(v time.Time) *RefreshTokenQuery {
	return q.w("logout_time<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) LogoutTime_LessEqual(v time.Time) *RefreshTokenQuery {
	return q.w("logout_time<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) LogoutTime_Greater(v time.Time) *RefreshTokenQuery {
	return q.w("logout_time>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) LogoutTime_GreaterEqual(v time.Time) *RefreshTokenQuery {
	return q.w("logout_time>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_Equal(v time.Time) *RefreshTokenQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_NotEqual(v time.Time) *RefreshTokenQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_Less(v time.Time) *RefreshTokenQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_LessEqual(v time.Time) *RefreshTokenQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_Greater(v time.Time) *RefreshTokenQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_GreaterEqual(v time.Time) *RefreshTokenQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_Equal(v time.Time) *RefreshTokenQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_NotEqual(v time.Time) *RefreshTokenQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_Less(v time.Time) *RefreshTokenQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_LessEqual(v time.Time) *RefreshTokenQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_Greater(v time.Time) *RefreshTokenQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_GreaterEqual(v time.Time) *RefreshTokenQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type RefreshTokenDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewRefreshTokenDao(db *DB) (t *RefreshTokenDao, err error) {
	t = &RefreshTokenDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *RefreshTokenDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *RefreshTokenDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO refresh_token (user_id,refresh_token,is_logout,logout_time) VALUES (?,?,?,?)")
	return err
}

func (dao *RefreshTokenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE refresh_token SET user_id=?,refresh_token=?,is_logout=?,logout_time=? WHERE id=?")
	return err
}

func (dao *RefreshTokenDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM refresh_token WHERE id=?")
	return err
}

func (dao *RefreshTokenDao) Insert(ctx context.Context, tx *wrap.Tx, e *RefreshToken) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.UserId, e.RefreshToken, e.IsLogout, e.LogoutTime)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *RefreshTokenDao) Update(ctx context.Context, tx *wrap.Tx, e *RefreshToken) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.UserId, e.RefreshToken, e.IsLogout, e.LogoutTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *RefreshTokenDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *RefreshTokenDao) scanRow(row *wrap.Row) (*RefreshToken, error) {
	e := &RefreshToken{}
	err := row.Scan(&e.Id, &e.UserId, &e.RefreshToken, &e.IsLogout, &e.LogoutTime, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *RefreshTokenDao) scanRows(rows *wrap.Rows) (list []*RefreshToken, err error) {
	list = make([]*RefreshToken, 0)
	for rows.Next() {
		e := RefreshToken{}
		err = rows.Scan(&e.Id, &e.UserId, &e.RefreshToken, &e.IsLogout, &e.LogoutTime, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *RefreshTokenDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*RefreshToken, error) {
	querySql := "SELECT " + REFRESH_TOKEN_ALL_FIELDS_STRING + " FROM refresh_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *RefreshTokenDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*RefreshToken, err error) {
	querySql := "SELECT " + REFRESH_TOKEN_ALL_FIELDS_STRING + " FROM refresh_token " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *RefreshTokenDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM refresh_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *RefreshTokenDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM refresh_token " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *RefreshTokenDao) GetQuery() *RefreshTokenQuery {
	return NewRefreshTokenQuery(dao)
}

const USER_TABLE_NAME = "user"

type USER_FIELD string

const USER_FIELD_ID = USER_FIELD("id")
const USER_FIELD_USER_ID = USER_FIELD("user_id")
const USER_FIELD_USER_NAME = USER_FIELD("user_name")
const USER_FIELD_USER_ICON = USER_FIELD("user_icon")
const USER_FIELD_CREATE_TIME = USER_FIELD("create_time")
const USER_FIELD_UPDATE_TIME = USER_FIELD("update_time")

const USER_ALL_FIELDS_STRING = "id,user_id,user_name,user_icon,create_time,update_time"

var USER_ALL_FIELDS = []string{
	"id",
	"user_id",
	"user_name",
	"user_icon",
	"create_time",
	"update_time",
}

type User struct {
	Id         uint64 //size=20
	UserId     string //size=32
	UserName   string //size=32
	UserIcon   string //size=256
	CreateTime time.Time
	UpdateTime time.Time
}

type UserQuery struct {
	BaseQuery
	dao *UserDao
}

func NewUserQuery(dao *UserDao) *UserQuery {
	q := &UserQuery{}
	q.dao = dao

	return q
}

func (q *UserQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*User, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *UserQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*User, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *UserQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *UserQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *UserQuery) ForUpdate() *UserQuery {
	q.forUpdate = true
	return q
}

func (q *UserQuery) ForShare() *UserQuery {
	q.forShare = true
	return q
}

func (q *UserQuery) GroupBy(fields ...USER_FIELD) *UserQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *UserQuery) Limit(startIncluded int64, count int64) *UserQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *UserQuery) OrderBy(fieldName USER_FIELD, asc bool) *UserQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *UserQuery) OrderByGroupCount(asc bool) *UserQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *UserQuery) w(format string, a ...interface{}) *UserQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *UserQuery) Left() *UserQuery  { return q.w(" ( ") }
func (q *UserQuery) Right() *UserQuery { return q.w(" ) ") }
func (q *UserQuery) And() *UserQuery   { return q.w(" AND ") }
func (q *UserQuery) Or() *UserQuery    { return q.w(" OR ") }
func (q *UserQuery) Not() *UserQuery   { return q.w(" NOT ") }

func (q *UserQuery) Id_Equal(v uint64) *UserQuery        { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *UserQuery) Id_NotEqual(v uint64) *UserQuery     { return q.w("id<>'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) Id_Less(v uint64) *UserQuery         { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) Id_LessEqual(v uint64) *UserQuery    { return q.w("id<='" + fmt.Sprint(v) + "'") }
func (q *UserQuery) Id_Greater(v uint64) *UserQuery      { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) Id_GreaterEqual(v uint64) *UserQuery { return q.w("id>='" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserId_Equal(v string) *UserQuery    { return q.w("user_id='" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserId_NotEqual(v string) *UserQuery { return q.w("user_id<>'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserId_Less(v string) *UserQuery     { return q.w("user_id<'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserId_LessEqual(v string) *UserQuery {
	return q.w("user_id<='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserId_Greater(v string) *UserQuery { return q.w("user_id>'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserId_GreaterEqual(v string) *UserQuery {
	return q.w("user_id>='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserName_Equal(v string) *UserQuery { return q.w("user_name='" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserName_NotEqual(v string) *UserQuery {
	return q.w("user_name<>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserName_Less(v string) *UserQuery { return q.w("user_name<'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserName_LessEqual(v string) *UserQuery {
	return q.w("user_name<='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserName_Greater(v string) *UserQuery {
	return q.w("user_name>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserName_GreaterEqual(v string) *UserQuery {
	return q.w("user_name>='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserIcon_Equal(v string) *UserQuery { return q.w("user_icon='" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserIcon_NotEqual(v string) *UserQuery {
	return q.w("user_icon<>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserIcon_Less(v string) *UserQuery { return q.w("user_icon<'" + fmt.Sprint(v) + "'") }
func (q *UserQuery) UserIcon_LessEqual(v string) *UserQuery {
	return q.w("user_icon<='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserIcon_Greater(v string) *UserQuery {
	return q.w("user_icon>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UserIcon_GreaterEqual(v string) *UserQuery {
	return q.w("user_icon>='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) CreateTime_Equal(v time.Time) *UserQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) CreateTime_NotEqual(v time.Time) *UserQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) CreateTime_Less(v time.Time) *UserQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) CreateTime_LessEqual(v time.Time) *UserQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) CreateTime_Greater(v time.Time) *UserQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) CreateTime_GreaterEqual(v time.Time) *UserQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UpdateTime_Equal(v time.Time) *UserQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UpdateTime_NotEqual(v time.Time) *UserQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UpdateTime_Less(v time.Time) *UserQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UpdateTime_LessEqual(v time.Time) *UserQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UpdateTime_Greater(v time.Time) *UserQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *UserQuery) UpdateTime_GreaterEqual(v time.Time) *UserQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type UserDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewUserDao(db *DB) (t *UserDao, err error) {
	t = &UserDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *UserDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO user (user_id,user_name,user_icon) VALUES (?,?,?)")
	return err
}

func (dao *UserDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE user SET user_id=?,user_name=?,user_icon=? WHERE id=?")
	return err
}

func (dao *UserDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM user WHERE id=?")
	return err
}

func (dao *UserDao) Insert(ctx context.Context, tx *wrap.Tx, e *User) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.UserId, e.UserName, e.UserIcon)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *UserDao) Update(ctx context.Context, tx *wrap.Tx, e *User) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.UserId, e.UserName, e.UserIcon, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserDao) scanRow(row *wrap.Row) (*User, error) {
	e := &User{}
	err := row.Scan(&e.Id, &e.UserId, &e.UserName, &e.UserIcon, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *UserDao) scanRows(rows *wrap.Rows) (list []*User, err error) {
	list = make([]*User, 0)
	for rows.Next() {
		e := User{}
		err = rows.Scan(&e.Id, &e.UserId, &e.UserName, &e.UserIcon, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *UserDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*User, error) {
	querySql := "SELECT " + USER_ALL_FIELDS_STRING + " FROM user " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *UserDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*User, err error) {
	querySql := "SELECT " + USER_ALL_FIELDS_STRING + " FROM user " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *UserDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM user " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *UserDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM user " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *UserDao) GetQuery() *UserQuery {
	return NewUserQuery(dao)
}

type DB struct {
	wrap.DB
	AccessToken      *AccessTokenDao
	AccountOperation *AccountOperationDao
	LoginSmsCode     *LoginSmsCodeDao
	OauthAccount     *OauthAccountDao
	OauthState       *OauthStateDao
	PhoneAccount     *PhoneAccountDao
	RefreshToken     *RefreshTokenDao
	User             *UserDao
}

func NewDB() (d *DB, err error) {
	d = &DB{}

	connectionString := os.Getenv("DB")
	if connectionString == "" {
		return nil, fmt.Errorf("DB env nil")
	}
	connectionString += "/neuron-user?parseTime=true"
	db, err := wrap.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	d.DB = *db

	err = d.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	d.AccessToken, err = NewAccessTokenDao(d)
	if err != nil {
		return nil, err
	}

	d.AccountOperation, err = NewAccountOperationDao(d)
	if err != nil {
		return nil, err
	}

	d.LoginSmsCode, err = NewLoginSmsCodeDao(d)
	if err != nil {
		return nil, err
	}

	d.OauthAccount, err = NewOauthAccountDao(d)
	if err != nil {
		return nil, err
	}

	d.OauthState, err = NewOauthStateDao(d)
	if err != nil {
		return nil, err
	}

	d.PhoneAccount, err = NewPhoneAccountDao(d)
	if err != nil {
		return nil, err
	}

	d.RefreshToken, err = NewRefreshTokenDao(d)
	if err != nil {
		return nil, err
	}

	d.User, err = NewUserDao(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
