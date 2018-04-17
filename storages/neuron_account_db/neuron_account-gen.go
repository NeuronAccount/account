package neuron_account_db

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

func (q *BaseQuery) groupBy(fields ...string) {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = v
	}
}

func (q *BaseQuery) setLimit(startIncluded int64, count int64) {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
}

func (q *BaseQuery) orderBy(fieldName string, asc bool) {
	if q.order != "" {
		q.order += ","
	}
	q.order += fieldName + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}
}

func (q *BaseQuery) orderByGroupCount(asc bool) {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}
}

func (q *BaseQuery) setWhere(format string, a ...interface{}) {
	q.where += fmt.Sprintf(format, a...)
}

const ACCESS_TOKEN_TABLE_NAME = "access_token"

type ACCESS_TOKEN_FIELD string

const ACCESS_TOKEN_FIELD_ID = ACCESS_TOKEN_FIELD("id")
const ACCESS_TOKEN_FIELD_USER_ID = ACCESS_TOKEN_FIELD("user_id")
const ACCESS_TOKEN_FIELD_ACCESS_TOKEN = ACCESS_TOKEN_FIELD("access_token")
const ACCESS_TOKEN_FIELD_CREATE_TIME = ACCESS_TOKEN_FIELD("create_time")
const ACCESS_TOKEN_FIELD_UPDATE_TIME = ACCESS_TOKEN_FIELD("update_time")

const ACCESS_TOKEN_ALL_FIELDS_STRING = "id,user_id,access_token,create_time,update_time"

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
	q.setLimit(startIncluded, count)
	return q
}

