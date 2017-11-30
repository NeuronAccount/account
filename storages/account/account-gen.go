package account

import (
	"bytes"
	"context"
	"fmt"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/sql/wrap"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"strings"
	"time"
)

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

const ACCOUNT_ALL_FIELDS_STRING = "id,account_id,phone_number,email_address,password_hash,oauth_provider,oauth_account_id,create_time,update_time"

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
}

type Account struct {
	Id             int64  //size=20
	AccountId      string //size=128
	PhoneNumber    string //size=32
	EmailAddress   string //size=128
	PasswordHash   string //size=128
	OauthProvider  string //size=128
	OauthAccountId string //size=128
	CreateTime     time.Time
	UpdateTime     time.Time
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

func (q *AccountQuery) Id_Equal(v int64) *AccountQuery        { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_NotEqual(v int64) *AccountQuery     { return q.w("id<>'" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_Less(v int64) *AccountQuery         { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_LessEqual(v int64) *AccountQuery    { return q.w("id<='" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_Greater(v int64) *AccountQuery      { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *AccountQuery) Id_GreaterEqual(v int64) *AccountQuery { return q.w("id>='" + fmt.Sprint(v) + "'") }
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
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO account (account_id,phone_number,email_address,password_hash,oauth_provider,oauth_account_id,create_time,update_time) VALUES (?,?,?,?,?,?,?,?)")
	return err
}

func (dao *AccountDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE account SET account_id=?,phone_number=?,email_address=?,password_hash=?,oauth_provider=?,oauth_account_id=?,create_time=?,update_time=? WHERE id=?")
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

	result, err := stmt.Exec(ctx, e.AccountId, e.PhoneNumber, e.EmailAddress, e.PasswordHash, e.OauthProvider, e.OauthAccountId, e.CreateTime, e.UpdateTime)
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

	_, err = stmt.Exec(ctx, e.AccountId, e.PhoneNumber, e.EmailAddress, e.PasswordHash, e.OauthProvider, e.OauthAccountId, e.CreateTime, e.UpdateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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
	err := row.Scan(&e.Id, &e.AccountId, &e.PhoneNumber, &e.EmailAddress, &e.PasswordHash, &e.OauthProvider, &e.OauthAccountId, &e.CreateTime, &e.UpdateTime)
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
		err = rows.Scan(&e.Id, &e.AccountId, &e.PhoneNumber, &e.EmailAddress, &e.PasswordHash, &e.OauthProvider, &e.OauthAccountId, &e.CreateTime, &e.UpdateTime)
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

const ACCOUNT_ID_GEN_TABLE_NAME = "account_id_gen"

type ACCOUNT_ID_GEN_FIELD string

const ACCOUNT_ID_GEN_FIELD_ID = ACCOUNT_ID_GEN_FIELD("id")
const ACCOUNT_ID_GEN_FIELD_MAX_ID = ACCOUNT_ID_GEN_FIELD("max_id")

const ACCOUNT_ID_GEN_ALL_FIELDS_STRING = "id,max_id"

var ACCOUNT_ID_GEN_ALL_FIELDS = []string{
	"id",
	"max_id",
}

type AccountIdGen struct {
	Id    int64 //size=20
	MaxId int64 //size=20
}

type AccountIdGenQuery struct {
	BaseQuery
	dao *AccountIdGenDao
}

func NewAccountIdGenQuery(dao *AccountIdGenDao) *AccountIdGenQuery {
	q := &AccountIdGenQuery{}
	q.dao = dao

	return q
}

func (q *AccountIdGenQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*AccountIdGen, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *AccountIdGenQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*AccountIdGen, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *AccountIdGenQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *AccountIdGenQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *AccountIdGenQuery) ForUpdate() *AccountIdGenQuery {
	q.forUpdate = true
	return q
}

func (q *AccountIdGenQuery) ForShare() *AccountIdGenQuery {
	q.forShare = true
	return q
}

func (q *AccountIdGenQuery) GroupBy(fields ...ACCOUNT_ID_GEN_FIELD) *AccountIdGenQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *AccountIdGenQuery) Limit(startIncluded int64, count int64) *AccountIdGenQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *AccountIdGenQuery) OrderBy(fieldName ACCOUNT_ID_GEN_FIELD, asc bool) *AccountIdGenQuery {
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

func (q *AccountIdGenQuery) OrderByGroupCount(asc bool) *AccountIdGenQuery {
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

func (q *AccountIdGenQuery) w(format string, a ...interface{}) *AccountIdGenQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *AccountIdGenQuery) Left() *AccountIdGenQuery  { return q.w(" ( ") }
func (q *AccountIdGenQuery) Right() *AccountIdGenQuery { return q.w(" ) ") }
func (q *AccountIdGenQuery) And() *AccountIdGenQuery   { return q.w(" AND ") }
func (q *AccountIdGenQuery) Or() *AccountIdGenQuery    { return q.w(" OR ") }
func (q *AccountIdGenQuery) Not() *AccountIdGenQuery   { return q.w(" NOT ") }

func (q *AccountIdGenQuery) Id_Equal(v int64) *AccountIdGenQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) Id_NotEqual(v int64) *AccountIdGenQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) Id_Less(v int64) *AccountIdGenQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) Id_LessEqual(v int64) *AccountIdGenQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) Id_Greater(v int64) *AccountIdGenQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) Id_GreaterEqual(v int64) *AccountIdGenQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) MaxId_Equal(v int64) *AccountIdGenQuery {
	return q.w("max_id='" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) MaxId_NotEqual(v int64) *AccountIdGenQuery {
	return q.w("max_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) MaxId_Less(v int64) *AccountIdGenQuery {
	return q.w("max_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) MaxId_LessEqual(v int64) *AccountIdGenQuery {
	return q.w("max_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) MaxId_Greater(v int64) *AccountIdGenQuery {
	return q.w("max_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountIdGenQuery) MaxId_GreaterEqual(v int64) *AccountIdGenQuery {
	return q.w("max_id>='" + fmt.Sprint(v) + "'")
}

type AccountIdGenDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewAccountIdGenDao(db *DB) (t *AccountIdGenDao, err error) {
	t = &AccountIdGenDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *AccountIdGenDao) init() (err error) {
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
func (dao *AccountIdGenDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO account_id_gen (max_id) VALUES (?)")
	return err
}

func (dao *AccountIdGenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE account_id_gen SET max_id=? WHERE id=?")
	return err
}

func (dao *AccountIdGenDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM account_id_gen WHERE id=?")
	return err
}

func (dao *AccountIdGenDao) Insert(ctx context.Context, tx *wrap.Tx, e *AccountIdGen) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.MaxId)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *AccountIdGenDao) Update(ctx context.Context, tx *wrap.Tx, e *AccountIdGen) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.MaxId, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountIdGenDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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

func (dao *AccountIdGenDao) scanRow(row *wrap.Row) (*AccountIdGen, error) {
	e := &AccountIdGen{}
	err := row.Scan(&e.Id, &e.MaxId)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *AccountIdGenDao) scanRows(rows *wrap.Rows) (list []*AccountIdGen, err error) {
	list = make([]*AccountIdGen, 0)
	for rows.Next() {
		e := AccountIdGen{}
		err = rows.Scan(&e.Id, &e.MaxId)
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

func (dao *AccountIdGenDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*AccountIdGen, error) {
	querySql := "SELECT " + ACCOUNT_ID_GEN_ALL_FIELDS_STRING + " FROM account_id_gen " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *AccountIdGenDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*AccountIdGen, err error) {
	querySql := "SELECT " + ACCOUNT_ID_GEN_ALL_FIELDS_STRING + " FROM account_id_gen " + query
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

func (dao *AccountIdGenDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM account_id_gen " + query
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

func (dao *AccountIdGenDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM account_id_gen " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *AccountIdGenDao) GetQuery() *AccountIdGenQuery {
	return NewAccountIdGenQuery(dao)
}

const ACCOUNT_OPERATION_TABLE_NAME = "account_operation"

type ACCOUNT_OPERATION_FIELD string

const ACCOUNT_OPERATION_FIELD_ID = ACCOUNT_OPERATION_FIELD("id")
const ACCOUNT_OPERATION_FIELD_CREATE_TIME = ACCOUNT_OPERATION_FIELD("create_time")

const ACCOUNT_OPERATION_ALL_FIELDS_STRING = "id,create_time"

var ACCOUNT_OPERATION_ALL_FIELDS = []string{
	"id",
	"create_time",
}

type AccountOperation struct {
	Id         int64 //size=20
	CreateTime time.Time
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

func (q *AccountOperationQuery) Id_Equal(v int64) *AccountOperationQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_NotEqual(v int64) *AccountOperationQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_Less(v int64) *AccountOperationQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_LessEqual(v int64) *AccountOperationQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_Greater(v int64) *AccountOperationQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) Id_GreaterEqual(v int64) *AccountOperationQuery {
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
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO account_operation (create_time) VALUES (?)")
	return err
}

func (dao *AccountOperationDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE account_operation SET create_time=? WHERE id=?")
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

	result, err := stmt.Exec(ctx, e.CreateTime)
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

	_, err = stmt.Exec(ctx, e.CreateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountOperationDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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
	err := row.Scan(&e.Id, &e.CreateTime)
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
		err = rows.Scan(&e.Id, &e.CreateTime)
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

const LOGIN_HISTORY_TABLE_NAME = "login_history"

type LOGIN_HISTORY_FIELD string

const LOGIN_HISTORY_FIELD_ID = LOGIN_HISTORY_FIELD("id")
const LOGIN_HISTORY_FIELD_ACCOUNT_ID = LOGIN_HISTORY_FIELD("account_id")
const LOGIN_HISTORY_FIELD_CREATE_TIME = LOGIN_HISTORY_FIELD("create_time")

const LOGIN_HISTORY_ALL_FIELDS_STRING = "id,account_id,create_time"

var LOGIN_HISTORY_ALL_FIELDS = []string{
	"id",
	"account_id",
	"create_time",
}

type LoginHistory struct {
	Id         int64  //size=20
	AccountId  string //size=128
	CreateTime time.Time
}

type LoginHistoryQuery struct {
	BaseQuery
	dao *LoginHistoryDao
}

func NewLoginHistoryQuery(dao *LoginHistoryDao) *LoginHistoryQuery {
	q := &LoginHistoryQuery{}
	q.dao = dao

	return q
}

func (q *LoginHistoryQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*LoginHistory, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *LoginHistoryQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*LoginHistory, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *LoginHistoryQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *LoginHistoryQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *LoginHistoryQuery) ForUpdate() *LoginHistoryQuery {
	q.forUpdate = true
	return q
}

func (q *LoginHistoryQuery) ForShare() *LoginHistoryQuery {
	q.forShare = true
	return q
}

func (q *LoginHistoryQuery) GroupBy(fields ...LOGIN_HISTORY_FIELD) *LoginHistoryQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *LoginHistoryQuery) Limit(startIncluded int64, count int64) *LoginHistoryQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *LoginHistoryQuery) OrderBy(fieldName LOGIN_HISTORY_FIELD, asc bool) *LoginHistoryQuery {
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

func (q *LoginHistoryQuery) OrderByGroupCount(asc bool) *LoginHistoryQuery {
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

func (q *LoginHistoryQuery) w(format string, a ...interface{}) *LoginHistoryQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *LoginHistoryQuery) Left() *LoginHistoryQuery  { return q.w(" ( ") }
func (q *LoginHistoryQuery) Right() *LoginHistoryQuery { return q.w(" ) ") }
func (q *LoginHistoryQuery) And() *LoginHistoryQuery   { return q.w(" AND ") }
func (q *LoginHistoryQuery) Or() *LoginHistoryQuery    { return q.w(" OR ") }
func (q *LoginHistoryQuery) Not() *LoginHistoryQuery   { return q.w(" NOT ") }

func (q *LoginHistoryQuery) Id_Equal(v int64) *LoginHistoryQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) Id_NotEqual(v int64) *LoginHistoryQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) Id_Less(v int64) *LoginHistoryQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) Id_LessEqual(v int64) *LoginHistoryQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) Id_Greater(v int64) *LoginHistoryQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) Id_GreaterEqual(v int64) *LoginHistoryQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) AccountId_Equal(v string) *LoginHistoryQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) AccountId_NotEqual(v string) *LoginHistoryQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) AccountId_Less(v string) *LoginHistoryQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) AccountId_LessEqual(v string) *LoginHistoryQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) AccountId_Greater(v string) *LoginHistoryQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) AccountId_GreaterEqual(v string) *LoginHistoryQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) CreateTime_Equal(v time.Time) *LoginHistoryQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) CreateTime_NotEqual(v time.Time) *LoginHistoryQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) CreateTime_Less(v time.Time) *LoginHistoryQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) CreateTime_LessEqual(v time.Time) *LoginHistoryQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) CreateTime_Greater(v time.Time) *LoginHistoryQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *LoginHistoryQuery) CreateTime_GreaterEqual(v time.Time) *LoginHistoryQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}

type LoginHistoryDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewLoginHistoryDao(db *DB) (t *LoginHistoryDao, err error) {
	t = &LoginHistoryDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *LoginHistoryDao) init() (err error) {
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
func (dao *LoginHistoryDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO login_history (account_id,create_time) VALUES (?,?)")
	return err
}

func (dao *LoginHistoryDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE login_history SET account_id=?,create_time=? WHERE id=?")
	return err
}

func (dao *LoginHistoryDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM login_history WHERE id=?")
	return err
}

func (dao *LoginHistoryDao) Insert(ctx context.Context, tx *wrap.Tx, e *LoginHistory) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.AccountId, e.CreateTime)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *LoginHistoryDao) Update(ctx context.Context, tx *wrap.Tx, e *LoginHistory) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.AccountId, e.CreateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *LoginHistoryDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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

func (dao *LoginHistoryDao) scanRow(row *wrap.Row) (*LoginHistory, error) {
	e := &LoginHistory{}
	err := row.Scan(&e.Id, &e.AccountId, &e.CreateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *LoginHistoryDao) scanRows(rows *wrap.Rows) (list []*LoginHistory, err error) {
	list = make([]*LoginHistory, 0)
	for rows.Next() {
		e := LoginHistory{}
		err = rows.Scan(&e.Id, &e.AccountId, &e.CreateTime)
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

func (dao *LoginHistoryDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*LoginHistory, error) {
	querySql := "SELECT " + LOGIN_HISTORY_ALL_FIELDS_STRING + " FROM login_history " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *LoginHistoryDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*LoginHistory, err error) {
	querySql := "SELECT " + LOGIN_HISTORY_ALL_FIELDS_STRING + " FROM login_history " + query
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

func (dao *LoginHistoryDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM login_history " + query
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

func (dao *LoginHistoryDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM login_history " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *LoginHistoryDao) GetQuery() *LoginHistoryQuery {
	return NewLoginHistoryQuery(dao)
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
	Id          int64  //size=20
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

func (q *SmsCodeQuery) Id_Equal(v int64) *SmsCodeQuery        { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_NotEqual(v int64) *SmsCodeQuery     { return q.w("id<>'" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_Less(v int64) *SmsCodeQuery         { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_LessEqual(v int64) *SmsCodeQuery    { return q.w("id<='" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_Greater(v int64) *SmsCodeQuery      { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *SmsCodeQuery) Id_GreaterEqual(v int64) *SmsCodeQuery { return q.w("id>='" + fmt.Sprint(v) + "'") }
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
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO sms_code (scene_type,phone_number,sms_code,create_time) VALUES (?,?,?,?)")
	return err
}

func (dao *SmsCodeDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE sms_code SET scene_type=?,phone_number=?,sms_code=?,create_time=? WHERE id=?")
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

	result, err := stmt.Exec(ctx, e.SceneType, e.PhoneNumber, e.SmsCode, e.CreateTime)
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

	_, err = stmt.Exec(ctx, e.SceneType, e.PhoneNumber, e.SmsCode, e.CreateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *SmsCodeDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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

type DB struct {
	wrap.DB
	Account          *AccountDao
	AccountIdGen     *AccountIdGenDao
	AccountOperation *AccountOperationDao
	LoginHistory     *LoginHistoryDao
	SmsCode          *SmsCodeDao
}

func NewDB(connectionString string) (d *DB, err error) {
	if connectionString == "" {
		return nil, fmt.Errorf("connectionString nil")
	}

	d = &DB{}

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

	d.AccountIdGen, err = NewAccountIdGenDao(d)
	if err != nil {
		return nil, err
	}

	d.AccountOperation, err = NewAccountOperationDao(d)
	if err != nil {
		return nil, err
	}

	d.LoginHistory, err = NewLoginHistoryDao(d)
	if err != nil {
		return nil, err
	}

	d.SmsCode, err = NewSmsCodeDao(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
