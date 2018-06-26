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

type QueryBase struct {
	where              *bytes.Buffer
	whereParams        []interface{}
	groupByFields      []string
	groupByOrders      []bool
	orderByFields      []string
	orderByOrders      []bool
	hasLimit           bool
	limitStartIncluded int64
	limitCount         int64
	forUpdate          bool
	forShare           bool
	updateFields       []string
	updateParams       []interface{}
}

func (q *QueryBase) buildSelectQuery() (queryString string, params []interface{}) {
	query := bytes.NewBufferString("")

	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams...)
	}

	groupByCount := len(q.groupByFields)
	if groupByCount > 0 {
		groupByItems := make([]string, groupByCount)
		for i, v := range q.groupByFields {
			if q.groupByOrders[i] {
				groupByItems[i] = v + " ASC"
			} else {
				groupByItems[i] = v + " DESC"
			}
		}
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(groupByItems, ","))
	}

	orderByCount := len(q.orderByFields)
	if orderByCount > 0 {
		orderByItems := make([]string, orderByCount)
		for i, v := range q.orderByFields {
			if q.orderByOrders[i] {
				orderByItems[i] = v + " ASC"
			} else {
				orderByItems[i] = v + " DESC"
			}
		}
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(orderByItems, ","))
	}

	if q.hasLimit {
		query.WriteString(fmt.Sprintf(" LIMIT %d,%d", q.limitStartIncluded, q.limitCount))
	}

	if q.forUpdate {
		query.WriteString(" FOR UPDATE")
	}

	if q.forShare {
		query.WriteString(" LOCK IN SHARE MODE")
	}

	return query.String(), params
}

type AccessToken struct {
	Id          uint64 //size=20
	UserId      string //size=32
	AccessToken string //size=1024
	CreateTime  time.Time
	UpdateTime  time.Time
}

type AccessTokenQuery struct {
	QueryBase
	dao *AccessTokenDao
}

func (dao *AccessTokenDao) Query() *AccessTokenQuery {
	q := &AccessTokenQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *AccessTokenQuery) Left() *AccessTokenQuery {
	q.where.WriteString(" (")
	return q
}

func (q *AccessTokenQuery) Right() *AccessTokenQuery {
	q.where.WriteString(" )")
	return q
}

