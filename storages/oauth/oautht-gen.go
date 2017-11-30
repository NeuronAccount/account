package oauth

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

const ACCESS_TOKEN_TABLE_NAME = "access_token"

type ACCESS_TOKEN_FIELD string

const ACCESS_TOKEN_FIELD_ID = ACCESS_TOKEN_FIELD("id")
const ACCESS_TOKEN_FIELD_ACCESS_TOKEN = ACCESS_TOKEN_FIELD("access_token")
const ACCESS_TOKEN_FIELD_CLIENT_ID = ACCESS_TOKEN_FIELD("client_id")
const ACCESS_TOKEN_FIELD_ACCOUNT_ID = ACCESS_TOKEN_FIELD("account_id")
const ACCESS_TOKEN_FIELD_EXPIRE_SECONDS = ACCESS_TOKEN_FIELD("expire_seconds")
const ACCESS_TOKEN_FIELD_OAUTH_SCOPE = ACCESS_TOKEN_FIELD("oauth_scope")
const ACCESS_TOKEN_FIELD_CREATE_TIME = ACCESS_TOKEN_FIELD("create_time")
const ACCESS_TOKEN_FIELD_UPDATE_TIME = ACCESS_TOKEN_FIELD("update_time")

const ACCESS_TOKEN_ALL_FIELDS_STRING = "id,access_token,client_id,account_id,expire_seconds,oauth_scope,create_time,update_time"

var ACCESS_TOKEN_ALL_FIELDS = []string{
	"id",
	"access_token",
	"client_id",
	"account_id",
	"expire_seconds",
	"oauth_scope",
	"create_time",
	"update_time",
}

