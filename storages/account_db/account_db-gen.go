package account_db

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

const ACCOUNT_TABLE_NAME = "account"

type ACCOUNT_FIELD string

const ACCOUNT_FIELD_ID = ACCOUNT_FIELD("id")
const ACCOUNT_FIELD_ACCOUNT_ID = ACCOUNT_FIELD("account_id")
const ACCOUNT_FIELD_PHONE_NUMBER = ACCOUNT_FIELD("phone_number")
const ACCOUNT_FIELD_EMAIL_ADDRESS = ACCOUNT_FIELD("email_address")
const ACCOUNT_FIELD_PASSWORD_HASH = ACCOUNT_FIELD("password_hash")
const ACCOUNT_FIELD_OAUTH_PROVIDER = ACCOUNT_FIELD("oauth_provider")
const ACCOUNT_FIELD_OAUTH_ACCOUNT_ID = ACCOUNT_FIELD("oauth_account_id")
const ACCOUNT_FIELD_CREATE_TIME = ACCOUNT_FIELD("create_time")
const ACCOUNT_FIELD_UPDATE_TIME = ACCOUNT_FIELD("update_time")
const ACCOUNT_FIELD_UPDATE_VERSION = ACCOUNT_FIELD("update_version")

const ACCOUNT_ALL_FIELDS_STRING = "id,account_id,phone_number,email_address,password_hash,oauth_provider,oauth_account_id,create_time,update_time,update_version"

var ACCOUNT_ALL_FIELDS = []string{
	"id",
	"account_id",
	"phone_number",
	"email_address",
	"password_hash",
	"oauth_provider",
	"oauth_account_id",
	"create_time",
	"update_time",
	"update_version",
}

type Account struct {
	Id             uint64         //size=20
	AccountId      string         //size=128
	PhoneNumber    sql.NullString //size=32
	EmailAddress   sql.NullString //size=128
	PasswordHash   string         //size=128
	OauthProvider  sql.NullString //size=128
	OauthAccountId sql.NullString //size=128
	CreateTime     time.Time
	UpdateTime     time.Time
	UpdateVersion  uint64 //size=20
}

type AccountQuery struct {
	BaseQuery
	dao *AccountDao
}

func NewAccountQuery(dao *AccountDao) *AccountQuery {
	q := &AccountQuery{}
	q.dao = dao

	return q
}

func (q *AccountQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*Account, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *AccountQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*Account, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *AccountQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *AccountQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *AccountQuery) ForUpdate() *AccountQuery {
	q.forUpdate = true
	return q
}

func (q *AccountQuery) ForShare() *AccountQuery {
	q.forShare = true
	return q
}