func (q *AccessTokenQuery) And() *AccessTokenQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *AccessTokenQuery) Or() *AccessTokenQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *AccessTokenQuery) Not() *AccessTokenQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *AccessTokenQuery) IdEqual(v uint64) *AccessTokenQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) IdNotEqual(v uint64) *AccessTokenQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) IdLess(v uint64) *AccessTokenQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) IdLessEqual(v uint64) *AccessTokenQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) IdGreater(v uint64) *AccessTokenQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) IdGreaterEqual(v uint64) *AccessTokenQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) IdIn(items []uint64) *AccessTokenQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccessTokenQuery) UserIdEqual(v string) *AccessTokenQuery {
	q.where.WriteString(" user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UserIdNotEqual(v string) *AccessTokenQuery {
	q.where.WriteString(" user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UserIdIn(items []string) *AccessTokenQuery {
	q.where.WriteString(" user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccessTokenQuery) AccessTokenEqual(v string) *AccessTokenQuery {
	q.where.WriteString(" access_token=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) AccessTokenNotEqual(v string) *AccessTokenQuery {
	q.where.WriteString(" access_token<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) AccessTokenIn(items []string) *AccessTokenQuery {
	q.where.WriteString(" access_token IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccessTokenQuery) CreateTimeEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) CreateTimeNotEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) CreateTimeLess(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) CreateTimeLessEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) CreateTimeGreater(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) CreateTimeGreaterEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UpdateTimeEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" update_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UpdateTimeNotEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" update_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UpdateTimeLess(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" update_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UpdateTimeLessEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" update_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UpdateTimeGreater(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" update_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) UpdateTimeGreaterEqual(v time.Time) *AccessTokenQuery {
	q.where.WriteString(" update_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccessTokenQuery) GroupByUserId(asc bool) *AccessTokenQuery {
	q.groupByFields = append(q.groupByFields, "user_id")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *AccessTokenQuery) OrderById(asc bool) *AccessTokenQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccessTokenQuery) OrderByUserId(asc bool) *AccessTokenQuery {
	q.orderByFields = append(q.orderByFields, "user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccessTokenQuery) OrderByAccessToken(asc bool) *AccessTokenQuery {
	q.orderByFields = append(q.orderByFields, "access_token")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccessTokenQuery) OrderByCreateTime(asc bool) *AccessTokenQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccessTokenQuery) OrderByUpdateTime(asc bool) *AccessTokenQuery {
	q.orderByFields = append(q.orderByFields, "update_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccessTokenQuery) OrderByGroupCount(asc bool) *AccessTokenQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccessTokenQuery) Limit(startIncluded int64, count int64) *AccessTokenQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *AccessTokenQuery) ForUpdate() *AccessTokenQuery {
	q.forUpdate = true
	return q
}

func (q *AccessTokenQuery) ForShare() *AccessTokenQuery {
	q.forShare = true
	return q
}

func (q *AccessTokenQuery) Select(ctx context.Context, tx *wrap.Tx) (e *AccessToken, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,access_token,create_time,update_time FROM access_token ")
	query.WriteString(queryString)
	e = &AccessToken{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.UserId, &e.AccessToken, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *AccessTokenQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*AccessToken, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,access_token,create_time,update_time FROM access_token ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
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

func (q *AccessTokenQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM access_token ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *AccessTokenQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM access_token ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *AccessTokenQuery) SetUserId(v string) *AccessTokenQuery {
	q.updateFields = append(q.updateFields, "user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccessTokenQuery) SetAccessToken(v string) *AccessTokenQuery {
	q.updateFields = append(q.updateFields, "access_token")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccessTokenQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE access_token SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *AccessTokenQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM access_token WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type AccessTokenDao struct {
	logger *zap.Logger
	db     *DB
}

func NewAccessTokenDao(db *DB) (t *AccessTokenDao, err error) {
	t = &AccessTokenDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *AccessTokenDao) Insert(ctx context.Context, tx *wrap.Tx, e *AccessToken, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO access_token (user_id,access_token) VALUES (?,?)")
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_id=VALUES(user_id)")
	}
	params := []interface{}{e.UserId, e.AccessToken}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *AccessTokenDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*AccessToken, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO access_token (user_id,access_token) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?)", len(list), ","))
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_id=VALUES(user_id)")
	}
	params := make([]interface{}, len(list)*2)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.UserId
		params[offset+1] = e.AccessToken
		offset += 2
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *AccessTokenDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM AccessToken WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *AccessTokenDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *AccessToken) (result *wrap.Result, err error) {
	query := "UPDATE access_token SET user_id=?,access_token=? WHERE id=?"
	params := []interface{}{e.UserId, e.AccessToken, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *AccessTokenDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *AccessToken, err error) {
	query := "SELECT id,user_id,access_token,create_time,update_time FROM access_token WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &AccessToken{}
	err = row.Scan(&e.Id, &e.UserId, &e.AccessToken, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
}

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
	QueryBase
	dao *AccountOperationDao
}

func (dao *AccountOperationDao) Query() *AccountOperationQuery {
	q := &AccountOperationQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *AccountOperationQuery) Left() *AccountOperationQuery {
	q.where.WriteString(" (")
	return q
}

func (q *AccountOperationQuery) Right() *AccountOperationQuery {
	q.where.WriteString(" )")
	return q
}

func (q *AccountOperationQuery) And() *AccountOperationQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *AccountOperationQuery) Or() *AccountOperationQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *AccountOperationQuery) Not() *AccountOperationQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *AccountOperationQuery) IdEqual(v uint64) *AccountOperationQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) IdNotEqual(v uint64) *AccountOperationQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) IdLess(v uint64) *AccountOperationQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) IdLessEqual(v uint64) *AccountOperationQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) IdGreater(v uint64) *AccountOperationQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) IdGreaterEqual(v uint64) *AccountOperationQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) IdIn(items []uint64) *AccountOperationQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccountOperationQuery) UserIdEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) UserIdNotEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) UserIdIn(items []string) *AccountOperationQuery {
	q.where.WriteString(" user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccountOperationQuery) OperationTypeEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" operationType=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) OperationTypeNotEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" operationType<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) OperationTypeIn(items []string) *AccountOperationQuery {
	q.where.WriteString(" operationType IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccountOperationQuery) UserAgentEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" user_agent=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) UserAgentNotEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" user_agent<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) UserAgentIn(items []string) *AccountOperationQuery {
	q.where.WriteString(" user_agent IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccountOperationQuery) PhoneEncryptedEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" phone_encrypted=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) PhoneEncryptedNotEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" phone_encrypted<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) PhoneEncryptedIn(items []string) *AccountOperationQuery {
	q.where.WriteString(" phone_encrypted IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccountOperationQuery) SmsSceneEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" sms_scene=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) SmsSceneNotEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" sms_scene<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) SmsSceneIn(items []string) *AccountOperationQuery {
	q.where.WriteString(" sms_scene IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccountOperationQuery) OtherUserIdEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" other_user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) OtherUserIdNotEqual(v string) *AccountOperationQuery {
	q.where.WriteString(" other_user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) OtherUserIdIn(items []string) *AccountOperationQuery {
	q.where.WriteString(" other_user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *AccountOperationQuery) CreateTimeEqual(v time.Time) *AccountOperationQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) CreateTimeNotEqual(v time.Time) *AccountOperationQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) CreateTimeLess(v time.Time) *AccountOperationQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) CreateTimeLessEqual(v time.Time) *AccountOperationQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) CreateTimeGreater(v time.Time) *AccountOperationQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) CreateTimeGreaterEqual(v time.Time) *AccountOperationQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *AccountOperationQuery) GroupByUserId(asc bool) *AccountOperationQuery {
	q.groupByFields = append(q.groupByFields, "user_id")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *AccountOperationQuery) GroupByOperationType(asc bool) *AccountOperationQuery {
	q.groupByFields = append(q.groupByFields, "operationType")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *AccountOperationQuery) GroupByUserAgent(asc bool) *AccountOperationQuery {
	q.groupByFields = append(q.groupByFields, "user_agent")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *AccountOperationQuery) GroupByPhoneEncrypted(asc bool) *AccountOperationQuery {
	q.groupByFields = append(q.groupByFields, "phone_encrypted")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *AccountOperationQuery) GroupBySmsScene(asc bool) *AccountOperationQuery {
	q.groupByFields = append(q.groupByFields, "sms_scene")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *AccountOperationQuery) GroupByOtherUserId(asc bool) *AccountOperationQuery {
	q.groupByFields = append(q.groupByFields, "other_user_id")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderById(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderByUserId(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderByOperationType(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "operationType")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderByUserAgent(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "user_agent")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderByPhoneEncrypted(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "phone_encrypted")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderBySmsScene(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "sms_scene")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderByOtherUserId(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "other_user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderByCreateTime(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) OrderByGroupCount(asc bool) *AccountOperationQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *AccountOperationQuery) Limit(startIncluded int64, count int64) *AccountOperationQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *AccountOperationQuery) ForUpdate() *AccountOperationQuery {
	q.forUpdate = true
	return q
}

func (q *AccountOperationQuery) ForShare() *AccountOperationQuery {
	q.forShare = true
	return q
}

func (q *AccountOperationQuery) Select(ctx context.Context, tx *wrap.Tx) (e *AccountOperation, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,operationType,user_agent,phone_encrypted,sms_scene,other_user_id,create_time FROM account_operation ")
	query.WriteString(queryString)
	e = &AccountOperation{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.UserId, &e.OperationType, &e.UserAgent, &e.PhoneEncrypted, &e.SmsScene, &e.OtherUserId, &e.CreateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *AccountOperationQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*AccountOperation, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,operationType,user_agent,phone_encrypted,sms_scene,other_user_id,create_time FROM account_operation ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
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

func (q *AccountOperationQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM account_operation ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *AccountOperationQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM account_operation ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *AccountOperationQuery) SetUserId(v string) *AccountOperationQuery {
	q.updateFields = append(q.updateFields, "user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccountOperationQuery) SetOperationType(v string) *AccountOperationQuery {
	q.updateFields = append(q.updateFields, "operationType")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccountOperationQuery) SetUserAgent(v string) *AccountOperationQuery {
	q.updateFields = append(q.updateFields, "user_agent")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccountOperationQuery) SetPhoneEncrypted(v string) *AccountOperationQuery {
	q.updateFields = append(q.updateFields, "phone_encrypted")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccountOperationQuery) SetSmsScene(v string) *AccountOperationQuery {
	q.updateFields = append(q.updateFields, "sms_scene")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccountOperationQuery) SetOtherUserId(v string) *AccountOperationQuery {
	q.updateFields = append(q.updateFields, "other_user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *AccountOperationQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE account_operation SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *AccountOperationQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM account_operation WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type AccountOperationDao struct {
	logger *zap.Logger
	db     *DB
}

func NewAccountOperationDao(db *DB) (t *AccountOperationDao, err error) {
	t = &AccountOperationDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *AccountOperationDao) Insert(ctx context.Context, tx *wrap.Tx, e *AccountOperation) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO account_operation (user_id,operationType,user_agent,phone_encrypted,sms_scene,other_user_id) VALUES (?,?,?,?,?,?)")
	params := []interface{}{e.UserId, e.OperationType, e.UserAgent, e.PhoneEncrypted, e.SmsScene, e.OtherUserId}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *AccountOperationDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*AccountOperation) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO account_operation (user_id,operationType,user_agent,phone_encrypted,sms_scene,other_user_id) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?,?,?,?,?)", len(list), ","))
	params := make([]interface{}, len(list)*6)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.UserId
		params[offset+1] = e.OperationType
		params[offset+2] = e.UserAgent
		params[offset+3] = e.PhoneEncrypted
		params[offset+4] = e.SmsScene
		params[offset+5] = e.OtherUserId
		offset += 6
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *AccountOperationDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM AccountOperation WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *AccountOperationDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *AccountOperation) (result *wrap.Result, err error) {
	query := "UPDATE account_operation SET user_id=?,operationType=?,user_agent=?,phone_encrypted=?,sms_scene=?,other_user_id=? WHERE id=?"
	params := []interface{}{e.UserId, e.OperationType, e.UserAgent, e.PhoneEncrypted, e.SmsScene, e.OtherUserId, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *AccountOperationDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *AccountOperation, err error) {
	query := "SELECT id,user_id,operationType,user_agent,phone_encrypted,sms_scene,other_user_id,create_time FROM account_operation WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &AccountOperation{}
	err = row.Scan(&e.Id, &e.UserId, &e.OperationType, &e.UserAgent, &e.PhoneEncrypted, &e.SmsScene, &e.OtherUserId, &e.CreateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
}

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
	QueryBase
	dao *OauthAccountDao
}

func (dao *OauthAccountDao) Query() *OauthAccountQuery {
	q := &OauthAccountQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *OauthAccountQuery) Left() *OauthAccountQuery {
	q.where.WriteString(" (")
	return q
}

func (q *OauthAccountQuery) Right() *OauthAccountQuery {
	q.where.WriteString(" )")
	return q
}

func (q *OauthAccountQuery) And() *OauthAccountQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *OauthAccountQuery) Or() *OauthAccountQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *OauthAccountQuery) Not() *OauthAccountQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *OauthAccountQuery) IdEqual(v uint64) *OauthAccountQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) IdNotEqual(v uint64) *OauthAccountQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) IdLess(v uint64) *OauthAccountQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) IdLessEqual(v uint64) *OauthAccountQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) IdGreater(v uint64) *OauthAccountQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) IdGreaterEqual(v uint64) *OauthAccountQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) IdIn(items []uint64) *OauthAccountQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthAccountQuery) UserIdEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UserIdNotEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UserIdIn(items []string) *OauthAccountQuery {
	q.where.WriteString(" user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthAccountQuery) ProviderIdEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" providerId=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) ProviderIdNotEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" providerId<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) ProviderIdIn(items []string) *OauthAccountQuery {
	q.where.WriteString(" providerId IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthAccountQuery) ProviderNameEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" provider_name=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) ProviderNameNotEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" provider_name<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) ProviderNameIn(items []string) *OauthAccountQuery {
	q.where.WriteString(" provider_name IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthAccountQuery) OpenIdEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" open_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) OpenIdNotEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" open_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) OpenIdIn(items []string) *OauthAccountQuery {
	q.where.WriteString(" open_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthAccountQuery) UserNameEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" user_name=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UserNameNotEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" user_name<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UserNameIn(items []string) *OauthAccountQuery {
	q.where.WriteString(" user_name IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthAccountQuery) UserIconEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" user_icon=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UserIconNotEqual(v string) *OauthAccountQuery {
	q.where.WriteString(" user_icon<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UserIconIn(items []string) *OauthAccountQuery {
	q.where.WriteString(" user_icon IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthAccountQuery) CreateTimeEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) CreateTimeNotEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) CreateTimeLess(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) CreateTimeLessEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) CreateTimeGreater(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) CreateTimeGreaterEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UpdateTimeEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" update_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UpdateTimeNotEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" update_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UpdateTimeLess(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" update_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UpdateTimeLessEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" update_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UpdateTimeGreater(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" update_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) UpdateTimeGreaterEqual(v time.Time) *OauthAccountQuery {
	q.where.WriteString(" update_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthAccountQuery) GroupByUserId(asc bool) *OauthAccountQuery {
	q.groupByFields = append(q.groupByFields, "user_id")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthAccountQuery) GroupByProviderId(asc bool) *OauthAccountQuery {
	q.groupByFields = append(q.groupByFields, "providerId")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthAccountQuery) GroupByProviderName(asc bool) *OauthAccountQuery {
	q.groupByFields = append(q.groupByFields, "provider_name")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthAccountQuery) GroupByOpenId(asc bool) *OauthAccountQuery {
	q.groupByFields = append(q.groupByFields, "open_id")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthAccountQuery) GroupByUserName(asc bool) *OauthAccountQuery {
	q.groupByFields = append(q.groupByFields, "user_name")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthAccountQuery) GroupByUserIcon(asc bool) *OauthAccountQuery {
	q.groupByFields = append(q.groupByFields, "user_icon")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderById(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByUserId(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByProviderId(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "providerId")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByProviderName(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "provider_name")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByOpenId(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "open_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByUserName(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "user_name")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByUserIcon(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "user_icon")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByCreateTime(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByUpdateTime(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "update_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) OrderByGroupCount(asc bool) *OauthAccountQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthAccountQuery) Limit(startIncluded int64, count int64) *OauthAccountQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *OauthAccountQuery) ForUpdate() *OauthAccountQuery {
	q.forUpdate = true
	return q
}

func (q *OauthAccountQuery) ForShare() *OauthAccountQuery {
	q.forShare = true
	return q
}

func (q *OauthAccountQuery) Select(ctx context.Context, tx *wrap.Tx) (e *OauthAccount, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,providerId,provider_name,open_id,user_name,user_icon,create_time,update_time FROM oauth_account ")
	query.WriteString(queryString)
	e = &OauthAccount{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.UserId, &e.ProviderId, &e.ProviderName, &e.OpenId, &e.UserName, &e.UserIcon, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *OauthAccountQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*OauthAccount, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,providerId,provider_name,open_id,user_name,user_icon,create_time,update_time FROM oauth_account ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
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

func (q *OauthAccountQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM oauth_account ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *OauthAccountQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM oauth_account ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *OauthAccountQuery) SetUserId(v string) *OauthAccountQuery {
	q.updateFields = append(q.updateFields, "user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthAccountQuery) SetProviderId(v string) *OauthAccountQuery {
	q.updateFields = append(q.updateFields, "providerId")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthAccountQuery) SetProviderName(v string) *OauthAccountQuery {
	q.updateFields = append(q.updateFields, "provider_name")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthAccountQuery) SetOpenId(v string) *OauthAccountQuery {
	q.updateFields = append(q.updateFields, "open_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthAccountQuery) SetUserName(v string) *OauthAccountQuery {
	q.updateFields = append(q.updateFields, "user_name")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthAccountQuery) SetUserIcon(v string) *OauthAccountQuery {
	q.updateFields = append(q.updateFields, "user_icon")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthAccountQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE oauth_account SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *OauthAccountQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM oauth_account WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type OauthAccountDao struct {
	logger *zap.Logger
	db     *DB
}

func NewOauthAccountDao(db *DB) (t *OauthAccountDao, err error) {
	t = &OauthAccountDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *OauthAccountDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthAccount, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO oauth_account (user_id,providerId,provider_name,open_id,user_name,user_icon) VALUES (?,?,?,?,?,?)")
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_id=VALUES(user_id),provider_name=VALUES(provider_name),user_name=VALUES(user_name),user_icon=VALUES(user_icon)")
	}
	params := []interface{}{e.UserId, e.ProviderId, e.ProviderName, e.OpenId, e.UserName, e.UserIcon}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *OauthAccountDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*OauthAccount, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO oauth_account (user_id,providerId,provider_name,open_id,user_name,user_icon) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?,?,?,?,?)", len(list), ","))
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_id=VALUES(user_id),provider_name=VALUES(provider_name),user_name=VALUES(user_name),user_icon=VALUES(user_icon)")
	}
	params := make([]interface{}, len(list)*6)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.UserId
		params[offset+1] = e.ProviderId
		params[offset+2] = e.ProviderName
		params[offset+3] = e.OpenId
		params[offset+4] = e.UserName
		params[offset+5] = e.UserIcon
		offset += 6
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *OauthAccountDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM OauthAccount WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *OauthAccountDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *OauthAccount) (result *wrap.Result, err error) {
	query := "UPDATE oauth_account SET user_id=?,providerId=?,provider_name=?,open_id=?,user_name=?,user_icon=? WHERE id=?"
	params := []interface{}{e.UserId, e.ProviderId, e.ProviderName, e.OpenId, e.UserName, e.UserIcon, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *OauthAccountDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *OauthAccount, err error) {
	query := "SELECT id,user_id,providerId,provider_name,open_id,user_name,user_icon,create_time,update_time FROM oauth_account WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &OauthAccount{}
	err = row.Scan(&e.Id, &e.UserId, &e.ProviderId, &e.ProviderName, &e.OpenId, &e.UserName, &e.UserIcon, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
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
	QueryBase
	dao *OauthStateDao
}

func (dao *OauthStateDao) Query() *OauthStateQuery {
	q := &OauthStateQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *OauthStateQuery) Left() *OauthStateQuery {
	q.where.WriteString(" (")
	return q
}

func (q *OauthStateQuery) Right() *OauthStateQuery {
	q.where.WriteString(" )")
	return q
}

func (q *OauthStateQuery) And() *OauthStateQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *OauthStateQuery) Or() *OauthStateQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *OauthStateQuery) Not() *OauthStateQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *OauthStateQuery) IdEqual(v uint64) *OauthStateQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IdNotEqual(v uint64) *OauthStateQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IdLess(v uint64) *OauthStateQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IdLessEqual(v uint64) *OauthStateQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IdGreater(v uint64) *OauthStateQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IdGreaterEqual(v uint64) *OauthStateQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IdIn(items []uint64) *OauthStateQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthStateQuery) OauthStateEqual(v string) *OauthStateQuery {
	q.where.WriteString(" oauth_state=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) OauthStateNotEqual(v string) *OauthStateQuery {
	q.where.WriteString(" oauth_state<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) OauthStateIn(items []string) *OauthStateQuery {
	q.where.WriteString(" oauth_state IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthStateQuery) IsUsedEqual(v int32) *OauthStateQuery {
	q.where.WriteString(" is_used=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IsUsedNotEqual(v int32) *OauthStateQuery {
	q.where.WriteString(" is_used<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IsUsedLess(v int32) *OauthStateQuery {
	q.where.WriteString(" is_used<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IsUsedLessEqual(v int32) *OauthStateQuery {
	q.where.WriteString(" is_used<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IsUsedGreater(v int32) *OauthStateQuery {
	q.where.WriteString(" is_used>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IsUsedGreaterEqual(v int32) *OauthStateQuery {
	q.where.WriteString(" is_used>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) IsUsedIn(items []int32) *OauthStateQuery {
	q.where.WriteString(" is_used IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthStateQuery) UserAgentEqual(v string) *OauthStateQuery {
	q.where.WriteString(" user_agent=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UserAgentNotEqual(v string) *OauthStateQuery {
	q.where.WriteString(" user_agent<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UserAgentIn(items []string) *OauthStateQuery {
	q.where.WriteString(" user_agent IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *OauthStateQuery) CreateTimeEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) CreateTimeNotEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) CreateTimeLess(v time.Time) *OauthStateQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) CreateTimeLessEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) CreateTimeGreater(v time.Time) *OauthStateQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) CreateTimeGreaterEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UpdateTimeEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" update_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UpdateTimeNotEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" update_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UpdateTimeLess(v time.Time) *OauthStateQuery {
	q.where.WriteString(" update_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UpdateTimeLessEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" update_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UpdateTimeGreater(v time.Time) *OauthStateQuery {
	q.where.WriteString(" update_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) UpdateTimeGreaterEqual(v time.Time) *OauthStateQuery {
	q.where.WriteString(" update_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *OauthStateQuery) GroupByIsUsed(asc bool) *OauthStateQuery {
	q.groupByFields = append(q.groupByFields, "is_used")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthStateQuery) GroupByUserAgent(asc bool) *OauthStateQuery {
	q.groupByFields = append(q.groupByFields, "user_agent")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *OauthStateQuery) OrderById(asc bool) *OauthStateQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthStateQuery) OrderByOauthState(asc bool) *OauthStateQuery {
	q.orderByFields = append(q.orderByFields, "oauth_state")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthStateQuery) OrderByIsUsed(asc bool) *OauthStateQuery {
	q.orderByFields = append(q.orderByFields, "is_used")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthStateQuery) OrderByUserAgent(asc bool) *OauthStateQuery {
	q.orderByFields = append(q.orderByFields, "user_agent")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthStateQuery) OrderByCreateTime(asc bool) *OauthStateQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthStateQuery) OrderByUpdateTime(asc bool) *OauthStateQuery {
	q.orderByFields = append(q.orderByFields, "update_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthStateQuery) OrderByGroupCount(asc bool) *OauthStateQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *OauthStateQuery) Limit(startIncluded int64, count int64) *OauthStateQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *OauthStateQuery) ForUpdate() *OauthStateQuery {
	q.forUpdate = true
	return q
}

func (q *OauthStateQuery) ForShare() *OauthStateQuery {
	q.forShare = true
	return q
}

func (q *OauthStateQuery) Select(ctx context.Context, tx *wrap.Tx) (e *OauthState, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,oauth_state,is_used,user_agent,create_time,update_time FROM oauth_state ")
	query.WriteString(queryString)
	e = &OauthState{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.OauthState, &e.IsUsed, &e.UserAgent, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *OauthStateQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*OauthState, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,oauth_state,is_used,user_agent,create_time,update_time FROM oauth_state ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
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

func (q *OauthStateQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM oauth_state ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *OauthStateQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM oauth_state ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *OauthStateQuery) SetOauthState(v string) *OauthStateQuery {
	q.updateFields = append(q.updateFields, "oauth_state")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthStateQuery) SetIsUsed(v int32) *OauthStateQuery {
	q.updateFields = append(q.updateFields, "is_used")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthStateQuery) SetUserAgent(v string) *OauthStateQuery {
	q.updateFields = append(q.updateFields, "user_agent")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *OauthStateQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE oauth_state SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *OauthStateQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM oauth_state WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type OauthStateDao struct {
	logger *zap.Logger
	db     *DB
}

func NewOauthStateDao(db *DB) (t *OauthStateDao, err error) {
	t = &OauthStateDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *OauthStateDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthState, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO oauth_state (oauth_state,is_used,user_agent) VALUES (?,?,?)")
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE is_used=VALUES(is_used),user_agent=VALUES(user_agent)")
	}
	params := []interface{}{e.OauthState, e.IsUsed, e.UserAgent}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *OauthStateDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*OauthState, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO oauth_state (oauth_state,is_used,user_agent) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?,?)", len(list), ","))
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE is_used=VALUES(is_used),user_agent=VALUES(user_agent)")
	}
	params := make([]interface{}, len(list)*3)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.OauthState
		params[offset+1] = e.IsUsed
		params[offset+2] = e.UserAgent
		offset += 3
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *OauthStateDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM OauthState WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *OauthStateDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *OauthState) (result *wrap.Result, err error) {
	query := "UPDATE oauth_state SET oauth_state=?,is_used=?,user_agent=? WHERE id=?"
	params := []interface{}{e.OauthState, e.IsUsed, e.UserAgent, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *OauthStateDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *OauthState, err error) {
	query := "SELECT id,oauth_state,is_used,user_agent,create_time,update_time FROM oauth_state WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &OauthState{}
	err = row.Scan(&e.Id, &e.OauthState, &e.IsUsed, &e.UserAgent, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
}

type PhoneAccount struct {
	Id             uint64 //size=20
	UserId         string //size=32
	PhoneEncrypted string //size=32
	CreateTime     time.Time
	UpdateTime     time.Time
}

type PhoneAccountQuery struct {
	QueryBase
	dao *PhoneAccountDao
}

func (dao *PhoneAccountDao) Query() *PhoneAccountQuery {
	q := &PhoneAccountQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *PhoneAccountQuery) Left() *PhoneAccountQuery {
	q.where.WriteString(" (")
	return q
}

func (q *PhoneAccountQuery) Right() *PhoneAccountQuery {
	q.where.WriteString(" )")
	return q
}

func (q *PhoneAccountQuery) And() *PhoneAccountQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *PhoneAccountQuery) Or() *PhoneAccountQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *PhoneAccountQuery) Not() *PhoneAccountQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *PhoneAccountQuery) IdEqual(v uint64) *PhoneAccountQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) IdNotEqual(v uint64) *PhoneAccountQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) IdLess(v uint64) *PhoneAccountQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) IdLessEqual(v uint64) *PhoneAccountQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) IdGreater(v uint64) *PhoneAccountQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) IdGreaterEqual(v uint64) *PhoneAccountQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) IdIn(items []uint64) *PhoneAccountQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *PhoneAccountQuery) UserIdEqual(v string) *PhoneAccountQuery {
	q.where.WriteString(" user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UserIdNotEqual(v string) *PhoneAccountQuery {
	q.where.WriteString(" user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UserIdIn(items []string) *PhoneAccountQuery {
	q.where.WriteString(" user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *PhoneAccountQuery) PhoneEncryptedEqual(v string) *PhoneAccountQuery {
	q.where.WriteString(" phone_encrypted=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) PhoneEncryptedNotEqual(v string) *PhoneAccountQuery {
	q.where.WriteString(" phone_encrypted<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) PhoneEncryptedIn(items []string) *PhoneAccountQuery {
	q.where.WriteString(" phone_encrypted IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *PhoneAccountQuery) CreateTimeEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) CreateTimeNotEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) CreateTimeLess(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) CreateTimeLessEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) CreateTimeGreater(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) CreateTimeGreaterEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UpdateTimeEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" update_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UpdateTimeNotEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" update_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UpdateTimeLess(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" update_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UpdateTimeLessEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" update_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UpdateTimeGreater(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" update_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) UpdateTimeGreaterEqual(v time.Time) *PhoneAccountQuery {
	q.where.WriteString(" update_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *PhoneAccountQuery) GroupByUserId(asc bool) *PhoneAccountQuery {
	q.groupByFields = append(q.groupByFields, "user_id")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *PhoneAccountQuery) OrderById(asc bool) *PhoneAccountQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *PhoneAccountQuery) OrderByUserId(asc bool) *PhoneAccountQuery {
	q.orderByFields = append(q.orderByFields, "user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *PhoneAccountQuery) OrderByPhoneEncrypted(asc bool) *PhoneAccountQuery {
	q.orderByFields = append(q.orderByFields, "phone_encrypted")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *PhoneAccountQuery) OrderByCreateTime(asc bool) *PhoneAccountQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *PhoneAccountQuery) OrderByUpdateTime(asc bool) *PhoneAccountQuery {
	q.orderByFields = append(q.orderByFields, "update_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *PhoneAccountQuery) OrderByGroupCount(asc bool) *PhoneAccountQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *PhoneAccountQuery) Limit(startIncluded int64, count int64) *PhoneAccountQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *PhoneAccountQuery) ForUpdate() *PhoneAccountQuery {
	q.forUpdate = true
	return q
}

func (q *PhoneAccountQuery) ForShare() *PhoneAccountQuery {
	q.forShare = true
	return q
}

func (q *PhoneAccountQuery) Select(ctx context.Context, tx *wrap.Tx) (e *PhoneAccount, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,phone_encrypted,create_time,update_time FROM phone_account ")
	query.WriteString(queryString)
	e = &PhoneAccount{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.UserId, &e.PhoneEncrypted, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *PhoneAccountQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*PhoneAccount, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,phone_encrypted,create_time,update_time FROM phone_account ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
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

func (q *PhoneAccountQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM phone_account ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *PhoneAccountQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM phone_account ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *PhoneAccountQuery) SetUserId(v string) *PhoneAccountQuery {
	q.updateFields = append(q.updateFields, "user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *PhoneAccountQuery) SetPhoneEncrypted(v string) *PhoneAccountQuery {
	q.updateFields = append(q.updateFields, "phone_encrypted")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *PhoneAccountQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE phone_account SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *PhoneAccountQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM phone_account WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type PhoneAccountDao struct {
	logger *zap.Logger
	db     *DB
}

func NewPhoneAccountDao(db *DB) (t *PhoneAccountDao, err error) {
	t = &PhoneAccountDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *PhoneAccountDao) Insert(ctx context.Context, tx *wrap.Tx, e *PhoneAccount, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO phone_account (user_id,phone_encrypted) VALUES (?,?)")
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_id=VALUES(user_id)")
	}
	params := []interface{}{e.UserId, e.PhoneEncrypted}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *PhoneAccountDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*PhoneAccount, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO phone_account (user_id,phone_encrypted) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?)", len(list), ","))
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_id=VALUES(user_id)")
	}
	params := make([]interface{}, len(list)*2)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.UserId
		params[offset+1] = e.PhoneEncrypted
		offset += 2
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *PhoneAccountDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM PhoneAccount WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *PhoneAccountDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *PhoneAccount) (result *wrap.Result, err error) {
	query := "UPDATE phone_account SET user_id=?,phone_encrypted=? WHERE id=?"
	params := []interface{}{e.UserId, e.PhoneEncrypted, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *PhoneAccountDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *PhoneAccount, err error) {
	query := "SELECT id,user_id,phone_encrypted,create_time,update_time FROM phone_account WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &PhoneAccount{}
	err = row.Scan(&e.Id, &e.UserId, &e.PhoneEncrypted, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
}

type RefreshToken struct {
	Id           uint64 //size=20
	UserId       string //size=32
	RefreshToken string //size=128
	CreateTime   time.Time
	UpdateTime   time.Time
}

type RefreshTokenQuery struct {
	QueryBase
	dao *RefreshTokenDao
}

func (dao *RefreshTokenDao) Query() *RefreshTokenQuery {
	q := &RefreshTokenQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *RefreshTokenQuery) Left() *RefreshTokenQuery {
	q.where.WriteString(" (")
	return q
}

func (q *RefreshTokenQuery) Right() *RefreshTokenQuery {
	q.where.WriteString(" )")
	return q
}

func (q *RefreshTokenQuery) And() *RefreshTokenQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *RefreshTokenQuery) Or() *RefreshTokenQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *RefreshTokenQuery) Not() *RefreshTokenQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *RefreshTokenQuery) IdEqual(v uint64) *RefreshTokenQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) IdNotEqual(v uint64) *RefreshTokenQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) IdLess(v uint64) *RefreshTokenQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) IdLessEqual(v uint64) *RefreshTokenQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) IdGreater(v uint64) *RefreshTokenQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) IdGreaterEqual(v uint64) *RefreshTokenQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) IdIn(items []uint64) *RefreshTokenQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *RefreshTokenQuery) UserIdEqual(v string) *RefreshTokenQuery {
	q.where.WriteString(" user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UserIdNotEqual(v string) *RefreshTokenQuery {
	q.where.WriteString(" user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UserIdIn(items []string) *RefreshTokenQuery {
	q.where.WriteString(" user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *RefreshTokenQuery) RefreshTokenEqual(v string) *RefreshTokenQuery {
	q.where.WriteString(" refresh_token=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) RefreshTokenNotEqual(v string) *RefreshTokenQuery {
	q.where.WriteString(" refresh_token<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) RefreshTokenIn(items []string) *RefreshTokenQuery {
	q.where.WriteString(" refresh_token IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *RefreshTokenQuery) CreateTimeEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) CreateTimeNotEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) CreateTimeLess(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) CreateTimeLessEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) CreateTimeGreater(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) CreateTimeGreaterEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UpdateTimeEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" update_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UpdateTimeNotEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" update_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UpdateTimeLess(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" update_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UpdateTimeLessEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" update_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UpdateTimeGreater(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" update_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) UpdateTimeGreaterEqual(v time.Time) *RefreshTokenQuery {
	q.where.WriteString(" update_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *RefreshTokenQuery) OrderById(asc bool) *RefreshTokenQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *RefreshTokenQuery) OrderByUserId(asc bool) *RefreshTokenQuery {
	q.orderByFields = append(q.orderByFields, "user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *RefreshTokenQuery) OrderByRefreshToken(asc bool) *RefreshTokenQuery {
	q.orderByFields = append(q.orderByFields, "refresh_token")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *RefreshTokenQuery) OrderByCreateTime(asc bool) *RefreshTokenQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *RefreshTokenQuery) OrderByUpdateTime(asc bool) *RefreshTokenQuery {
	q.orderByFields = append(q.orderByFields, "update_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *RefreshTokenQuery) OrderByGroupCount(asc bool) *RefreshTokenQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *RefreshTokenQuery) Limit(startIncluded int64, count int64) *RefreshTokenQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *RefreshTokenQuery) ForUpdate() *RefreshTokenQuery {
	q.forUpdate = true
	return q
}

func (q *RefreshTokenQuery) ForShare() *RefreshTokenQuery {
	q.forShare = true
	return q
}

func (q *RefreshTokenQuery) Select(ctx context.Context, tx *wrap.Tx) (e *RefreshToken, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,refresh_token,create_time,update_time FROM refresh_token ")
	query.WriteString(queryString)
	e = &RefreshToken{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.UserId, &e.RefreshToken, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *RefreshTokenQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*RefreshToken, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,refresh_token,create_time,update_time FROM refresh_token ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		e := RefreshToken{}
		err = rows.Scan(&e.Id, &e.UserId, &e.RefreshToken, &e.CreateTime, &e.UpdateTime)
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

func (q *RefreshTokenQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM refresh_token ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *RefreshTokenQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM refresh_token ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *RefreshTokenQuery) SetUserId(v string) *RefreshTokenQuery {
	q.updateFields = append(q.updateFields, "user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *RefreshTokenQuery) SetRefreshToken(v string) *RefreshTokenQuery {
	q.updateFields = append(q.updateFields, "refresh_token")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *RefreshTokenQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE refresh_token SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *RefreshTokenQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM refresh_token WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type RefreshTokenDao struct {
	logger *zap.Logger
	db     *DB
}

func NewRefreshTokenDao(db *DB) (t *RefreshTokenDao, err error) {
	t = &RefreshTokenDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *RefreshTokenDao) Insert(ctx context.Context, tx *wrap.Tx, e *RefreshToken, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO refresh_token (user_id,refresh_token) VALUES (?,?)")
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE ")
	}
	params := []interface{}{e.UserId, e.RefreshToken}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *RefreshTokenDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*RefreshToken, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO refresh_token (user_id,refresh_token) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?)", len(list), ","))
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE ")
	}
	params := make([]interface{}, len(list)*2)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.UserId
		params[offset+1] = e.RefreshToken
		offset += 2
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *RefreshTokenDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM RefreshToken WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *RefreshTokenDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *RefreshToken) (result *wrap.Result, err error) {
	query := "UPDATE refresh_token SET user_id=?,refresh_token=? WHERE id=?"
	params := []interface{}{e.UserId, e.RefreshToken, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *RefreshTokenDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *RefreshToken, err error) {
	query := "SELECT id,user_id,refresh_token,create_time,update_time FROM refresh_token WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &RefreshToken{}
	err = row.Scan(&e.Id, &e.UserId, &e.RefreshToken, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
}

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
	QueryBase
	dao *SmsCodeDao
}

func (dao *SmsCodeDao) Query() *SmsCodeQuery {
	q := &SmsCodeQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *SmsCodeQuery) Left() *SmsCodeQuery {
	q.where.WriteString(" (")
	return q
}

func (q *SmsCodeQuery) Right() *SmsCodeQuery {
	q.where.WriteString(" )")
	return q
}

func (q *SmsCodeQuery) And() *SmsCodeQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *SmsCodeQuery) Or() *SmsCodeQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *SmsCodeQuery) Not() *SmsCodeQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *SmsCodeQuery) IdEqual(v uint64) *SmsCodeQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) IdNotEqual(v uint64) *SmsCodeQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) IdLess(v uint64) *SmsCodeQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) IdLessEqual(v uint64) *SmsCodeQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) IdGreater(v uint64) *SmsCodeQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) IdGreaterEqual(v uint64) *SmsCodeQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) IdIn(items []uint64) *SmsCodeQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *SmsCodeQuery) SmsSceneEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" sms_scene=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) SmsSceneNotEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" sms_scene<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) SmsSceneIn(items []string) *SmsCodeQuery {
	q.where.WriteString(" sms_scene IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *SmsCodeQuery) PhoneEncryptedEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" phone_encrypted=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) PhoneEncryptedNotEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" phone_encrypted<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) PhoneEncryptedIn(items []string) *SmsCodeQuery {
	q.where.WriteString(" phone_encrypted IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *SmsCodeQuery) SmsCodeEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" sms_code=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) SmsCodeNotEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" sms_code<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) SmsCodeIn(items []string) *SmsCodeQuery {
	q.where.WriteString(" sms_code IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *SmsCodeQuery) UserIdEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UserIdNotEqual(v string) *SmsCodeQuery {
	q.where.WriteString(" user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UserIdIn(items []string) *SmsCodeQuery {
	q.where.WriteString(" user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *SmsCodeQuery) CreateTimeEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) CreateTimeNotEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) CreateTimeLess(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) CreateTimeLessEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) CreateTimeGreater(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) CreateTimeGreaterEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UpdateTimeEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" update_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UpdateTimeNotEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" update_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UpdateTimeLess(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" update_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UpdateTimeLessEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" update_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UpdateTimeGreater(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" update_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) UpdateTimeGreaterEqual(v time.Time) *SmsCodeQuery {
	q.where.WriteString(" update_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *SmsCodeQuery) GroupBySmsScene(asc bool) *SmsCodeQuery {
	q.groupByFields = append(q.groupByFields, "sms_scene")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *SmsCodeQuery) GroupByPhoneEncrypted(asc bool) *SmsCodeQuery {
	q.groupByFields = append(q.groupByFields, "phone_encrypted")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *SmsCodeQuery) GroupBySmsCode(asc bool) *SmsCodeQuery {
	q.groupByFields = append(q.groupByFields, "sms_code")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *SmsCodeQuery) GroupByUserId(asc bool) *SmsCodeQuery {
	q.groupByFields = append(q.groupByFields, "user_id")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderById(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderBySmsScene(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "sms_scene")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderByPhoneEncrypted(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "phone_encrypted")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderBySmsCode(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "sms_code")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderByUserId(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderByCreateTime(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderByUpdateTime(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "update_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) OrderByGroupCount(asc bool) *SmsCodeQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *SmsCodeQuery) Limit(startIncluded int64, count int64) *SmsCodeQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *SmsCodeQuery) ForUpdate() *SmsCodeQuery {
	q.forUpdate = true
	return q
}

func (q *SmsCodeQuery) ForShare() *SmsCodeQuery {
	q.forShare = true
	return q
}

func (q *SmsCodeQuery) Select(ctx context.Context, tx *wrap.Tx) (e *SmsCode, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,sms_scene,phone_encrypted,sms_code,user_id,create_time,update_time FROM sms_code ")
	query.WriteString(queryString)
	e = &SmsCode{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.SmsScene, &e.PhoneEncrypted, &e.SmsCode, &e.UserId, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *SmsCodeQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*SmsCode, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,sms_scene,phone_encrypted,sms_code,user_id,create_time,update_time FROM sms_code ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
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

func (q *SmsCodeQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM sms_code ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *SmsCodeQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM sms_code ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *SmsCodeQuery) SetSmsScene(v string) *SmsCodeQuery {
	q.updateFields = append(q.updateFields, "sms_scene")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *SmsCodeQuery) SetPhoneEncrypted(v string) *SmsCodeQuery {
	q.updateFields = append(q.updateFields, "phone_encrypted")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *SmsCodeQuery) SetSmsCode(v string) *SmsCodeQuery {
	q.updateFields = append(q.updateFields, "sms_code")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *SmsCodeQuery) SetUserId(v string) *SmsCodeQuery {
	q.updateFields = append(q.updateFields, "user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *SmsCodeQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE sms_code SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *SmsCodeQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM sms_code WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type SmsCodeDao struct {
	logger *zap.Logger
	db     *DB
}

func NewSmsCodeDao(db *DB) (t *SmsCodeDao, err error) {
	t = &SmsCodeDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *SmsCodeDao) Insert(ctx context.Context, tx *wrap.Tx, e *SmsCode) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO sms_code (sms_scene,phone_encrypted,sms_code,user_id) VALUES (?,?,?,?)")
	params := []interface{}{e.SmsScene, e.PhoneEncrypted, e.SmsCode, e.UserId}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *SmsCodeDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*SmsCode) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO sms_code (sms_scene,phone_encrypted,sms_code,user_id) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?,?,?)", len(list), ","))
	params := make([]interface{}, len(list)*4)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.SmsScene
		params[offset+1] = e.PhoneEncrypted
		params[offset+2] = e.SmsCode
		params[offset+3] = e.UserId
		offset += 4
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *SmsCodeDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM SmsCode WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *SmsCodeDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *SmsCode) (result *wrap.Result, err error) {
	query := "UPDATE sms_code SET sms_scene=?,phone_encrypted=?,sms_code=?,user_id=? WHERE id=?"
	params := []interface{}{e.SmsScene, e.PhoneEncrypted, e.SmsCode, e.UserId, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *SmsCodeDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *SmsCode, err error) {
	query := "SELECT id,sms_scene,phone_encrypted,sms_code,user_id,create_time,update_time FROM sms_code WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &SmsCode{}
	err = row.Scan(&e.Id, &e.SmsScene, &e.PhoneEncrypted, &e.SmsCode, &e.UserId, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
}

type UserInfo struct {
	Id           uint64 //size=20
	UserId       string //size=32
	UserName     string //size=32
	UserIcon     string //size=256
	PasswordHash string //size=1024
	CreateTime   time.Time
	UpdateTime   time.Time
}

type UserInfoQuery struct {
	QueryBase
	dao *UserInfoDao
}

func (dao *UserInfoDao) Query() *UserInfoQuery {
	q := &UserInfoQuery{}
	q.dao = dao
	q.where = bytes.NewBufferString("")
	return q
}

func (q *UserInfoQuery) Left() *UserInfoQuery {
	q.where.WriteString(" (")
	return q
}

func (q *UserInfoQuery) Right() *UserInfoQuery {
	q.where.WriteString(" )")
	return q
}

func (q *UserInfoQuery) And() *UserInfoQuery {
	q.where.WriteString(" AND")
	return q
}

func (q *UserInfoQuery) Or() *UserInfoQuery {
	q.where.WriteString(" OR")
	return q
}

func (q *UserInfoQuery) Not() *UserInfoQuery {
	q.where.WriteString(" NOT")
	return q
}

func (q *UserInfoQuery) IdEqual(v uint64) *UserInfoQuery {
	q.where.WriteString(" id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) IdNotEqual(v uint64) *UserInfoQuery {
	q.where.WriteString(" id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) IdLess(v uint64) *UserInfoQuery {
	q.where.WriteString(" id<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) IdLessEqual(v uint64) *UserInfoQuery {
	q.where.WriteString(" id<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) IdGreater(v uint64) *UserInfoQuery {
	q.where.WriteString(" id>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) IdGreaterEqual(v uint64) *UserInfoQuery {
	q.where.WriteString(" id>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) IdIn(items []uint64) *UserInfoQuery {
	q.where.WriteString(" id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *UserInfoQuery) UserIdEqual(v string) *UserInfoQuery {
	q.where.WriteString(" user_id=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UserIdNotEqual(v string) *UserInfoQuery {
	q.where.WriteString(" user_id<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UserIdIn(items []string) *UserInfoQuery {
	q.where.WriteString(" user_id IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *UserInfoQuery) UserNameEqual(v string) *UserInfoQuery {
	q.where.WriteString(" user_name=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UserNameNotEqual(v string) *UserInfoQuery {
	q.where.WriteString(" user_name<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UserNameIn(items []string) *UserInfoQuery {
	q.where.WriteString(" user_name IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *UserInfoQuery) UserIconEqual(v string) *UserInfoQuery {
	q.where.WriteString(" user_icon=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UserIconNotEqual(v string) *UserInfoQuery {
	q.where.WriteString(" user_icon<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UserIconIn(items []string) *UserInfoQuery {
	q.where.WriteString(" user_icon IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *UserInfoQuery) PasswordHashEqual(v string) *UserInfoQuery {
	q.where.WriteString(" password_hash=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) PasswordHashNotEqual(v string) *UserInfoQuery {
	q.where.WriteString(" password_hash<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) PasswordHashIn(items []string) *UserInfoQuery {
	q.where.WriteString(" password_hash IN(")
	q.where.WriteString(wrap.RepeatWithSeparator("?", len(items), ","))
	q.where.WriteString(")")
	q.whereParams = append(q.whereParams, items)
	return q
}

func (q *UserInfoQuery) CreateTimeEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" create_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) CreateTimeNotEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" create_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) CreateTimeLess(v time.Time) *UserInfoQuery {
	q.where.WriteString(" create_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) CreateTimeLessEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" create_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) CreateTimeGreater(v time.Time) *UserInfoQuery {
	q.where.WriteString(" create_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) CreateTimeGreaterEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" create_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UpdateTimeEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" update_time=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UpdateTimeNotEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" update_time<>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UpdateTimeLess(v time.Time) *UserInfoQuery {
	q.where.WriteString(" update_time<?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UpdateTimeLessEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" update_time<=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UpdateTimeGreater(v time.Time) *UserInfoQuery {
	q.where.WriteString(" update_time>?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) UpdateTimeGreaterEqual(v time.Time) *UserInfoQuery {
	q.where.WriteString(" update_time>=?")
	q.whereParams = append(q.whereParams, v)
	return q
}

func (q *UserInfoQuery) GroupByUserIcon(asc bool) *UserInfoQuery {
	q.groupByFields = append(q.groupByFields, "user_icon")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *UserInfoQuery) GroupByPasswordHash(asc bool) *UserInfoQuery {
	q.groupByFields = append(q.groupByFields, "password_hash")
	q.groupByOrders = append(q.groupByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderById(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderByUserId(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "user_id")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderByUserName(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "user_name")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderByUserIcon(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "user_icon")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderByPasswordHash(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "password_hash")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderByCreateTime(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "create_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderByUpdateTime(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "update_time")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) OrderByGroupCount(asc bool) *UserInfoQuery {
	q.orderByFields = append(q.orderByFields, "count(*)")
	q.orderByOrders = append(q.orderByOrders, asc)
	return q
}

func (q *UserInfoQuery) Limit(startIncluded int64, count int64) *UserInfoQuery {
	q.hasLimit = true
	q.limitStartIncluded = startIncluded
	q.limitCount = count
	return q
}

func (q *UserInfoQuery) ForUpdate() *UserInfoQuery {
	q.forUpdate = true
	return q
}

func (q *UserInfoQuery) ForShare() *UserInfoQuery {
	q.forShare = true
	return q
}

func (q *UserInfoQuery) Select(ctx context.Context, tx *wrap.Tx) (e *UserInfo, err error) {
	if !q.hasLimit {
		q.limitCount = 1
		q.hasLimit = true
	}

	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,user_name,user_icon,password_hash,create_time,update_time FROM user_info ")
	query.WriteString(queryString)
	e = &UserInfo{}
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&e.Id, &e.UserId, &e.UserName, &e.UserIcon, &e.PasswordHash, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}

	return e, err
}

func (q *UserInfoQuery) SelectList(ctx context.Context, tx *wrap.Tx) (list []*UserInfo, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT id,user_id,user_name,user_icon,password_hash,create_time,update_time FROM user_info ")
	query.WriteString(queryString)
	rows, err := q.dao.db.Query(ctx, tx, query.String(), params...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		e := UserInfo{}
		err = rows.Scan(&e.Id, &e.UserId, &e.UserName, &e.UserIcon, &e.PasswordHash, &e.CreateTime, &e.UpdateTime)
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

func (q *UserInfoQuery) SelectCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT COUNT(*) FROM user_info ")
	query.WriteString(queryString)
	row := q.dao.db.QueryRow(ctx, tx, query.String(), params...)
	err = row.Scan(&count)

	return count, err
}

func (q *UserInfoQuery) SelectGroupBy(ctx context.Context, tx *wrap.Tx, withCount bool) (rows *wrap.Rows, err error) {
	queryString, params := q.buildSelectQuery()
	query := bytes.NewBufferString("")
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(q.groupByFields, ","))
	if withCount {
		query.WriteString(",Count(*) ")
	}
	query.WriteString(" FROM user_info ")
	query.WriteString(queryString)

	return q.dao.db.Query(ctx, tx, query.String(), params...)
}

func (q *UserInfoQuery) SetUserId(v string) *UserInfoQuery {
	q.updateFields = append(q.updateFields, "user_id")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *UserInfoQuery) SetUserName(v string) *UserInfoQuery {
	q.updateFields = append(q.updateFields, "user_name")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *UserInfoQuery) SetUserIcon(v string) *UserInfoQuery {
	q.updateFields = append(q.updateFields, "user_icon")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *UserInfoQuery) SetPasswordHash(v string) *UserInfoQuery {
	q.updateFields = append(q.updateFields, "password_hash")
	q.updateParams = append(q.updateParams, v)
	return q
}

func (q *UserInfoQuery) Update(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	var params []interface{}
	params = append(params, q.updateParams)
	query.WriteString("UPDATE user_info SET ")
	updateItems := make([]string, len(q.updateFields))
	for i, v := range q.updateFields {
		updateItems[i] = v + "=?"
	}
	query.WriteString(strings.Join(updateItems, ","))
	where := q.where.String()
	if where != "" {
		query.WriteString(" WHERE ")
		query.WriteString(where)
		params = append(params, q.whereParams)
	}

	return q.dao.db.Exec(ctx, tx, query.String(), params...)
}

func (q *UserInfoQuery) Delete(ctx context.Context, tx *wrap.Tx) (result *wrap.Result, err error) {
	query := "DELETE FROM user_info WHERE " + q.where.String()
	return q.dao.db.Exec(ctx, tx, query, q.whereParams...)
}

type UserInfoDao struct {
	logger *zap.Logger
	db     *DB
}

func NewUserInfoDao(db *DB) (t *UserInfoDao, err error) {
	t = &UserInfoDao{}
	t.logger = log.TypedLogger(t)
	t.db = db

	return t, nil
}

func (dao *UserInfoDao) Insert(ctx context.Context, tx *wrap.Tx, e *UserInfo, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO user_info (user_id,user_name,user_icon,password_hash) VALUES (?,?,?,?)")
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_icon=VALUES(user_icon),password_hash=VALUES(password_hash)")
	}
	params := []interface{}{e.UserId, e.UserName, e.UserIcon, e.PasswordHash}
	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *UserInfoDao) BatchInsert(ctx context.Context, tx *wrap.Tx, list []*UserInfo, onDuplicatedKeyUpdate bool) (result *wrap.Result, err error) {
	query := bytes.NewBufferString("")
	query.WriteString("INSERT INTO user_info (user_id,user_name,user_icon,password_hash) VALUES ")
	query.WriteString(wrap.RepeatWithSeparator("(?,?,?,?)", len(list), ","))
	if onDuplicatedKeyUpdate {
		query.WriteString(" ON DUPLICATED KEY UPDATE user_icon=VALUES(user_icon),password_hash=VALUES(password_hash)")
	}
	params := make([]interface{}, len(list)*4)
	offset := 0
	for _, e := range list {
		params[offset+0] = e.UserId
		params[offset+1] = e.UserName
		params[offset+2] = e.UserIcon
		params[offset+3] = e.PasswordHash
		offset += 4
	}

	return dao.db.Exec(ctx, tx, query.String(), params...)
}

func (dao *UserInfoDao) DeleteById(ctx context.Context, tx *wrap.Tx, id uint64) (result *wrap.Result, err error) {
	query := "DELETE FROM UserInfo WHERE id=?"
	return dao.db.Exec(ctx, tx, query, id)
}

func (dao *UserInfoDao) UpdateById(ctx context.Context, tx *wrap.Tx, e *UserInfo) (result *wrap.Result, err error) {
	query := "UPDATE user_info SET user_id=?,user_name=?,user_icon=?,password_hash=? WHERE id=?"
	params := []interface{}{e.UserId, e.UserName, e.UserIcon, e.PasswordHash, e.Id}
	return dao.db.Exec(ctx, tx, query, params...)
}

func (dao *UserInfoDao) SelectById(ctx context.Context, tx *wrap.Tx, id int64) (e *UserInfo, err error) {
	query := "SELECT id,user_id,user_name,user_icon,password_hash,create_time,update_time FROM user_info WHERE id=?"
	row := dao.db.QueryRow(ctx, tx, query, id)
	e = &UserInfo{}
	err = row.Scan(&e.Id, &e.UserId, &e.UserName, &e.UserIcon, &e.PasswordHash, &e.CreateTime, &e.UpdateTime)
	if err == wrap.ErrNoRows {
		return nil, nil
	}
	return e, err
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