type AccessToken struct {
	Id            int64  //size=20
	AccessToken   string //size=128
	ClientId      string //size=32
	AccountId     string //size=32
	ExpireSeconds int64  //size=20
	OauthScope    string //size=256
	CreateTime    time.Time
	UpdateTime    time.Time
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

func (q *AccessTokenQuery) Id_Equal(v int64) *AccessTokenQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_NotEqual(v int64) *AccessTokenQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_Less(v int64) *AccessTokenQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_LessEqual(v int64) *AccessTokenQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_Greater(v int64) *AccessTokenQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) Id_GreaterEqual(v int64) *AccessTokenQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
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
func (q *AccessTokenQuery) ClientId_Equal(v string) *AccessTokenQuery {
	return q.w("client_id='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ClientId_NotEqual(v string) *AccessTokenQuery {
	return q.w("client_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ClientId_Less(v string) *AccessTokenQuery {
	return q.w("client_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ClientId_LessEqual(v string) *AccessTokenQuery {
	return q.w("client_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ClientId_Greater(v string) *AccessTokenQuery {
	return q.w("client_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ClientId_GreaterEqual(v string) *AccessTokenQuery {
	return q.w("client_id>='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccountId_Equal(v string) *AccessTokenQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccountId_NotEqual(v string) *AccessTokenQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccountId_Less(v string) *AccessTokenQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccountId_LessEqual(v string) *AccessTokenQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccountId_Greater(v string) *AccessTokenQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) AccountId_GreaterEqual(v string) *AccessTokenQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ExpireSeconds_Equal(v int64) *AccessTokenQuery {
	return q.w("expire_seconds='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ExpireSeconds_NotEqual(v int64) *AccessTokenQuery {
	return q.w("expire_seconds<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ExpireSeconds_Less(v int64) *AccessTokenQuery {
	return q.w("expire_seconds<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ExpireSeconds_LessEqual(v int64) *AccessTokenQuery {
	return q.w("expire_seconds<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ExpireSeconds_Greater(v int64) *AccessTokenQuery {
	return q.w("expire_seconds>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) ExpireSeconds_GreaterEqual(v int64) *AccessTokenQuery {
	return q.w("expire_seconds>='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) OauthScope_Equal(v string) *AccessTokenQuery {
	return q.w("oauth_scope='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) OauthScope_NotEqual(v string) *AccessTokenQuery {
	return q.w("oauth_scope<>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) OauthScope_Less(v string) *AccessTokenQuery {
	return q.w("oauth_scope<'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) OauthScope_LessEqual(v string) *AccessTokenQuery {
	return q.w("oauth_scope<='" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) OauthScope_Greater(v string) *AccessTokenQuery {
	return q.w("oauth_scope>'" + fmt.Sprint(v) + "'")
}
func (q *AccessTokenQuery) OauthScope_GreaterEqual(v string) *AccessTokenQuery {
	return q.w("oauth_scope>='" + fmt.Sprint(v) + "'")
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
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO access_token (access_token,client_id,account_id,expire_seconds,oauth_scope,create_time,update_time) VALUES (?,?,?,?,?,?,?)")
	return err
}

func (dao *AccessTokenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE access_token SET access_token=?,client_id=?,account_id=?,expire_seconds=?,oauth_scope=?,create_time=?,update_time=? WHERE id=?")
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

	result, err := stmt.Exec(ctx, e.AccessToken, e.ClientId, e.AccountId, e.ExpireSeconds, e.OauthScope, e.CreateTime, e.UpdateTime)
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

	_, err = stmt.Exec(ctx, e.AccessToken, e.ClientId, e.AccountId, e.ExpireSeconds, e.OauthScope, e.CreateTime, e.UpdateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AccessTokenDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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
	err := row.Scan(&e.Id, &e.AccessToken, &e.ClientId, &e.AccountId, &e.ExpireSeconds, &e.OauthScope, &e.CreateTime, &e.UpdateTime)
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
		err = rows.Scan(&e.Id, &e.AccessToken, &e.ClientId, &e.AccountId, &e.ExpireSeconds, &e.OauthScope, &e.CreateTime, &e.UpdateTime)
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

const AUTHORIZATION_CODE_TABLE_NAME = "authorization_code"

type AUTHORIZATION_CODE_FIELD string

const AUTHORIZATION_CODE_FIELD_ID = AUTHORIZATION_CODE_FIELD("id")
const AUTHORIZATION_CODE_FIELD_AUTHORIZATION_CODE = AUTHORIZATION_CODE_FIELD("authorization_code")
const AUTHORIZATION_CODE_FIELD_CLIENT_ID = AUTHORIZATION_CODE_FIELD("client_id")
const AUTHORIZATION_CODE_FIELD_ACCOUNT_ID = AUTHORIZATION_CODE_FIELD("account_id")
const AUTHORIZATION_CODE_FIELD_REDIRECT_URI = AUTHORIZATION_CODE_FIELD("redirect_uri")
const AUTHORIZATION_CODE_FIELD_OAUTH_SCOPE = AUTHORIZATION_CODE_FIELD("oauth_scope")
const AUTHORIZATION_CODE_FIELD_EXPIRE_SECONDS = AUTHORIZATION_CODE_FIELD("expire_seconds")
const AUTHORIZATION_CODE_FIELD_CREATE_TIME = AUTHORIZATION_CODE_FIELD("create_time")
const AUTHORIZATION_CODE_FIELD_UPDATE_TIME = AUTHORIZATION_CODE_FIELD("update_time")

const AUTHORIZATION_CODE_ALL_FIELDS_STRING = "id,authorization_code,client_id,account_id,redirect_uri,oauth_scope,expire_seconds,create_time,update_time"

var AUTHORIZATION_CODE_ALL_FIELDS = []string{
	"id",
	"authorization_code",
	"client_id",
	"account_id",
	"redirect_uri",
	"oauth_scope",
	"expire_seconds",
	"create_time",
	"update_time",
}

type AuthorizationCode struct {
	Id                int64  //size=20
	AuthorizationCode string //size=128
	ClientId          string //size=32
	AccountId         string //size=32
	RedirectUri       string //size=256
	OauthScope        string //size=256
	ExpireSeconds     int64  //size=20
	CreateTime        time.Time
	UpdateTime        time.Time
}

type AuthorizationCodeQuery struct {
	BaseQuery
	dao *AuthorizationCodeDao
}

func NewAuthorizationCodeQuery(dao *AuthorizationCodeDao) *AuthorizationCodeQuery {
	q := &AuthorizationCodeQuery{}
	q.dao = dao

	return q
}

func (q *AuthorizationCodeQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*AuthorizationCode, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *AuthorizationCodeQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*AuthorizationCode, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *AuthorizationCodeQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *AuthorizationCodeQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *AuthorizationCodeQuery) ForUpdate() *AuthorizationCodeQuery {
	q.forUpdate = true
	return q
}

func (q *AuthorizationCodeQuery) ForShare() *AuthorizationCodeQuery {
	q.forShare = true
	return q
}

func (q *AuthorizationCodeQuery) GroupBy(fields ...AUTHORIZATION_CODE_FIELD) *AuthorizationCodeQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *AuthorizationCodeQuery) Limit(startIncluded int64, count int64) *AuthorizationCodeQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *AuthorizationCodeQuery) OrderBy(fieldName AUTHORIZATION_CODE_FIELD, asc bool) *AuthorizationCodeQuery {
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

func (q *AuthorizationCodeQuery) OrderByGroupCount(asc bool) *AuthorizationCodeQuery {
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

func (q *AuthorizationCodeQuery) w(format string, a ...interface{}) *AuthorizationCodeQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *AuthorizationCodeQuery) Left() *AuthorizationCodeQuery  { return q.w(" ( ") }
func (q *AuthorizationCodeQuery) Right() *AuthorizationCodeQuery { return q.w(" ) ") }
func (q *AuthorizationCodeQuery) And() *AuthorizationCodeQuery   { return q.w(" AND ") }
func (q *AuthorizationCodeQuery) Or() *AuthorizationCodeQuery    { return q.w(" OR ") }
func (q *AuthorizationCodeQuery) Not() *AuthorizationCodeQuery   { return q.w(" NOT ") }

func (q *AuthorizationCodeQuery) Id_Equal(v int64) *AuthorizationCodeQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) Id_NotEqual(v int64) *AuthorizationCodeQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) Id_Less(v int64) *AuthorizationCodeQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) Id_LessEqual(v int64) *AuthorizationCodeQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) Id_Greater(v int64) *AuthorizationCodeQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) Id_GreaterEqual(v int64) *AuthorizationCodeQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AuthorizationCode_Equal(v string) *AuthorizationCodeQuery {
	return q.w("authorization_code='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AuthorizationCode_NotEqual(v string) *AuthorizationCodeQuery {
	return q.w("authorization_code<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AuthorizationCode_Less(v string) *AuthorizationCodeQuery {
	return q.w("authorization_code<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AuthorizationCode_LessEqual(v string) *AuthorizationCodeQuery {
	return q.w("authorization_code<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AuthorizationCode_Greater(v string) *AuthorizationCodeQuery {
	return q.w("authorization_code>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AuthorizationCode_GreaterEqual(v string) *AuthorizationCodeQuery {
	return q.w("authorization_code>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ClientId_Equal(v string) *AuthorizationCodeQuery {
	return q.w("client_id='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ClientId_NotEqual(v string) *AuthorizationCodeQuery {
	return q.w("client_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ClientId_Less(v string) *AuthorizationCodeQuery {
	return q.w("client_id<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ClientId_LessEqual(v string) *AuthorizationCodeQuery {
	return q.w("client_id<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ClientId_Greater(v string) *AuthorizationCodeQuery {
	return q.w("client_id>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ClientId_GreaterEqual(v string) *AuthorizationCodeQuery {
	return q.w("client_id>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AccountId_Equal(v string) *AuthorizationCodeQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AccountId_NotEqual(v string) *AuthorizationCodeQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AccountId_Less(v string) *AuthorizationCodeQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AccountId_LessEqual(v string) *AuthorizationCodeQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AccountId_Greater(v string) *AuthorizationCodeQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) AccountId_GreaterEqual(v string) *AuthorizationCodeQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) RedirectUri_Equal(v string) *AuthorizationCodeQuery {
	return q.w("redirect_uri='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) RedirectUri_NotEqual(v string) *AuthorizationCodeQuery {
	return q.w("redirect_uri<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) RedirectUri_Less(v string) *AuthorizationCodeQuery {
	return q.w("redirect_uri<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) RedirectUri_LessEqual(v string) *AuthorizationCodeQuery {
	return q.w("redirect_uri<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) RedirectUri_Greater(v string) *AuthorizationCodeQuery {
	return q.w("redirect_uri>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) RedirectUri_GreaterEqual(v string) *AuthorizationCodeQuery {
	return q.w("redirect_uri>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) OauthScope_Equal(v string) *AuthorizationCodeQuery {
	return q.w("oauth_scope='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) OauthScope_NotEqual(v string) *AuthorizationCodeQuery {
	return q.w("oauth_scope<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) OauthScope_Less(v string) *AuthorizationCodeQuery {
	return q.w("oauth_scope<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) OauthScope_LessEqual(v string) *AuthorizationCodeQuery {
	return q.w("oauth_scope<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) OauthScope_Greater(v string) *AuthorizationCodeQuery {
	return q.w("oauth_scope>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) OauthScope_GreaterEqual(v string) *AuthorizationCodeQuery {
	return q.w("oauth_scope>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ExpireSeconds_Equal(v int64) *AuthorizationCodeQuery {
	return q.w("expire_seconds='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ExpireSeconds_NotEqual(v int64) *AuthorizationCodeQuery {
	return q.w("expire_seconds<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ExpireSeconds_Less(v int64) *AuthorizationCodeQuery {
	return q.w("expire_seconds<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ExpireSeconds_LessEqual(v int64) *AuthorizationCodeQuery {
	return q.w("expire_seconds<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ExpireSeconds_Greater(v int64) *AuthorizationCodeQuery {
	return q.w("expire_seconds>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) ExpireSeconds_GreaterEqual(v int64) *AuthorizationCodeQuery {
	return q.w("expire_seconds>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) CreateTime_Equal(v time.Time) *AuthorizationCodeQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) CreateTime_NotEqual(v time.Time) *AuthorizationCodeQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) CreateTime_Less(v time.Time) *AuthorizationCodeQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) CreateTime_LessEqual(v time.Time) *AuthorizationCodeQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) CreateTime_Greater(v time.Time) *AuthorizationCodeQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) CreateTime_GreaterEqual(v time.Time) *AuthorizationCodeQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) UpdateTime_Equal(v time.Time) *AuthorizationCodeQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) UpdateTime_NotEqual(v time.Time) *AuthorizationCodeQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) UpdateTime_Less(v time.Time) *AuthorizationCodeQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) UpdateTime_LessEqual(v time.Time) *AuthorizationCodeQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) UpdateTime_Greater(v time.Time) *AuthorizationCodeQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *AuthorizationCodeQuery) UpdateTime_GreaterEqual(v time.Time) *AuthorizationCodeQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type AuthorizationCodeDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewAuthorizationCodeDao(db *DB) (t *AuthorizationCodeDao, err error) {
	t = &AuthorizationCodeDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *AuthorizationCodeDao) init() (err error) {
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
func (dao *AuthorizationCodeDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO authorization_code (authorization_code,client_id,account_id,redirect_uri,oauth_scope,expire_seconds,create_time,update_time) VALUES (?,?,?,?,?,?,?,?)")
	return err
}

func (dao *AuthorizationCodeDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE authorization_code SET authorization_code=?,client_id=?,account_id=?,redirect_uri=?,oauth_scope=?,expire_seconds=?,create_time=?,update_time=? WHERE id=?")
	return err
}

func (dao *AuthorizationCodeDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM authorization_code WHERE id=?")
	return err
}

func (dao *AuthorizationCodeDao) Insert(ctx context.Context, tx *wrap.Tx, e *AuthorizationCode) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.AuthorizationCode, e.ClientId, e.AccountId, e.RedirectUri, e.OauthScope, e.ExpireSeconds, e.CreateTime, e.UpdateTime)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *AuthorizationCodeDao) Update(ctx context.Context, tx *wrap.Tx, e *AuthorizationCode) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.AuthorizationCode, e.ClientId, e.AccountId, e.RedirectUri, e.OauthScope, e.ExpireSeconds, e.CreateTime, e.UpdateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *AuthorizationCodeDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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

func (dao *AuthorizationCodeDao) scanRow(row *wrap.Row) (*AuthorizationCode, error) {
	e := &AuthorizationCode{}
	err := row.Scan(&e.Id, &e.AuthorizationCode, &e.ClientId, &e.AccountId, &e.RedirectUri, &e.OauthScope, &e.ExpireSeconds, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *AuthorizationCodeDao) scanRows(rows *wrap.Rows) (list []*AuthorizationCode, err error) {
	list = make([]*AuthorizationCode, 0)
	for rows.Next() {
		e := AuthorizationCode{}
		err = rows.Scan(&e.Id, &e.AuthorizationCode, &e.ClientId, &e.AccountId, &e.RedirectUri, &e.OauthScope, &e.ExpireSeconds, &e.CreateTime, &e.UpdateTime)
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

func (dao *AuthorizationCodeDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*AuthorizationCode, error) {
	querySql := "SELECT " + AUTHORIZATION_CODE_ALL_FIELDS_STRING + " FROM authorization_code " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *AuthorizationCodeDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*AuthorizationCode, err error) {
	querySql := "SELECT " + AUTHORIZATION_CODE_ALL_FIELDS_STRING + " FROM authorization_code " + query
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

func (dao *AuthorizationCodeDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM authorization_code " + query
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

func (dao *AuthorizationCodeDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM authorization_code " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *AuthorizationCodeDao) GetQuery() *AuthorizationCodeQuery {
	return NewAuthorizationCodeQuery(dao)
}

const OAUTH_CLIENT_TABLE_NAME = "oauth_client"

type OAUTH_CLIENT_FIELD string

const OAUTH_CLIENT_FIELD_ID = OAUTH_CLIENT_FIELD("id")
const OAUTH_CLIENT_FIELD_CLIENT_ID = OAUTH_CLIENT_FIELD("client_id")
const OAUTH_CLIENT_FIELD_ACCOUNT_ID = OAUTH_CLIENT_FIELD("account_id")
const OAUTH_CLIENT_FIELD_PASSWORD_HASH = OAUTH_CLIENT_FIELD("password_hash")
const OAUTH_CLIENT_FIELD_REDIRECT_URI = OAUTH_CLIENT_FIELD("redirect_uri")
const OAUTH_CLIENT_FIELD_CREATE_TIME = OAUTH_CLIENT_FIELD("create_time")
const OAUTH_CLIENT_FIELD_UPDATE_TIME = OAUTH_CLIENT_FIELD("update_time")

const OAUTH_CLIENT_ALL_FIELDS_STRING = "id,client_id,account_id,password_hash,redirect_uri,create_time,update_time"

var OAUTH_CLIENT_ALL_FIELDS = []string{
	"id",
	"client_id",
	"account_id",
	"password_hash",
	"redirect_uri",
	"create_time",
	"update_time",
}

type OauthClient struct {
	Id           int64  //size=20
	ClientId     string //size=32
	AccountId    string //size=32
	PasswordHash string //size=128
	RedirectUri  string //size=256
	CreateTime   time.Time
	UpdateTime   time.Time
}

type OauthClientQuery struct {
	BaseQuery
	dao *OauthClientDao
}

func NewOauthClientQuery(dao *OauthClientDao) *OauthClientQuery {
	q := &OauthClientQuery{}
	q.dao = dao

	return q
}

func (q *OauthClientQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*OauthClient, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *OauthClientQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*OauthClient, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *OauthClientQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *OauthClientQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *OauthClientQuery) ForUpdate() *OauthClientQuery {
	q.forUpdate = true
	return q
}

func (q *OauthClientQuery) ForShare() *OauthClientQuery {
	q.forShare = true
	return q
}

func (q *OauthClientQuery) GroupBy(fields ...OAUTH_CLIENT_FIELD) *OauthClientQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *OauthClientQuery) Limit(startIncluded int64, count int64) *OauthClientQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *OauthClientQuery) OrderBy(fieldName OAUTH_CLIENT_FIELD, asc bool) *OauthClientQuery {
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

func (q *OauthClientQuery) OrderByGroupCount(asc bool) *OauthClientQuery {
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

func (q *OauthClientQuery) w(format string, a ...interface{}) *OauthClientQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *OauthClientQuery) Left() *OauthClientQuery  { return q.w(" ( ") }
func (q *OauthClientQuery) Right() *OauthClientQuery { return q.w(" ) ") }
func (q *OauthClientQuery) And() *OauthClientQuery   { return q.w(" AND ") }
func (q *OauthClientQuery) Or() *OauthClientQuery    { return q.w(" OR ") }
func (q *OauthClientQuery) Not() *OauthClientQuery   { return q.w(" NOT ") }

func (q *OauthClientQuery) Id_Equal(v int64) *OauthClientQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) Id_NotEqual(v int64) *OauthClientQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) Id_Less(v int64) *OauthClientQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) Id_LessEqual(v int64) *OauthClientQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) Id_Greater(v int64) *OauthClientQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) Id_GreaterEqual(v int64) *OauthClientQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) ClientId_Equal(v string) *OauthClientQuery {
	return q.w("client_id='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) ClientId_NotEqual(v string) *OauthClientQuery {
	return q.w("client_id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) ClientId_Less(v string) *OauthClientQuery {
	return q.w("client_id<'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) ClientId_LessEqual(v string) *OauthClientQuery {
	return q.w("client_id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) ClientId_Greater(v string) *OauthClientQuery {
	return q.w("client_id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) ClientId_GreaterEqual(v string) *OauthClientQuery {
	return q.w("client_id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) AccountId_Equal(v string) *OauthClientQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) AccountId_NotEqual(v string) *OauthClientQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) AccountId_Less(v string) *OauthClientQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) AccountId_LessEqual(v string) *OauthClientQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) AccountId_Greater(v string) *OauthClientQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) AccountId_GreaterEqual(v string) *OauthClientQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) PasswordHash_Equal(v string) *OauthClientQuery {
	return q.w("password_hash='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) PasswordHash_NotEqual(v string) *OauthClientQuery {
	return q.w("password_hash<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) PasswordHash_Less(v string) *OauthClientQuery {
	return q.w("password_hash<'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) PasswordHash_LessEqual(v string) *OauthClientQuery {
	return q.w("password_hash<='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) PasswordHash_Greater(v string) *OauthClientQuery {
	return q.w("password_hash>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) PasswordHash_GreaterEqual(v string) *OauthClientQuery {
	return q.w("password_hash>='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) RedirectUri_Equal(v string) *OauthClientQuery {
	return q.w("redirect_uri='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) RedirectUri_NotEqual(v string) *OauthClientQuery {
	return q.w("redirect_uri<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) RedirectUri_Less(v string) *OauthClientQuery {
	return q.w("redirect_uri<'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) RedirectUri_LessEqual(v string) *OauthClientQuery {
	return q.w("redirect_uri<='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) RedirectUri_Greater(v string) *OauthClientQuery {
	return q.w("redirect_uri>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) RedirectUri_GreaterEqual(v string) *OauthClientQuery {
	return q.w("redirect_uri>='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) CreateTime_Equal(v time.Time) *OauthClientQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) CreateTime_NotEqual(v time.Time) *OauthClientQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) CreateTime_Less(v time.Time) *OauthClientQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) CreateTime_LessEqual(v time.Time) *OauthClientQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) CreateTime_Greater(v time.Time) *OauthClientQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) CreateTime_GreaterEqual(v time.Time) *OauthClientQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) UpdateTime_Equal(v time.Time) *OauthClientQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) UpdateTime_NotEqual(v time.Time) *OauthClientQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) UpdateTime_Less(v time.Time) *OauthClientQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) UpdateTime_LessEqual(v time.Time) *OauthClientQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) UpdateTime_Greater(v time.Time) *OauthClientQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthClientQuery) UpdateTime_GreaterEqual(v time.Time) *OauthClientQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type OauthClientDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewOauthClientDao(db *DB) (t *OauthClientDao, err error) {
	t = &OauthClientDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *OauthClientDao) init() (err error) {
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
func (dao *OauthClientDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO oauth_client (client_id,account_id,password_hash,redirect_uri,create_time,update_time) VALUES (?,?,?,?,?,?)")
	return err
}

func (dao *OauthClientDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE oauth_client SET client_id=?,account_id=?,password_hash=?,redirect_uri=?,create_time=?,update_time=? WHERE id=?")
	return err
}

func (dao *OauthClientDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM oauth_client WHERE id=?")
	return err
}

func (dao *OauthClientDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthClient) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.ClientId, e.AccountId, e.PasswordHash, e.RedirectUri, e.CreateTime, e.UpdateTime)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *OauthClientDao) Update(ctx context.Context, tx *wrap.Tx, e *OauthClient) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.ClientId, e.AccountId, e.PasswordHash, e.RedirectUri, e.CreateTime, e.UpdateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthClientDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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

func (dao *OauthClientDao) scanRow(row *wrap.Row) (*OauthClient, error) {
	e := &OauthClient{}
	err := row.Scan(&e.Id, &e.ClientId, &e.AccountId, &e.PasswordHash, &e.RedirectUri, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *OauthClientDao) scanRows(rows *wrap.Rows) (list []*OauthClient, err error) {
	list = make([]*OauthClient, 0)
	for rows.Next() {
		e := OauthClient{}
		err = rows.Scan(&e.Id, &e.ClientId, &e.AccountId, &e.PasswordHash, &e.RedirectUri, &e.CreateTime, &e.UpdateTime)
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

func (dao *OauthClientDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*OauthClient, error) {
	querySql := "SELECT " + OAUTH_CLIENT_ALL_FIELDS_STRING + " FROM oauth_client " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *OauthClientDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*OauthClient, err error) {
	querySql := "SELECT " + OAUTH_CLIENT_ALL_FIELDS_STRING + " FROM oauth_client " + query
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

func (dao *OauthClientDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM oauth_client " + query
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

func (dao *OauthClientDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM oauth_client " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *OauthClientDao) GetQuery() *OauthClientQuery {
	return NewOauthClientQuery(dao)
}

const OAUTH_SCOPE_TABLE_NAME = "oauth_scope"

type OAUTH_SCOPE_FIELD string

const OAUTH_SCOPE_FIELD_ID = OAUTH_SCOPE_FIELD("id")
const OAUTH_SCOPE_FIELD_OAUTH_SCOPE = OAUTH_SCOPE_FIELD("oauth_scope")
const OAUTH_SCOPE_FIELD_SCOPE_DESC = OAUTH_SCOPE_FIELD("scope_desc")
const OAUTH_SCOPE_FIELD_CREATE_TIME = OAUTH_SCOPE_FIELD("create_time")
const OAUTH_SCOPE_FIELD_UPDATE_TIME = OAUTH_SCOPE_FIELD("update_time")

const OAUTH_SCOPE_ALL_FIELDS_STRING = "id,oauth_scope,scope_desc,create_time,update_time"

var OAUTH_SCOPE_ALL_FIELDS = []string{
	"id",
	"oauth_scope",
	"scope_desc",
	"create_time",
	"update_time",
}

type OauthScope struct {
	Id         int64  //size=20
	OauthScope string //size=256
	ScopeDesc  string //size=1024
	CreateTime time.Time
	UpdateTime string //size=45
}

type OauthScopeQuery struct {
	BaseQuery
	dao *OauthScopeDao
}

func NewOauthScopeQuery(dao *OauthScopeDao) *OauthScopeQuery {
	q := &OauthScopeQuery{}
	q.dao = dao

	return q
}

func (q *OauthScopeQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*OauthScope, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *OauthScopeQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*OauthScope, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *OauthScopeQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *OauthScopeQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *OauthScopeQuery) ForUpdate() *OauthScopeQuery {
	q.forUpdate = true
	return q
}

func (q *OauthScopeQuery) ForShare() *OauthScopeQuery {
	q.forShare = true
	return q
}

func (q *OauthScopeQuery) GroupBy(fields ...OAUTH_SCOPE_FIELD) *OauthScopeQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *OauthScopeQuery) Limit(startIncluded int64, count int64) *OauthScopeQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *OauthScopeQuery) OrderBy(fieldName OAUTH_SCOPE_FIELD, asc bool) *OauthScopeQuery {
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

func (q *OauthScopeQuery) OrderByGroupCount(asc bool) *OauthScopeQuery {
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

func (q *OauthScopeQuery) w(format string, a ...interface{}) *OauthScopeQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *OauthScopeQuery) Left() *OauthScopeQuery  { return q.w(" ( ") }
func (q *OauthScopeQuery) Right() *OauthScopeQuery { return q.w(" ) ") }
func (q *OauthScopeQuery) And() *OauthScopeQuery   { return q.w(" AND ") }
func (q *OauthScopeQuery) Or() *OauthScopeQuery    { return q.w(" OR ") }
func (q *OauthScopeQuery) Not() *OauthScopeQuery   { return q.w(" NOT ") }

func (q *OauthScopeQuery) Id_Equal(v int64) *OauthScopeQuery { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *OauthScopeQuery) Id_NotEqual(v int64) *OauthScopeQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) Id_Less(v int64) *OauthScopeQuery { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *OauthScopeQuery) Id_LessEqual(v int64) *OauthScopeQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) Id_Greater(v int64) *OauthScopeQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) Id_GreaterEqual(v int64) *OauthScopeQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) OauthScope_Equal(v string) *OauthScopeQuery {
	return q.w("oauth_scope='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) OauthScope_NotEqual(v string) *OauthScopeQuery {
	return q.w("oauth_scope<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) OauthScope_Less(v string) *OauthScopeQuery {
	return q.w("oauth_scope<'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) OauthScope_LessEqual(v string) *OauthScopeQuery {
	return q.w("oauth_scope<='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) OauthScope_Greater(v string) *OauthScopeQuery {
	return q.w("oauth_scope>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) OauthScope_GreaterEqual(v string) *OauthScopeQuery {
	return q.w("oauth_scope>='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) ScopeDesc_Equal(v string) *OauthScopeQuery {
	return q.w("scope_desc='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) ScopeDesc_NotEqual(v string) *OauthScopeQuery {
	return q.w("scope_desc<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) ScopeDesc_Less(v string) *OauthScopeQuery {
	return q.w("scope_desc<'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) ScopeDesc_LessEqual(v string) *OauthScopeQuery {
	return q.w("scope_desc<='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) ScopeDesc_Greater(v string) *OauthScopeQuery {
	return q.w("scope_desc>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) ScopeDesc_GreaterEqual(v string) *OauthScopeQuery {
	return q.w("scope_desc>='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) CreateTime_Equal(v time.Time) *OauthScopeQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) CreateTime_NotEqual(v time.Time) *OauthScopeQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) CreateTime_Less(v time.Time) *OauthScopeQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) CreateTime_LessEqual(v time.Time) *OauthScopeQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) CreateTime_Greater(v time.Time) *OauthScopeQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) CreateTime_GreaterEqual(v time.Time) *OauthScopeQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) UpdateTime_Equal(v string) *OauthScopeQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) UpdateTime_NotEqual(v string) *OauthScopeQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) UpdateTime_Less(v string) *OauthScopeQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) UpdateTime_LessEqual(v string) *OauthScopeQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) UpdateTime_Greater(v string) *OauthScopeQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthScopeQuery) UpdateTime_GreaterEqual(v string) *OauthScopeQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type OauthScopeDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewOauthScopeDao(db *DB) (t *OauthScopeDao, err error) {
	t = &OauthScopeDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *OauthScopeDao) init() (err error) {
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
func (dao *OauthScopeDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO oauth_scope (oauth_scope,scope_desc,create_time,update_time) VALUES (?,?,?,?)")
	return err
}

func (dao *OauthScopeDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE oauth_scope SET oauth_scope=?,scope_desc=?,create_time=?,update_time=? WHERE id=?")
	return err
}

func (dao *OauthScopeDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM oauth_scope WHERE id=?")
	return err
}

func (dao *OauthScopeDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthScope) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.OauthScope, e.ScopeDesc, e.CreateTime, e.UpdateTime)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *OauthScopeDao) Update(ctx context.Context, tx *wrap.Tx, e *OauthScope) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.OauthScope, e.ScopeDesc, e.CreateTime, e.UpdateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthScopeDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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

func (dao *OauthScopeDao) scanRow(row *wrap.Row) (*OauthScope, error) {
	e := &OauthScope{}
	err := row.Scan(&e.Id, &e.OauthScope, &e.ScopeDesc, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *OauthScopeDao) scanRows(rows *wrap.Rows) (list []*OauthScope, err error) {
	list = make([]*OauthScope, 0)
	for rows.Next() {
		e := OauthScope{}
		err = rows.Scan(&e.Id, &e.OauthScope, &e.ScopeDesc, &e.CreateTime, &e.UpdateTime)
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

func (dao *OauthScopeDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*OauthScope, error) {
	querySql := "SELECT " + OAUTH_SCOPE_ALL_FIELDS_STRING + " FROM oauth_scope " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *OauthScopeDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*OauthScope, err error) {
	querySql := "SELECT " + OAUTH_SCOPE_ALL_FIELDS_STRING + " FROM oauth_scope " + query
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

func (dao *OauthScopeDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM oauth_scope " + query
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

func (dao *OauthScopeDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM oauth_scope " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *OauthScopeDao) GetQuery() *OauthScopeQuery {
	return NewOauthScopeQuery(dao)
}

const REFRESH_TOKEN_TABLE_NAME = "refresh_token"

type REFRESH_TOKEN_FIELD string

const REFRESH_TOKEN_FIELD_ID = REFRESH_TOKEN_FIELD("id")
const REFRESH_TOKEN_FIELD_REFRESH_TOKEN = REFRESH_TOKEN_FIELD("refresh_token")
const REFRESH_TOKEN_FIELD_CLIENT_ID = REFRESH_TOKEN_FIELD("client_id")
const REFRESH_TOKEN_FIELD_ACCOUNT_ID = REFRESH_TOKEN_FIELD("account_id")
const REFRESH_TOKEN_FIELD_EXPIRE_SECONDS = REFRESH_TOKEN_FIELD("expire_seconds")
const REFRESH_TOKEN_FIELD_OAUTH_SCOPE = REFRESH_TOKEN_FIELD("oauth_scope")
const REFRESH_TOKEN_FIELD_CREATE_TIME = REFRESH_TOKEN_FIELD("create_time")
const REFRESH_TOKEN_FIELD_UPDATE_TIME = REFRESH_TOKEN_FIELD("update_time")

const REFRESH_TOKEN_ALL_FIELDS_STRING = "id,refresh_token,client_id,account_id,expire_seconds,oauth_scope,create_time,update_time"

var REFRESH_TOKEN_ALL_FIELDS = []string{
	"id",
	"refresh_token",
	"client_id",
	"account_id",
	"expire_seconds",
	"oauth_scope",
	"create_time",
	"update_time",
}

type RefreshToken struct {
	Id            int64  //size=20
	RefreshToken  string //size=128
	ClientId      string //size=32
	AccountId     string //size=32
	ExpireSeconds int64  //size=20
	OauthScope    string //size=256
	CreateTime    time.Time
	UpdateTime    time.Time
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

func (q *RefreshTokenQuery) Id_Equal(v int64) *RefreshTokenQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_NotEqual(v int64) *RefreshTokenQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_Less(v int64) *RefreshTokenQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_LessEqual(v int64) *RefreshTokenQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_Greater(v int64) *RefreshTokenQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_GreaterEqual(v int64) *RefreshTokenQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
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
func (q *RefreshTokenQuery) ClientId_Equal(v string) *RefreshTokenQuery {
	return q.w("client_id='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ClientId_NotEqual(v string) *RefreshTokenQuery {
	return q.w("client_id<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ClientId_Less(v string) *RefreshTokenQuery {
	return q.w("client_id<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ClientId_LessEqual(v string) *RefreshTokenQuery {
	return q.w("client_id<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ClientId_Greater(v string) *RefreshTokenQuery {
	return q.w("client_id>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ClientId_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("client_id>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_Equal(v string) *RefreshTokenQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_NotEqual(v string) *RefreshTokenQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_Less(v string) *RefreshTokenQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_LessEqual(v string) *RefreshTokenQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_Greater(v string) *RefreshTokenQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ExpireSeconds_Equal(v int64) *RefreshTokenQuery {
	return q.w("expire_seconds='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ExpireSeconds_NotEqual(v int64) *RefreshTokenQuery {
	return q.w("expire_seconds<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ExpireSeconds_Less(v int64) *RefreshTokenQuery {
	return q.w("expire_seconds<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ExpireSeconds_LessEqual(v int64) *RefreshTokenQuery {
	return q.w("expire_seconds<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ExpireSeconds_Greater(v int64) *RefreshTokenQuery {
	return q.w("expire_seconds>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) ExpireSeconds_GreaterEqual(v int64) *RefreshTokenQuery {
	return q.w("expire_seconds>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthScope_Equal(v string) *RefreshTokenQuery {
	return q.w("oauth_scope='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthScope_NotEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_scope<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthScope_Less(v string) *RefreshTokenQuery {
	return q.w("oauth_scope<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthScope_LessEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_scope<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthScope_Greater(v string) *RefreshTokenQuery {
	return q.w("oauth_scope>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthScope_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_scope>='" + fmt.Sprint(v) + "'")
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
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO refresh_token (refresh_token,client_id,account_id,expire_seconds,oauth_scope,create_time,update_time) VALUES (?,?,?,?,?,?,?)")
	return err
}

func (dao *RefreshTokenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE refresh_token SET refresh_token=?,client_id=?,account_id=?,expire_seconds=?,oauth_scope=?,create_time=?,update_time=? WHERE id=?")
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

	result, err := stmt.Exec(ctx, e.RefreshToken, e.ClientId, e.AccountId, e.ExpireSeconds, e.OauthScope, e.CreateTime, e.UpdateTime)
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

	_, err = stmt.Exec(ctx, e.RefreshToken, e.ClientId, e.AccountId, e.ExpireSeconds, e.OauthScope, e.CreateTime, e.UpdateTime, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *RefreshTokenDao) Delete(ctx context.Context, tx *wrap.Tx, id int64) (err error) {
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
	err := row.Scan(&e.Id, &e.RefreshToken, &e.ClientId, &e.AccountId, &e.ExpireSeconds, &e.OauthScope, &e.CreateTime, &e.UpdateTime)
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
		err = rows.Scan(&e.Id, &e.RefreshToken, &e.ClientId, &e.AccountId, &e.ExpireSeconds, &e.OauthScope, &e.CreateTime, &e.UpdateTime)
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

type DB struct {
	wrap.DB
	AccessToken       *AccessTokenDao
	AuthorizationCode *AuthorizationCodeDao
	OauthClient       *OauthClientDao
	OauthScope        *OauthScopeDao
	RefreshToken      *RefreshTokenDao
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

	d.AccessToken, err = NewAccessTokenDao(d)
	if err != nil {
		return nil, err
	}

	d.AuthorizationCode, err = NewAuthorizationCodeDao(d)
	if err != nil {
		return nil, err
	}

	d.OauthClient, err = NewOauthClientDao(d)
	if err != nil {
		return nil, err
	}

	d.OauthScope, err = NewOauthScopeDao(d)
	if err != nil {
		return nil, err
	}

	d.RefreshToken, err = NewRefreshTokenDao(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