func (q *AccountQuery) GroupBy(fields ...ACCOUNT_FIELD) *AccountQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *AccountQuery) Limit(startIncluded int64, count int64) *AccountQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *AccountQuery) OrderBy(fieldName ACCOUNT_FIELD, asc bool) *AccountQuery {
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

func (q *AccountQuery) OrderByGroupCount(asc bool) *AccountQuery {
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

func (q *AccountQuery) w(format string, a ...interface{}) *AccountQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *AccountQuery) Left() *AccountQuery  { return q.w(" ( ") }
func (q *AccountQuery) Right() *AccountQuery { return q.w(" ) ") }
func (q *AccountQuery) And() *AccountQuery   { return q.w(" AND ") }
func (q *AccountQuery) Or() *AccountQuery    { return q.w(" OR ") }
func (q *AccountQuery) Not() *AccountQuery   { return q.w(" NOT ") }

func (q *AccountQuery) Id_Equal(v uint64) *AccountQuery     { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_NotEqual(v uint64) *AccountQuery  { return q.w("id<>'" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_Less(v uint64) *AccountQuery      { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_LessEqual(v uint64) *AccountQuery { return q.w("id<='" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_Greater(v uint64) *AccountQuery   { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_GreaterEqual(v uint64) *AccountQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) AccountId_Equal(v string) *AccountQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) AccountId_NotEqual(v string) *AccountQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) AccountId_Less(v string) *AccountQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) AccountId_LessEqual(v string) *AccountQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) AccountId_Greater(v string) *AccountQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) AccountId_GreaterEqual(v string) *AccountQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PhoneNumber_Equal(v string) *AccountQuery {
	return q.w("phone_number='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PhoneNumber_NotEqual(v string) *AccountQuery {
	return q.w("phone_number<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PhoneNumber_Less(v string) *AccountQuery {
	return q.w("phone_number<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PhoneNumber_LessEqual(v string) *AccountQuery {
	return q.w("phone_number<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PhoneNumber_Greater(v string) *AccountQuery {
	return q.w("phone_number>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PhoneNumber_GreaterEqual(v string) *AccountQuery {
	return q.w("phone_number>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PhoneNumber_IsNull() *AccountQuery  { return q.w("phone_number IS NULL") }
func (q *AccountQuery) PhoneNumber_NotNull() *AccountQuery { return q.w("phone_number IS NOT NULL") }
func (q *AccountQuery) EmailAddress_Equal(v string) *AccountQuery {
	return q.w("email_address='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) EmailAddress_NotEqual(v string) *AccountQuery {
	return q.w("email_address<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) EmailAddress_Less(v string) *AccountQuery {
	return q.w("email_address<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) EmailAddress_LessEqual(v string) *AccountQuery {
	return q.w("email_address<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) EmailAddress_Greater(v string) *AccountQuery {
	return q.w("email_address>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) EmailAddress_GreaterEqual(v string) *AccountQuery {
	return q.w("email_address>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) EmailAddress_IsNull() *AccountQuery  { return q.w("email_address IS NULL") }
func (q *AccountQuery) EmailAddress_NotNull() *AccountQuery { return q.w("email_address IS NOT NULL") }
func (q *AccountQuery) PasswordHash_Equal(v string) *AccountQuery {
	return q.w("password_hash='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PasswordHash_NotEqual(v string) *AccountQuery {
	return q.w("password_hash<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PasswordHash_Less(v string) *AccountQuery {
	return q.w("password_hash<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PasswordHash_LessEqual(v string) *AccountQuery {
	return q.w("password_hash<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PasswordHash_Greater(v string) *AccountQuery {
	return q.w("password_hash>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) PasswordHash_GreaterEqual(v string) *AccountQuery {
	return q.w("password_hash>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthProvider_Equal(v string) *AccountQuery {
	return q.w("oauth_provider='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthProvider_NotEqual(v string) *AccountQuery {
	return q.w("oauth_provider<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthProvider_Less(v string) *AccountQuery {
	return q.w("oauth_provider<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthProvider_LessEqual(v string) *AccountQuery {
	return q.w("oauth_provider<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthProvider_Greater(v string) *AccountQuery {
	return q.w("oauth_provider>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthProvider_GreaterEqual(v string) *AccountQuery {
	return q.w("oauth_provider>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthProvider_IsNull() *AccountQuery  { return q.w("oauth_provider IS NULL") }
func (q *AccountQuery) OauthProvider_NotNull() *AccountQuery { return q.w("oauth_provider IS NOT NULL") }
func (q *AccountQuery) OauthAccountId_Equal(v string) *AccountQuery {
	return q.w("oauth_account_id='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthAccountId_NotEqual(v string) *AccountQuery {
	return q.w("oauth_account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthAccountId_Less(v string) *AccountQuery {
	return q.w("oauth_account_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthAccountId_LessEqual(v string) *AccountQuery {
	return q.w("oauth_account_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthAccountId_Greater(v string) *AccountQuery {
	return q.w("oauth_account_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthAccountId_GreaterEqual(v string) *AccountQuery {
	return q.w("oauth_account_id>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) OauthAccountId_IsNull() *AccountQuery { return q.w("oauth_account_id IS NULL") }
func (q *AccountQuery) OauthAccountId_NotNull() *AccountQuery {
	return q.w("oauth_account_id IS NOT NULL")
}
func (q *AccountQuery) CreateTime_Equal(v time.Time) *AccountQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) CreateTime_NotEqual(v time.Time) *AccountQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) CreateTime_Less(v time.Time) *AccountQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) CreateTime_LessEqual(v time.Time) *AccountQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) CreateTime_Greater(v time.Time) *AccountQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) CreateTime_GreaterEqual(v time.Time) *AccountQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateTime_Equal(v time.Time) *AccountQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateTime_NotEqual(v time.Time) *AccountQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateTime_Less(v time.Time) *AccountQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateTime_LessEqual(v time.Time) *AccountQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateTime_Greater(v time.Time) *AccountQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateTime_GreaterEqual(v time.Time) *AccountQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateVersion_Equal(v uint64) *AccountQuery {
	return q.w("update_version='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateVersion_NotEqual(v uint64) *AccountQuery {
	return q.w("update_version<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateVersion_Less(v uint64) *AccountQuery {
	return q.w("update_version<'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateVersion_LessEqual(v uint64) *AccountQuery {
	return q.w("update_version<='" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateVersion_Greater(v uint64) *AccountQuery {
	return q.w("update_version>'" + fmt.Sprint(v) + "'")
}
func (q *AccountQuery) UpdateVersion_GreaterEqual(v uint64) *AccountQuery {
	return q.w("update_version>='" + fmt.Sprint(v) + "'")
}

type AccountDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewAccountDao(db *DB) (t *AccountDao, err error) {
	t = &AccountDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *AccountDao) init() (err error) {
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

func (dao *AccountDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO account (account_id,phone_number,email_address,password_hash,oauth_provider,oauth_account_id,update_version) VALUES (?,?,?,?,?,?,?)")
	return err
}

func (dao *AccountDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE account SET account_id=?,phone_number=?,email_address=?,password_hash=?,oauth_provider=?,oauth_account_id=?,update_version=update_version+1 WHERE id=? AND update_version=?")
	return err
}

func (dao *AccountDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM account WHERE id=?")
	return err
}

func (dao *AccountDao) Insert(ctx context.Context, tx *wrap.Tx, e *Account) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.AccountId, e.PhoneNumber, e.EmailAddress, e.PasswordHash, e.OauthProvider, e.OauthAccountId, e.UpdateVersion)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *AccountDao) Update(ctx context.Context, tx *wrap.Tx, e *Account) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.AccountId, e.PhoneNumber, e.EmailAddress, e.PasswordHash, e.OauthProvider, e.OauthAccountId, e.Id, e.UpdateVersion)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
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

func (dao *AccountDao) scanRow(row *wrap.Row) (*Account, error) {
	e := &Account{}
	err := row.Scan(&e.Id, &e.AccountId, &e.PhoneNumber, &e.EmailAddress, &e.PasswordHash, &e.OauthProvider, &e.OauthAccountId, &e.CreateTime, &e.UpdateTime, &e.UpdateVersion)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *AccountDao) scanRows(rows *wrap.Rows) (list []*Account, err error) {
	list = make([]*Account, 0)
	for rows.Next() {
		e := Account{}
		err = rows.Scan(&e.Id, &e.AccountId, &e.PhoneNumber, &e.EmailAddress, &e.PasswordHash, &e.OauthProvider, &e.OauthAccountId, &e.CreateTime, &e.UpdateTime, &e.UpdateVersion)
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

func (dao *AccountDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*Account, error) {
	querySql := "SELECT " + ACCOUNT_ALL_FIELDS_STRING + " FROM account " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *AccountDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*Account, err error) {
	querySql := "SELECT " + ACCOUNT_ALL_FIELDS_STRING + " FROM account " + query
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

func (dao *AccountDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM account " + query
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

func (dao *AccountDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM account " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *AccountDao) GetQuery() *AccountQuery {
	return NewAccountQuery(dao)
}

const ACCOUNT_OPERATION_TABLE_NAME = "account_operation"

type ACCOUNT_OPERATION_FIELD string

const ACCOUNT_OPERATION_FIELD_ID = ACCOUNT_OPERATION_FIELD("id")
const ACCOUNT_OPERATION_FIELD_CREATE_TIME = ACCOUNT_OPERATION_FIELD("create_time")
const ACCOUNT_OPERATION_FIELD_OPERATION_TYPE = ACCOUNT_OPERATION_FIELD("operation_type")
const ACCOUNT_OPERATION_FIELD_USER_AGENT = ACCOUNT_OPERATION_FIELD("user_agent")
const ACCOUNT_OPERATION_FIELD_ERROR_STATUS = ACCOUNT_OPERATION_FIELD("error_status")
const ACCOUNT_OPERATION_FIELD_ERROR_CODE = ACCOUNT_OPERATION_FIELD("error_code")
const ACCOUNT_OPERATION_FIELD_ERROR_MESSAGE = ACCOUNT_OPERATION_FIELD("error_message")
const ACCOUNT_OPERATION_FIELD_SMS_SCENE = ACCOUNT_OPERATION_FIELD("sms_scene")
const ACCOUNT_OPERATION_FIELD_PHONE_NUMBER = ACCOUNT_OPERATION_FIELD("phone_number")
const ACCOUNT_OPERATION_FIELD_LOGIN_NAME = ACCOUNT_OPERATION_FIELD("login_name")
const ACCOUNT_OPERATION_FIELD_ACCOUNT_ID = ACCOUNT_OPERATION_FIELD("account_id")

const ACCOUNT_OPERATION_ALL_FIELDS_STRING = "id,create_time,operation_type,user_agent,error_status,error_code,error_message,sms_scene,phone_number,login_name,account_id"

var ACCOUNT_OPERATION_ALL_FIELDS = []string{
	"id",
	"create_time",
	"operation_type",
	"user_agent",
	"error_status",
	"error_code",
	"error_message",
	"sms_scene",
	"phone_number",
	"login_name",
	"account_id",
}

type AccountOperation struct {
	Id            uint64 //size=20
	CreateTime    time.Time
	OperationType string //size=32
	UserAgent     string //size=256
	ErrorStatus   int32  //size=11
	ErrorCode     string //size=32
	ErrorMessage  string //size=256
	SmsScene      string //size=32
	PhoneNumber   string //size=16
	LoginName     string //size=64
	AccountId     string //size=128
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
func (q *AccountOperationQuery) OperationType_Equal(v string) *AccountOperationQuery {
	return q.w("operation_type='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_NotEqual(v string) *AccountOperationQuery {
	return q.w("operation_type<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_Less(v string) *AccountOperationQuery {
	return q.w("operation_type<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_LessEqual(v string) *AccountOperationQuery {
	return q.w("operation_type<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_Greater(v string) *AccountOperationQuery {
	return q.w("operation_type>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("operation_type>='" + fmt.Sprint(v) + "'")
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
func (q *AccountOperationQuery) ErrorStatus_Equal(v int32) *AccountOperationQuery {
	return q.w("error_status='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorStatus_NotEqual(v int32) *AccountOperationQuery {
	return q.w("error_status<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorStatus_Less(v int32) *AccountOperationQuery {
	return q.w("error_status<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorStatus_LessEqual(v int32) *AccountOperationQuery {
	return q.w("error_status<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorStatus_Greater(v int32) *AccountOperationQuery {
	return q.w("error_status>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorStatus_GreaterEqual(v int32) *AccountOperationQuery {
	return q.w("error_status>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorCode_Equal(v string) *AccountOperationQuery {
	return q.w("error_code='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorCode_NotEqual(v string) *AccountOperationQuery {
	return q.w("error_code<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorCode_Less(v string) *AccountOperationQuery {
	return q.w("error_code<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorCode_LessEqual(v string) *AccountOperationQuery {
	return q.w("error_code<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorCode_Greater(v string) *AccountOperationQuery {
	return q.w("error_code>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorCode_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("error_code>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorMessage_Equal(v string) *AccountOperationQuery {
	return q.w("error_message='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorMessage_NotEqual(v string) *AccountOperationQuery {
	return q.w("error_message<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorMessage_Less(v string) *AccountOperationQuery {
	return q.w("error_message<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorMessage_LessEqual(v string) *AccountOperationQuery {
	return q.w("error_message<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorMessage_Greater(v string) *AccountOperationQuery {
	return q.w("error_message>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) ErrorMessage_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("error_message>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_Equal(v string) *AccountOperationQuery {
	return q.w("sms_scene='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_NotEqual(v string) *AccountOperationQuery {
	return q.w("sms_scene<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_Less(v string) *AccountOperationQuery {
	return q.w("sms_scene<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_LessEqual(v string) *AccountOperationQuery {
	return q.w("sms_scene<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_Greater(v string) *AccountOperationQuery {
	return q.w("sms_scene>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("sms_scene>='" + fmt.Sprint(v) + "'")
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
func (q *AccountOperationQuery) LoginName_Equal(v string) *AccountOperationQuery {
	return q.w("login_name='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) LoginName_NotEqual(v string) *AccountOperationQuery {
	return q.w("login_name<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) LoginName_Less(v string) *AccountOperationQuery {
	return q.w("login_name<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) LoginName_LessEqual(v string) *AccountOperationQuery {
	return q.w("login_name<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) LoginName_Greater(v string) *AccountOperationQuery {
	return q.w("login_name>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) LoginName_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("login_name>='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) AccountId_Equal(v string) *AccountOperationQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) AccountId_NotEqual(v string) *AccountOperationQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) AccountId_Less(v string) *AccountOperationQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) AccountId_LessEqual(v string) *AccountOperationQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) AccountId_Greater(v string) *AccountOperationQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) AccountId_GreaterEqual(v string) *AccountOperationQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
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
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO account_operation (operation_type,user_agent,error_status,error_code,error_message,sms_scene,phone_number,login_name,account_id) VALUES (?,?,?,?,?,?,?,?,?)")
	return err
}

func (dao *AccountOperationDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE account_operation SET operation_type=?,user_agent=?,error_status=?,error_code=?,error_message=?,sms_scene=?,phone_number=?,login_name=?,account_id=? WHERE id=?")
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

	result, err := stmt.Exec(ctx, e.OperationType, e.UserAgent, e.ErrorStatus, e.ErrorCode, e.ErrorMessage, e.SmsScene, e.PhoneNumber, e.LoginName, e.AccountId)
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

	_, err = stmt.Exec(ctx, e.OperationType, e.UserAgent, e.ErrorStatus, e.ErrorCode, e.ErrorMessage, e.SmsScene, e.PhoneNumber, e.LoginName, e.AccountId, e.Id)
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
	err := row.Scan(&e.Id, &e.CreateTime, &e.OperationType, &e.UserAgent, &e.ErrorStatus, &e.ErrorCode, &e.ErrorMessage, &e.SmsScene, &e.PhoneNumber, &e.LoginName, &e.AccountId)
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
		err = rows.Scan(&e.Id, &e.CreateTime, &e.OperationType, &e.UserAgent, &e.ErrorStatus, &e.ErrorCode, &e.ErrorMessage, &e.SmsScene, &e.PhoneNumber, &e.LoginName, &e.AccountId)
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

const SMS_CODE_TABLE_NAME = "sms_code"

type SMS_CODE_FIELD string

const SMS_CODE_FIELD_ID = SMS_CODE_FIELD("id")
const SMS_CODE_FIELD_SCENE_TYPE = SMS_CODE_FIELD("scene_type")
const SMS_CODE_FIELD_PHONE_NUMBER = SMS_CODE_FIELD("phone_number")
const SMS_CODE_FIELD_SMS_CODE = SMS_CODE_FIELD("sms_code")
const SMS_CODE_FIELD_CREATE_TIME = SMS_CODE_FIELD("create_time")

const SMS_CODE_ALL_FIELDS_STRING = "id,scene_type,phone_number,sms_code,create_time"

var SMS_CODE_ALL_FIELDS = []string{
	"id",
	"scene_type",
	"phone_number",
	"sms_code",
	"create_time",
}

type SmsCode struct {
	Id          uint64 //size=20
	SceneType   string //size=32
	PhoneNumber string //size=32
	SmsCode     string //size=32
	CreateTime  time.Time
}

type SmsCodeQuery struct {
	BaseQuery
	dao *SmsCodeDao
}

func NewSmsCodeQuery(dao *SmsCodeDao) *SmsCodeQuery {
	q := &SmsCodeQuery{}
	q.dao = dao

	return q
}

func (q *SmsCodeQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*SmsCode, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *SmsCodeQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*SmsCode, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *SmsCodeQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *SmsCodeQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *SmsCodeQuery) ForUpdate() *SmsCodeQuery {
	q.forUpdate = true
	return q
}

func (q *SmsCodeQuery) ForShare() *SmsCodeQuery {
	q.forShare = true
	return q
}

func (q *SmsCodeQuery) GroupBy(fields ...SMS_CODE_FIELD) *SmsCodeQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *SmsCodeQuery) Limit(startIncluded int64, count int64) *SmsCodeQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *SmsCodeQuery) OrderBy(fieldName SMS_CODE_FIELD, asc bool) *SmsCodeQuery {
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

func (q *SmsCodeQuery) OrderByGroupCount(asc bool) *SmsCodeQuery {
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

func (q *SmsCodeQuery) w(format string, a ...interface{}) *SmsCodeQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *SmsCodeQuery) Left() *SmsCodeQuery  { return q.w(" ( ") }
func (q *SmsCodeQuery) Right() *SmsCodeQuery { return q.w(" ) ") }
func (q *SmsCodeQuery) And() *SmsCodeQuery   { return q.w(" AND ") }
func (q *SmsCodeQuery) Or() *SmsCodeQuery    { return q.w(" OR ") }
func (q *SmsCodeQuery) Not() *SmsCodeQuery   { return q.w(" NOT ") }

func (q *SmsCodeQuery) Id_Equal(v uint64) *SmsCodeQuery     { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_NotEqual(v uint64) *SmsCodeQuery  { return q.w("id<>'" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_Less(v uint64) *SmsCodeQuery      { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_LessEqual(v uint64) *SmsCodeQuery { return q.w("id<='" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_Greater(v uint64) *SmsCodeQuery   { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_GreaterEqual(v uint64) *SmsCodeQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SceneType_Equal(v string) *SmsCodeQuery {
	return q.w("scene_type='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SceneType_NotEqual(v string) *SmsCodeQuery {
	return q.w("scene_type<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SceneType_Less(v string) *SmsCodeQuery {
	return q.w("scene_type<'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SceneType_LessEqual(v string) *SmsCodeQuery {
	return q.w("scene_type<='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SceneType_Greater(v string) *SmsCodeQuery {
	return q.w("scene_type>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SceneType_GreaterEqual(v string) *SmsCodeQuery {
	return q.w("scene_type>='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneNumber_Equal(v string) *SmsCodeQuery {
	return q.w("phone_number='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneNumber_NotEqual(v string) *SmsCodeQuery {
	return q.w("phone_number<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneNumber_Less(v string) *SmsCodeQuery {
	return q.w("phone_number<'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneNumber_LessEqual(v string) *SmsCodeQuery {
	return q.w("phone_number<='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneNumber_Greater(v string) *SmsCodeQuery {
	return q.w("phone_number>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneNumber_GreaterEqual(v string) *SmsCodeQuery {
	return q.w("phone_number>='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_Equal(v string) *SmsCodeQuery {
	return q.w("sms_code='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_NotEqual(v string) *SmsCodeQuery {
	return q.w("sms_code<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_Less(v string) *SmsCodeQuery {
	return q.w("sms_code<'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_LessEqual(v string) *SmsCodeQuery {
	return q.w("sms_code<='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_Greater(v string) *SmsCodeQuery {
	return q.w("sms_code>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_GreaterEqual(v string) *SmsCodeQuery {
	return q.w("sms_code>='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) CreateTime_Equal(v time.Time) *SmsCodeQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) CreateTime_NotEqual(v time.Time) *SmsCodeQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) CreateTime_Less(v time.Time) *SmsCodeQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) CreateTime_LessEqual(v time.Time) *SmsCodeQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) CreateTime_Greater(v time.Time) *SmsCodeQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) CreateTime_GreaterEqual(v time.Time) *SmsCodeQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}

type SmsCodeDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewSmsCodeDao(db *DB) (t *SmsCodeDao, err error) {
	t = &SmsCodeDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *SmsCodeDao) init() (err error) {
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

func (dao *SmsCodeDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO sms_code (scene_type,phone_number,sms_code) VALUES (?,?,?)")
	return err
}

func (dao *SmsCodeDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE sms_code SET scene_type=?,phone_number=?,sms_code=? WHERE id=?")
	return err
}

func (dao *SmsCodeDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM sms_code WHERE id=?")
	return err
}

func (dao *SmsCodeDao) Insert(ctx context.Context, tx *wrap.Tx, e *SmsCode) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.SceneType, e.PhoneNumber, e.SmsCode)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *SmsCodeDao) Update(ctx context.Context, tx *wrap.Tx, e *SmsCode) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.SceneType, e.PhoneNumber, e.SmsCode, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *SmsCodeDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
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

func (dao *SmsCodeDao) scanRow(row *wrap.Row) (*SmsCode, error) {
	e := &SmsCode{}
	err := row.Scan(&e.Id, &e.SceneType, &e.PhoneNumber, &e.SmsCode, &e.CreateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *SmsCodeDao) scanRows(rows *wrap.Rows) (list []*SmsCode, err error) {
	list = make([]*SmsCode, 0)
	for rows.Next() {
		e := SmsCode{}
		err = rows.Scan(&e.Id, &e.SceneType, &e.PhoneNumber, &e.SmsCode, &e.CreateTime)
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

func (dao *SmsCodeDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*SmsCode, error) {
	querySql := "SELECT " + SMS_CODE_ALL_FIELDS_STRING + " FROM sms_code " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *SmsCodeDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*SmsCode, err error) {
	querySql := "SELECT " + SMS_CODE_ALL_FIELDS_STRING + " FROM sms_code " + query
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

func (dao *SmsCodeDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM sms_code " + query
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

func (dao *SmsCodeDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM sms_code " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *SmsCodeDao) GetQuery() *SmsCodeQuery {
	return NewSmsCodeQuery(dao)
}

const SMS_SCENE_TABLE_NAME = "sms_scene"

type SMS_SCENE_FIELD string

const SMS_SCENE_FIELD_ID = SMS_SCENE_FIELD("id")
const SMS_SCENE_FIELD_SCENE_TYPE = SMS_SCENE_FIELD("scene_type")
const SMS_SCENE_FIELD_SCENE_DESC = SMS_SCENE_FIELD("scene_desc")
const SMS_SCENE_FIELD_CREATE_TIME = SMS_SCENE_FIELD("create_time")
const SMS_SCENE_FIELD_UPDATE_TIME = SMS_SCENE_FIELD("update_time")

const SMS_SCENE_ALL_FIELDS_STRING = "id,scene_type,scene_desc,create_time,update_time"

var SMS_SCENE_ALL_FIELDS = []string{
	"id",
	"scene_type",
	"scene_desc",
	"create_time",
	"update_time",
}

type SmsScene struct {
	Id         uint64 //size=20
	SceneType  string //size=32
	SceneDesc  string //size=32
	CreateTime time.Time
	UpdateTime time.Time
}

type SmsSceneQuery struct {
	BaseQuery
	dao *SmsSceneDao
}

func NewSmsSceneQuery(dao *SmsSceneDao) *SmsSceneQuery {
	q := &SmsSceneQuery{}
	q.dao = dao

	return q
}

func (q *SmsSceneQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*SmsScene, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *SmsSceneQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*SmsScene, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *SmsSceneQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *SmsSceneQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *SmsSceneQuery) ForUpdate() *SmsSceneQuery {
	q.forUpdate = true
	return q
}

func (q *SmsSceneQuery) ForShare() *SmsSceneQuery {
	q.forShare = true
	return q
}

func (q *SmsSceneQuery) GroupBy(fields ...SMS_SCENE_FIELD) *SmsSceneQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *SmsSceneQuery) Limit(startIncluded int64, count int64) *SmsSceneQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *SmsSceneQuery) OrderBy(fieldName SMS_SCENE_FIELD, asc bool) *SmsSceneQuery {
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

func (q *SmsSceneQuery) OrderByGroupCount(asc bool) *SmsSceneQuery {
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

func (q *SmsSceneQuery) w(format string, a ...interface{}) *SmsSceneQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *SmsSceneQuery) Left() *SmsSceneQuery  { return q.w(" ( ") }
func (q *SmsSceneQuery) Right() *SmsSceneQuery { return q.w(" ) ") }
func (q *SmsSceneQuery) And() *SmsSceneQuery   { return q.w(" AND ") }
func (q *SmsSceneQuery) Or() *SmsSceneQuery    { return q.w(" OR ") }
func (q *SmsSceneQuery) Not() *SmsSceneQuery   { return q.w(" NOT ") }

func (q *SmsSceneQuery) Id_Equal(v uint64) *SmsSceneQuery     { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *SmsSceneQuery) Id_NotEqual(v uint64) *SmsSceneQuery  { return q.w("id<>'" + fmt.Sprint(v) + "'") }
func (q *SmsSceneQuery) Id_Less(v uint64) *SmsSceneQuery      { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *SmsSceneQuery) Id_LessEqual(v uint64) *SmsSceneQuery { return q.w("id<='" + fmt.Sprint(v) + "'") }
func (q *SmsSceneQuery) Id_Greater(v uint64) *SmsSceneQuery   { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *SmsSceneQuery) Id_GreaterEqual(v uint64) *SmsSceneQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneType_Equal(v string) *SmsSceneQuery {
	return q.w("scene_type='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneType_NotEqual(v string) *SmsSceneQuery {
	return q.w("scene_type<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneType_Less(v string) *SmsSceneQuery {
	return q.w("scene_type<'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneType_LessEqual(v string) *SmsSceneQuery {
	return q.w("scene_type<='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneType_Greater(v string) *SmsSceneQuery {
	return q.w("scene_type>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneType_GreaterEqual(v string) *SmsSceneQuery {
	return q.w("scene_type>='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneDesc_Equal(v string) *SmsSceneQuery {
	return q.w("scene_desc='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneDesc_NotEqual(v string) *SmsSceneQuery {
	return q.w("scene_desc<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneDesc_Less(v string) *SmsSceneQuery {
	return q.w("scene_desc<'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneDesc_LessEqual(v string) *SmsSceneQuery {
	return q.w("scene_desc<='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneDesc_Greater(v string) *SmsSceneQuery {
	return q.w("scene_desc>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) SceneDesc_GreaterEqual(v string) *SmsSceneQuery {
	return q.w("scene_desc>='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) CreateTime_Equal(v time.Time) *SmsSceneQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) CreateTime_NotEqual(v time.Time) *SmsSceneQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) CreateTime_Less(v time.Time) *SmsSceneQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) CreateTime_LessEqual(v time.Time) *SmsSceneQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) CreateTime_Greater(v time.Time) *SmsSceneQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) CreateTime_GreaterEqual(v time.Time) *SmsSceneQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) UpdateTime_Equal(v time.Time) *SmsSceneQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) UpdateTime_NotEqual(v time.Time) *SmsSceneQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) UpdateTime_Less(v time.Time) *SmsSceneQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) UpdateTime_LessEqual(v time.Time) *SmsSceneQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) UpdateTime_Greater(v time.Time) *SmsSceneQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *SmsSceneQuery) UpdateTime_GreaterEqual(v time.Time) *SmsSceneQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type SmsSceneDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewSmsSceneDao(db *DB) (t *SmsSceneDao, err error) {
	t = &SmsSceneDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *SmsSceneDao) init() (err error) {
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

func (dao *SmsSceneDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO sms_scene (scene_type,scene_desc) VALUES (?,?)")
	return err
}

func (dao *SmsSceneDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE sms_scene SET scene_type=?,scene_desc=? WHERE id=?")
	return err
}

func (dao *SmsSceneDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM sms_scene WHERE id=?")
	return err
}

func (dao *SmsSceneDao) Insert(ctx context.Context, tx *wrap.Tx, e *SmsScene) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.SceneType, e.SceneDesc)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *SmsSceneDao) Update(ctx context.Context, tx *wrap.Tx, e *SmsScene) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.SceneType, e.SceneDesc, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *SmsSceneDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
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

func (dao *SmsSceneDao) scanRow(row *wrap.Row) (*SmsScene, error) {
	e := &SmsScene{}
	err := row.Scan(&e.Id, &e.SceneType, &e.SceneDesc, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *SmsSceneDao) scanRows(rows *wrap.Rows) (list []*SmsScene, err error) {
	list = make([]*SmsScene, 0)
	for rows.Next() {
		e := SmsScene{}
		err = rows.Scan(&e.Id, &e.SceneType, &e.SceneDesc, &e.CreateTime, &e.UpdateTime)
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

func (dao *SmsSceneDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*SmsScene, error) {
	querySql := "SELECT " + SMS_SCENE_ALL_FIELDS_STRING + " FROM sms_scene " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *SmsSceneDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*SmsScene, err error) {
	querySql := "SELECT " + SMS_SCENE_ALL_FIELDS_STRING + " FROM sms_scene " + query
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

func (dao *SmsSceneDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM sms_scene " + query
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

func (dao *SmsSceneDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM sms_scene " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *SmsSceneDao) GetQuery() *SmsSceneQuery {
	return NewSmsSceneQuery(dao)
}

type DB struct {
	wrap.DB
	Account          *AccountDao
	AccountOperation *AccountOperationDao
	SmsCode          *SmsCodeDao
	SmsScene         *SmsSceneDao
}

func NewDB() (d *DB, err error) {
	d = &DB{}

	connectionString := os.Getenv("DB")
	if connectionString == "" {
		return nil, fmt.Errorf("DB env nil")
	}
	connectionString += "/account?parseTime=true"
	db, err := wrap.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	d.DB = *db

	err = d.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	d.Account, err = NewAccountDao(d)
	if err != nil {
		return nil, err
	}

	d.AccountOperation, err = NewAccountOperationDao(d)
	if err != nil {
		return nil, err
	}

	d.SmsCode, err = NewSmsCodeDao(d)
	if err != nil {
		return nil, err
	}

	d.SmsScene, err = NewSmsSceneDao(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