func (q *AccessTokenQuery) OrderBy(fieldName ACCESS_TOKEN_FIELD, asc bool) *AccessTokenQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *AccessTokenQuery) OrderByGroupCount(asc bool) *AccessTokenQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *AccessTokenQuery) w(format string, a ...interface{}) *AccessTokenQuery {
	q.setWhere(format, a...)
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
func (q *AccessTokenQuery) AccessToken_Equal(v string) *AccessTokenQuery {
	return q.w("access_token='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccessToken_NotEqual(v string) *AccessTokenQuery {
	return q.w("access_token<>'" + fmt.Sprint(v) + "'")
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

type AccessTokenUpdate struct {
	dao    *AccessTokenDao
	keys   []string
	values []interface{}
}

func NewAccessTokenUpdate(dao *AccessTokenDao) *AccessTokenUpdate {
	q := &AccessTokenUpdate{}
	q.dao = dao
	q.keys = make([]string, 0)
	q.values = make([]interface{}, 0)

	return q
}

func (u *AccessTokenUpdate) Update(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	if len(u.keys) == 0 {
		err = fmt.Errorf("AccessTokenUpdate没有设置更新字段")
		u.dao.logger.Error("AccessTokenUpdate", zap.Error(err))
		return err
	}
	s := "UPDATE access_token SET " + strings.Join(u.keys, ",") + " WHERE id=?"
	v := append(u.values, id)
	if tx == nil {
		_, err = u.dao.db.Exec(ctx, s, v)
	} else {
		_, err = tx.Exec(ctx, s, v)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *AccessTokenUpdate) UserId(v string) *AccessTokenUpdate {
	u.keys = append(u.keys, "user_id=?")
	u.values = append(u.values, v)
	return u
}

func (u *AccessTokenUpdate) AccessToken(v string) *AccessTokenUpdate {
	u.keys = append(u.keys, "access_token=?")
	u.values = append(u.values, v)
	return u
}

type AccessTokenDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
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

func (dao *AccessTokenDao) GetUpdate() *AccessTokenUpdate {
	return NewAccessTokenUpdate(dao)
}

const ACCOUNT_OPERATION_TABLE_NAME = "account_operation"

type ACCOUNT_OPERATION_FIELD string

const ACCOUNT_OPERATION_FIELD_ID = ACCOUNT_OPERATION_FIELD("id")
const ACCOUNT_OPERATION_FIELD_USER_ID = ACCOUNT_OPERATION_FIELD("user_id")
const ACCOUNT_OPERATION_FIELD_OPERATIONTYPE = ACCOUNT_OPERATION_FIELD("operationType")
const ACCOUNT_OPERATION_FIELD_USER_AGENT = ACCOUNT_OPERATION_FIELD("user_agent")
const ACCOUNT_OPERATION_FIELD_PHONE_ENCRYPTED = ACCOUNT_OPERATION_FIELD("phone_encrypted")
const ACCOUNT_OPERATION_FIELD_SMS_SCENE = ACCOUNT_OPERATION_FIELD("sms_scene")
const ACCOUNT_OPERATION_FIELD_OTHER_USER_ID = ACCOUNT_OPERATION_FIELD("other_user_id")
const ACCOUNT_OPERATION_FIELD_CREATE_TIME = ACCOUNT_OPERATION_FIELD("create_time")

const ACCOUNT_OPERATION_ALL_FIELDS_STRING = "id,user_id,operationType,user_agent,phone_encrypted,sms_scene,other_user_id,create_time"

type AccountOperation struct {
	Id             uint64 //size=20
	UserId         string //size=32
	OperationType  string //size=32
	UserAgent      string //size=256
	PhoneEncrypted string //size=32
	SmsScene       string //size=32
	OtherUserId    string //size=32
	CreateTime     time.Time
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
	q.setLimit(startIncluded, count)
	return q
}

func (q *AccountOperationQuery) OrderBy(fieldName ACCOUNT_OPERATION_FIELD, asc bool) *AccountOperationQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *AccountOperationQuery) OrderByGroupCount(asc bool) *AccountOperationQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *AccountOperationQuery) w(format string, a ...interface{}) *AccountOperationQuery {
	q.setWhere(format, a...)
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
func (q *AccountOperationQuery) OperationType_Equal(v string) *AccountOperationQuery {
	return q.w("operationType='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OperationType_NotEqual(v string) *AccountOperationQuery {
	return q.w("operationType<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_Equal(v string) *AccountOperationQuery {
	return q.w("user_agent='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) UserAgent_NotEqual(v string) *AccountOperationQuery {
	return q.w("user_agent<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneEncrypted_Equal(v string) *AccountOperationQuery {
	return q.w("phone_encrypted='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) PhoneEncrypted_NotEqual(v string) *AccountOperationQuery {
	return q.w("phone_encrypted<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_Equal(v string) *AccountOperationQuery {
	return q.w("sms_scene='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) SmsScene_NotEqual(v string) *AccountOperationQuery {
	return q.w("sms_scene<>'" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OtherUserId_Equal(v string) *AccountOperationQuery {
	return q.w("other_user_id='" + fmt.Sprint(v) + "'")
}
func (q *AccountOperationQuery) OtherUserId_NotEqual(v string) *AccountOperationQuery {
	return q.w("other_user_id<>'" + fmt.Sprint(v) + "'")
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

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccountOperationDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO account_operation (user_id,operationType,user_agent,phone_encrypted,sms_scene,other_user_id) VALUES (?,?,?,?,?,?)")
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

	result, err := stmt.Exec(ctx, e.UserId, e.OperationType, e.UserAgent, e.PhoneEncrypted, e.SmsScene, e.OtherUserId)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
	err := row.Scan(&e.Id, &e.UserId, &e.OperationType, &e.UserAgent, &e.PhoneEncrypted, &e.SmsScene, &e.OtherUserId, &e.CreateTime)
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
		err = rows.Scan(&e.Id, &e.UserId, &e.OperationType, &e.UserAgent, &e.PhoneEncrypted, &e.SmsScene, &e.OtherUserId, &e.CreateTime)
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

const OAUTH_ACCOUNT_TABLE_NAME = "oauth_account"

type OAUTH_ACCOUNT_FIELD string

const OAUTH_ACCOUNT_FIELD_ID = OAUTH_ACCOUNT_FIELD("id")
const OAUTH_ACCOUNT_FIELD_USER_ID = OAUTH_ACCOUNT_FIELD("user_id")
const OAUTH_ACCOUNT_FIELD_PROVIDERID = OAUTH_ACCOUNT_FIELD("providerId")
const OAUTH_ACCOUNT_FIELD_PROVIDER_NAME = OAUTH_ACCOUNT_FIELD("provider_name")
const OAUTH_ACCOUNT_FIELD_OPEN_ID = OAUTH_ACCOUNT_FIELD("open_id")
const OAUTH_ACCOUNT_FIELD_USER_NAME = OAUTH_ACCOUNT_FIELD("user_name")
const OAUTH_ACCOUNT_FIELD_USER_ICON = OAUTH_ACCOUNT_FIELD("user_icon")
const OAUTH_ACCOUNT_FIELD_CREATE_TIME = OAUTH_ACCOUNT_FIELD("create_time")
const OAUTH_ACCOUNT_FIELD_UPDATE_TIME = OAUTH_ACCOUNT_FIELD("update_time")

const OAUTH_ACCOUNT_ALL_FIELDS_STRING = "id,user_id,providerId,provider_name,open_id,user_name,user_icon,create_time,update_time"

type OauthAccount struct {
	Id           uint64 //size=20
	UserId       string //size=32
	ProviderId   string //size=32
	ProviderName string //size=32
	OpenId       string //size=128
	UserName     string //size=32
	UserIcon     string //size=256
	CreateTime   time.Time
	UpdateTime   time.Time
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
	q.setLimit(startIncluded, count)
	return q
}

func (q *OauthAccountQuery) OrderBy(fieldName OAUTH_ACCOUNT_FIELD, asc bool) *OauthAccountQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *OauthAccountQuery) OrderByGroupCount(asc bool) *OauthAccountQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *OauthAccountQuery) w(format string, a ...interface{}) *OauthAccountQuery {
	q.setWhere(format, a...)
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
func (q *OauthAccountQuery) ProviderId_Equal(v string) *OauthAccountQuery {
	return q.w("providerId='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) ProviderId_NotEqual(v string) *OauthAccountQuery {
	return q.w("providerId<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) ProviderName_Equal(v string) *OauthAccountQuery {
	return q.w("provider_name='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) ProviderName_NotEqual(v string) *OauthAccountQuery {
	return q.w("provider_name<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OpenId_Equal(v string) *OauthAccountQuery {
	return q.w("open_id='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) OpenId_NotEqual(v string) *OauthAccountQuery {
	return q.w("open_id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserName_Equal(v string) *OauthAccountQuery {
	return q.w("user_name='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserName_NotEqual(v string) *OauthAccountQuery {
	return q.w("user_name<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserIcon_Equal(v string) *OauthAccountQuery {
	return q.w("user_icon='" + fmt.Sprint(v) + "'")
}
func (q *OauthAccountQuery) UserIcon_NotEqual(v string) *OauthAccountQuery {
	return q.w("user_icon<>'" + fmt.Sprint(v) + "'")
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

type OauthAccountUpdate struct {
	dao    *OauthAccountDao
	keys   []string
	values []interface{}
}

func NewOauthAccountUpdate(dao *OauthAccountDao) *OauthAccountUpdate {
	q := &OauthAccountUpdate{}
	q.dao = dao
	q.keys = make([]string, 0)
	q.values = make([]interface{}, 0)

	return q
}

func (u *OauthAccountUpdate) Update(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	if len(u.keys) == 0 {
		err = fmt.Errorf("OauthAccountUpdate没有设置更新字段")
		u.dao.logger.Error("OauthAccountUpdate", zap.Error(err))
		return err
	}
	s := "UPDATE oauth_account SET " + strings.Join(u.keys, ",") + " WHERE id=?"
	v := append(u.values, id)
	if tx == nil {
		_, err = u.dao.db.Exec(ctx, s, v)
	} else {
		_, err = tx.Exec(ctx, s, v)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *OauthAccountUpdate) UserId(v string) *OauthAccountUpdate {
	u.keys = append(u.keys, "user_id=?")
	u.values = append(u.values, v)
	return u
}

func (u *OauthAccountUpdate) ProviderId(v string) *OauthAccountUpdate {
	u.keys = append(u.keys, "providerId=?")
	u.values = append(u.values, v)
	return u
}

func (u *OauthAccountUpdate) ProviderName(v string) *OauthAccountUpdate {
	u.keys = append(u.keys, "provider_name=?")
	u.values = append(u.values, v)
	return u
}

func (u *OauthAccountUpdate) OpenId(v string) *OauthAccountUpdate {
	u.keys = append(u.keys, "open_id=?")
	u.values = append(u.values, v)
	return u
}

func (u *OauthAccountUpdate) UserName(v string) *OauthAccountUpdate {
	u.keys = append(u.keys, "user_name=?")
	u.values = append(u.values, v)
	return u
}

func (u *OauthAccountUpdate) UserIcon(v string) *OauthAccountUpdate {
	u.keys = append(u.keys, "user_icon=?")
	u.values = append(u.values, v)
	return u
}

type OauthAccountDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
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

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthAccountDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO oauth_account (user_id,providerId,provider_name,open_id,user_name,user_icon) VALUES (?,?,?,?,?,?)")
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

	result, err := stmt.Exec(ctx, e.UserId, e.ProviderId, e.ProviderName, e.OpenId, e.UserName, e.UserIcon)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
	err := row.Scan(&e.Id, &e.UserId, &e.ProviderId, &e.ProviderName, &e.OpenId, &e.UserName, &e.UserIcon, &e.CreateTime, &e.UpdateTime)
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
		err = rows.Scan(&e.Id, &e.UserId, &e.ProviderId, &e.ProviderName, &e.OpenId, &e.UserName, &e.UserIcon, &e.CreateTime, &e.UpdateTime)
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

func (dao *OauthAccountDao) GetUpdate() *OauthAccountUpdate {
	return NewOauthAccountUpdate(dao)
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
	q.setLimit(startIncluded, count)
	return q
}

func (q *OauthStateQuery) OrderBy(fieldName OAUTH_STATE_FIELD, asc bool) *OauthStateQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *OauthStateQuery) OrderByGroupCount(asc bool) *OauthStateQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *OauthStateQuery) w(format string, a ...interface{}) *OauthStateQuery {
	q.setWhere(format, a...)
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

type OauthStateUpdate struct {
	dao    *OauthStateDao
	keys   []string
	values []interface{}
}

func NewOauthStateUpdate(dao *OauthStateDao) *OauthStateUpdate {
	q := &OauthStateUpdate{}
	q.dao = dao
	q.keys = make([]string, 0)
	q.values = make([]interface{}, 0)

	return q
}

func (u *OauthStateUpdate) Update(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	if len(u.keys) == 0 {
		err = fmt.Errorf("OauthStateUpdate没有设置更新字段")
		u.dao.logger.Error("OauthStateUpdate", zap.Error(err))
		return err
	}
	s := "UPDATE oauth_state SET " + strings.Join(u.keys, ",") + " WHERE id=?"
	v := append(u.values, id)
	if tx == nil {
		_, err = u.dao.db.Exec(ctx, s, v)
	} else {
		_, err = tx.Exec(ctx, s, v)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *OauthStateUpdate) OauthState(v string) *OauthStateUpdate {
	u.keys = append(u.keys, "oauth_state=?")
	u.values = append(u.values, v)
	return u
}

func (u *OauthStateUpdate) IsUsed(v int32) *OauthStateUpdate {
	u.keys = append(u.keys, "is_used=?")
	u.values = append(u.values, v)
	return u
}

func (u *OauthStateUpdate) UserAgent(v string) *OauthStateUpdate {
	u.keys = append(u.keys, "user_agent=?")
	u.values = append(u.values, v)
	return u
}

type OauthStateDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
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

func (dao *OauthStateDao) GetUpdate() *OauthStateUpdate {
	return NewOauthStateUpdate(dao)
}

const PHONE_ACCOUNT_TABLE_NAME = "phone_account"

type PHONE_ACCOUNT_FIELD string

const PHONE_ACCOUNT_FIELD_ID = PHONE_ACCOUNT_FIELD("id")
const PHONE_ACCOUNT_FIELD_USER_ID = PHONE_ACCOUNT_FIELD("user_id")
const PHONE_ACCOUNT_FIELD_PHONE_ENCRYPTED = PHONE_ACCOUNT_FIELD("phone_encrypted")
const PHONE_ACCOUNT_FIELD_CREATE_TIME = PHONE_ACCOUNT_FIELD("create_time")
const PHONE_ACCOUNT_FIELD_UPDATE_TIME = PHONE_ACCOUNT_FIELD("update_time")

const PHONE_ACCOUNT_ALL_FIELDS_STRING = "id,user_id,phone_encrypted,create_time,update_time"

type PhoneAccount struct {
	Id             uint64 //size=20
	UserId         string //size=32
	PhoneEncrypted string //size=32
	CreateTime     time.Time
	UpdateTime     time.Time
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
	q.setLimit(startIncluded, count)
	return q
}

func (q *PhoneAccountQuery) OrderBy(fieldName PHONE_ACCOUNT_FIELD, asc bool) *PhoneAccountQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *PhoneAccountQuery) OrderByGroupCount(asc bool) *PhoneAccountQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *PhoneAccountQuery) w(format string, a ...interface{}) *PhoneAccountQuery {
	q.setWhere(format, a...)
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
func (q *PhoneAccountQuery) PhoneEncrypted_Equal(v string) *PhoneAccountQuery {
	return q.w("phone_encrypted='" + fmt.Sprint(v) + "'")
}
func (q *PhoneAccountQuery) PhoneEncrypted_NotEqual(v string) *PhoneAccountQuery {
	return q.w("phone_encrypted<>'" + fmt.Sprint(v) + "'")
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

type PhoneAccountUpdate struct {
	dao    *PhoneAccountDao
	keys   []string
	values []interface{}
}

func NewPhoneAccountUpdate(dao *PhoneAccountDao) *PhoneAccountUpdate {
	q := &PhoneAccountUpdate{}
	q.dao = dao
	q.keys = make([]string, 0)
	q.values = make([]interface{}, 0)

	return q
}

func (u *PhoneAccountUpdate) Update(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	if len(u.keys) == 0 {
		err = fmt.Errorf("PhoneAccountUpdate没有设置更新字段")
		u.dao.logger.Error("PhoneAccountUpdate", zap.Error(err))
		return err
	}
	s := "UPDATE phone_account SET " + strings.Join(u.keys, ",") + " WHERE id=?"
	v := append(u.values, id)
	if tx == nil {
		_, err = u.dao.db.Exec(ctx, s, v)
	} else {
		_, err = tx.Exec(ctx, s, v)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *PhoneAccountUpdate) UserId(v string) *PhoneAccountUpdate {
	u.keys = append(u.keys, "user_id=?")
	u.values = append(u.values, v)
	return u
}

func (u *PhoneAccountUpdate) PhoneEncrypted(v string) *PhoneAccountUpdate {
	u.keys = append(u.keys, "phone_encrypted=?")
	u.values = append(u.values, v)
	return u
}

type PhoneAccountDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
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

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *PhoneAccountDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO phone_account (user_id,phone_encrypted) VALUES (?,?)")
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

	result, err := stmt.Exec(ctx, e.UserId, e.PhoneEncrypted)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
	err := row.Scan(&e.Id, &e.UserId, &e.PhoneEncrypted, &e.CreateTime, &e.UpdateTime)
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
		err = rows.Scan(&e.Id, &e.UserId, &e.PhoneEncrypted, &e.CreateTime, &e.UpdateTime)
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

func (dao *PhoneAccountDao) GetUpdate() *PhoneAccountUpdate {
	return NewPhoneAccountUpdate(dao)
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
	q.setLimit(startIncluded, count)
	return q
}

func (q *RefreshTokenQuery) OrderBy(fieldName REFRESH_TOKEN_FIELD, asc bool) *RefreshTokenQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *RefreshTokenQuery) OrderByGroupCount(asc bool) *RefreshTokenQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *RefreshTokenQuery) w(format string, a ...interface{}) *RefreshTokenQuery {
	q.setWhere(format, a...)
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
func (q *RefreshTokenQuery) RefreshToken_Equal(v string) *RefreshTokenQuery {
	return q.w("refresh_token='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_NotEqual(v string) *RefreshTokenQuery {
	return q.w("refresh_token<>'" + fmt.Sprint(v) + "'")
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

type RefreshTokenUpdate struct {
	dao    *RefreshTokenDao
	keys   []string
	values []interface{}
}

func NewRefreshTokenUpdate(dao *RefreshTokenDao) *RefreshTokenUpdate {
	q := &RefreshTokenUpdate{}
	q.dao = dao
	q.keys = make([]string, 0)
	q.values = make([]interface{}, 0)

	return q
}

func (u *RefreshTokenUpdate) Update(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	if len(u.keys) == 0 {
		err = fmt.Errorf("RefreshTokenUpdate没有设置更新字段")
		u.dao.logger.Error("RefreshTokenUpdate", zap.Error(err))
		return err
	}
	s := "UPDATE refresh_token SET " + strings.Join(u.keys, ",") + " WHERE id=?"
	v := append(u.values, id)
	if tx == nil {
		_, err = u.dao.db.Exec(ctx, s, v)
	} else {
		_, err = tx.Exec(ctx, s, v)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *RefreshTokenUpdate) UserId(v string) *RefreshTokenUpdate {
	u.keys = append(u.keys, "user_id=?")
	u.values = append(u.values, v)
	return u
}

func (u *RefreshTokenUpdate) RefreshToken(v string) *RefreshTokenUpdate {
	u.keys = append(u.keys, "refresh_token=?")
	u.values = append(u.values, v)
	return u
}

func (u *RefreshTokenUpdate) IsLogout(v int32) *RefreshTokenUpdate {
	u.keys = append(u.keys, "is_logout=?")
	u.values = append(u.values, v)
	return u
}

func (u *RefreshTokenUpdate) LogoutTime(v time.Time) *RefreshTokenUpdate {
	u.keys = append(u.keys, "logout_time=?")
	u.values = append(u.values, v)
	return u
}

type RefreshTokenDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
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

func (dao *RefreshTokenDao) GetUpdate() *RefreshTokenUpdate {
	return NewRefreshTokenUpdate(dao)
}

const SMS_CODE_TABLE_NAME = "sms_code"

type SMS_CODE_FIELD string

const SMS_CODE_FIELD_ID = SMS_CODE_FIELD("id")
const SMS_CODE_FIELD_SMS_SCENE = SMS_CODE_FIELD("sms_scene")
const SMS_CODE_FIELD_PHONE_ENCRYPTED = SMS_CODE_FIELD("phone_encrypted")
const SMS_CODE_FIELD_SMS_CODE = SMS_CODE_FIELD("sms_code")
const SMS_CODE_FIELD_USER_ID = SMS_CODE_FIELD("user_id")
const SMS_CODE_FIELD_CREATE_TIME = SMS_CODE_FIELD("create_time")
const SMS_CODE_FIELD_UPDATE_TIME = SMS_CODE_FIELD("update_time")

const SMS_CODE_ALL_FIELDS_STRING = "id,sms_scene,phone_encrypted,sms_code,user_id,create_time,update_time"

type SmsCode struct {
	Id             uint64 //size=20
	SmsScene       string //size=32
	PhoneEncrypted string //size=32
	SmsCode        string //size=8
	UserId         string //size=32
	CreateTime     time.Time
	UpdateTime     time.Time
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
	q.setLimit(startIncluded, count)
	return q
}

func (q *SmsCodeQuery) OrderBy(fieldName SMS_CODE_FIELD, asc bool) *SmsCodeQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *SmsCodeQuery) OrderByGroupCount(asc bool) *SmsCodeQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *SmsCodeQuery) w(format string, a ...interface{}) *SmsCodeQuery {
	q.setWhere(format, a...)
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
func (q *SmsCodeQuery) SmsScene_Equal(v string) *SmsCodeQuery {
	return q.w("sms_scene='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsScene_NotEqual(v string) *SmsCodeQuery {
	return q.w("sms_scene<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneEncrypted_Equal(v string) *SmsCodeQuery {
	return q.w("phone_encrypted='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) PhoneEncrypted_NotEqual(v string) *SmsCodeQuery {
	return q.w("phone_encrypted<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_Equal(v string) *SmsCodeQuery {
	return q.w("sms_code='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) SmsCode_NotEqual(v string) *SmsCodeQuery {
	return q.w("sms_code<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) UserId_Equal(v string) *SmsCodeQuery {
	return q.w("user_id='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) UserId_NotEqual(v string) *SmsCodeQuery {
	return q.w("user_id<>'" + fmt.Sprint(v) + "'")
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
func (q *SmsCodeQuery) UpdateTime_Equal(v time.Time) *SmsCodeQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) UpdateTime_NotEqual(v time.Time) *SmsCodeQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) UpdateTime_Less(v time.Time) *SmsCodeQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) UpdateTime_LessEqual(v time.Time) *SmsCodeQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) UpdateTime_Greater(v time.Time) *SmsCodeQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *SmsCodeQuery) UpdateTime_GreaterEqual(v time.Time) *SmsCodeQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type SmsCodeUpdate struct {
	dao    *SmsCodeDao
	keys   []string
	values []interface{}
}

func NewSmsCodeUpdate(dao *SmsCodeDao) *SmsCodeUpdate {
	q := &SmsCodeUpdate{}
	q.dao = dao
	q.keys = make([]string, 0)
	q.values = make([]interface{}, 0)

	return q
}

func (u *SmsCodeUpdate) Update(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	if len(u.keys) == 0 {
		err = fmt.Errorf("SmsCodeUpdate没有设置更新字段")
		u.dao.logger.Error("SmsCodeUpdate", zap.Error(err))
		return err
	}
	s := "UPDATE sms_code SET " + strings.Join(u.keys, ",") + " WHERE id=?"
	v := append(u.values, id)
	if tx == nil {
		_, err = u.dao.db.Exec(ctx, s, v)
	} else {
		_, err = tx.Exec(ctx, s, v)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *SmsCodeUpdate) SmsScene(v string) *SmsCodeUpdate {
	u.keys = append(u.keys, "sms_scene=?")
	u.values = append(u.values, v)
	return u
}

func (u *SmsCodeUpdate) PhoneEncrypted(v string) *SmsCodeUpdate {
	u.keys = append(u.keys, "phone_encrypted=?")
	u.values = append(u.values, v)
	return u
}

func (u *SmsCodeUpdate) SmsCode(v string) *SmsCodeUpdate {
	u.keys = append(u.keys, "sms_code=?")
	u.values = append(u.values, v)
	return u
}

func (u *SmsCodeUpdate) UserId(v string) *SmsCodeUpdate {
	u.keys = append(u.keys, "user_id=?")
	u.values = append(u.values, v)
	return u
}

type SmsCodeDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
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

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *SmsCodeDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO sms_code (sms_scene,phone_encrypted,sms_code,user_id) VALUES (?,?,?,?)")
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

	result, err := stmt.Exec(ctx, e.SmsScene, e.PhoneEncrypted, e.SmsCode, e.UserId)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
	err := row.Scan(&e.Id, &e.SmsScene, &e.PhoneEncrypted, &e.SmsCode, &e.UserId, &e.CreateTime, &e.UpdateTime)
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
		err = rows.Scan(&e.Id, &e.SmsScene, &e.PhoneEncrypted, &e.SmsCode, &e.UserId, &e.CreateTime, &e.UpdateTime)
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

func (dao *SmsCodeDao) GetUpdate() *SmsCodeUpdate {
	return NewSmsCodeUpdate(dao)
}

const USER_INFO_TABLE_NAME = "user_info"

type USER_INFO_FIELD string

const USER_INFO_FIELD_ID = USER_INFO_FIELD("id")
const USER_INFO_FIELD_USER_ID = USER_INFO_FIELD("user_id")
const USER_INFO_FIELD_USER_NAME = USER_INFO_FIELD("user_name")
const USER_INFO_FIELD_USER_ICON = USER_INFO_FIELD("user_icon")
const USER_INFO_FIELD_CREATE_TIME = USER_INFO_FIELD("create_time")
const USER_INFO_FIELD_UPDATE_TIME = USER_INFO_FIELD("update_time")

const USER_INFO_ALL_FIELDS_STRING = "id,user_id,user_name,user_icon,create_time,update_time"

type UserInfo struct {
	Id         uint64 //size=20
	UserId     string //size=32
	UserName   string //size=32
	UserIcon   string //size=256
	CreateTime time.Time
	UpdateTime time.Time
}

type UserInfoQuery struct {
	BaseQuery
	dao *UserInfoDao
}

func NewUserInfoQuery(dao *UserInfoDao) *UserInfoQuery {
	q := &UserInfoQuery{}
	q.dao = dao

	return q
}

func (q *UserInfoQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*UserInfo, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *UserInfoQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*UserInfo, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *UserInfoQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *UserInfoQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *UserInfoQuery) ForUpdate() *UserInfoQuery {
	q.forUpdate = true
	return q
}

func (q *UserInfoQuery) ForShare() *UserInfoQuery {
	q.forShare = true
	return q
}

func (q *UserInfoQuery) GroupBy(fields ...USER_INFO_FIELD) *UserInfoQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *UserInfoQuery) Limit(startIncluded int64, count int64) *UserInfoQuery {
	q.setLimit(startIncluded, count)
	return q
}

func (q *UserInfoQuery) OrderBy(fieldName USER_INFO_FIELD, asc bool) *UserInfoQuery {
	q.orderBy(string(fieldName), asc)
	return q
}

func (q *UserInfoQuery) OrderByGroupCount(asc bool) *UserInfoQuery {
	q.orderByGroupCount(asc)
	return q
}

func (q *UserInfoQuery) w(format string, a ...interface{}) *UserInfoQuery {
	q.setWhere(format, a...)
	return q
}

func (q *UserInfoQuery) Left() *UserInfoQuery  { return q.w(" ( ") }
func (q *UserInfoQuery) Right() *UserInfoQuery { return q.w(" ) ") }
func (q *UserInfoQuery) And() *UserInfoQuery   { return q.w(" AND ") }
func (q *UserInfoQuery) Or() *UserInfoQuery    { return q.w(" OR ") }
func (q *UserInfoQuery) Not() *UserInfoQuery   { return q.w(" NOT ") }

func (q *UserInfoQuery) Id_Equal(v uint64) *UserInfoQuery     { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *UserInfoQuery) Id_NotEqual(v uint64) *UserInfoQuery  { return q.w("id<>'" + fmt.Sprint(v) + "'") }
func (q *UserInfoQuery) Id_Less(v uint64) *UserInfoQuery      { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *UserInfoQuery) Id_LessEqual(v uint64) *UserInfoQuery { return q.w("id<='" + fmt.Sprint(v) + "'") }
func (q *UserInfoQuery) Id_Greater(v uint64) *UserInfoQuery   { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *UserInfoQuery) Id_GreaterEqual(v uint64) *UserInfoQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UserId_Equal(v string) *UserInfoQuery {
	return q.w("user_id='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UserId_NotEqual(v string) *UserInfoQuery {
	return q.w("user_id<>'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UserName_Equal(v string) *UserInfoQuery {
	return q.w("user_name='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UserName_NotEqual(v string) *UserInfoQuery {
	return q.w("user_name<>'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UserIcon_Equal(v string) *UserInfoQuery {
	return q.w("user_icon='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UserIcon_NotEqual(v string) *UserInfoQuery {
	return q.w("user_icon<>'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) CreateTime_Equal(v time.Time) *UserInfoQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) CreateTime_NotEqual(v time.Time) *UserInfoQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) CreateTime_Less(v time.Time) *UserInfoQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) CreateTime_LessEqual(v time.Time) *UserInfoQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) CreateTime_Greater(v time.Time) *UserInfoQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) CreateTime_GreaterEqual(v time.Time) *UserInfoQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UpdateTime_Equal(v time.Time) *UserInfoQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UpdateTime_NotEqual(v time.Time) *UserInfoQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UpdateTime_Less(v time.Time) *UserInfoQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UpdateTime_LessEqual(v time.Time) *UserInfoQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UpdateTime_Greater(v time.Time) *UserInfoQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *UserInfoQuery) UpdateTime_GreaterEqual(v time.Time) *UserInfoQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type UserInfoUpdate struct {
	dao    *UserInfoDao
	keys   []string
	values []interface{}
}

func NewUserInfoUpdate(dao *UserInfoDao) *UserInfoUpdate {
	q := &UserInfoUpdate{}
	q.dao = dao
	q.keys = make([]string, 0)
	q.values = make([]interface{}, 0)

	return q
}

func (u *UserInfoUpdate) Update(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	if len(u.keys) == 0 {
		err = fmt.Errorf("UserInfoUpdate没有设置更新字段")
		u.dao.logger.Error("UserInfoUpdate", zap.Error(err))
		return err
	}
	s := "UPDATE user_info SET " + strings.Join(u.keys, ",") + " WHERE id=?"
	v := append(u.values, id)
	if tx == nil {
		_, err = u.dao.db.Exec(ctx, s, v)
	} else {
		_, err = tx.Exec(ctx, s, v)
	}

	if err != nil {
		return err
	}

	return nil
}

func (u *UserInfoUpdate) UserId(v string) *UserInfoUpdate {
	u.keys = append(u.keys, "user_id=?")
	u.values = append(u.values, v)
	return u
}

func (u *UserInfoUpdate) UserName(v string) *UserInfoUpdate {
	u.keys = append(u.keys, "user_name=?")
	u.values = append(u.values, v)
	return u
}

func (u *UserInfoUpdate) UserIcon(v string) *UserInfoUpdate {
	u.keys = append(u.keys, "user_icon=?")
	u.values = append(u.values, v)
	return u
}

type UserInfoDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewUserInfoDao(db *DB) (t *UserInfoDao, err error) {
	t = &UserInfoDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *UserInfoDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserInfoDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO user_info (user_id,user_name,user_icon) VALUES (?,?,?)")
	return err
}

func (dao *UserInfoDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM user_info WHERE id=?")
	return err
}

func (dao *UserInfoDao) Insert(ctx context.Context, tx *wrap.Tx, e *UserInfo) (id int64, err error) {
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

func (dao *UserInfoDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
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

func (dao *UserInfoDao) scanRow(row *wrap.Row) (*UserInfo, error) {
	e := &UserInfo{}
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

func (dao *UserInfoDao) scanRows(rows *wrap.Rows) (list []*UserInfo, err error) {
	list = make([]*UserInfo, 0)
	for rows.Next() {
		e := UserInfo{}
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

func (dao *UserInfoDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*UserInfo, error) {
	querySql := "SELECT " + USER_INFO_ALL_FIELDS_STRING + " FROM user_info " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *UserInfoDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*UserInfo, err error) {
	querySql := "SELECT " + USER_INFO_ALL_FIELDS_STRING + " FROM user_info " + query
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

func (dao *UserInfoDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM user_info " + query
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

func (dao *UserInfoDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM user_info " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *UserInfoDao) GetQuery() *UserInfoQuery {
	return NewUserInfoQuery(dao)
}

func (dao *UserInfoDao) GetUpdate() *UserInfoUpdate {
	return NewUserInfoUpdate(dao)
}

type DB struct {
	wrap.DB
	AccessToken      *AccessTokenDao
	AccountOperation *AccountOperationDao
	OauthAccount     *OauthAccountDao
	OauthState       *OauthStateDao
	PhoneAccount     *PhoneAccountDao
	RefreshToken     *RefreshTokenDao
	SmsCode          *SmsCodeDao
	UserInfo         *UserInfoDao
}

func NewDB() (d *DB, err error) {
	d = &DB{}

	connectionString := os.Getenv("DB")
	if connectionString == "" {
		return nil, fmt.Errorf("DB env nil")
	}
	connectionString += "/neuron_account?parseTime=true"
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

	d.SmsCode, err = NewSmsCodeDao(d)
	if err != nil {
		return nil, err
	}

	d.UserInfo, err = NewUserInfoDao(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
