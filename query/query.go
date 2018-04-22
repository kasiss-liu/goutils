package query

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type query struct {
	conn     *sql.DB
	tx       *sql.Tx
	txStatus bool
	table    string
	where    []where
	fields   []string
	order    map[string]string
	group    []string
	save     []map[string]interface{}
	offset   int
	limit    int
	sql      string
}

type where struct {
	field   string
	compare string
	value   interface{}
}

type queryResult struct {
	Columns []string
	RowsNum int
	Value   []map[string]string
}
type execResult struct {
	AffectedRows int
	Result       bool
}

//mysql 配置结构
type DbConfig struct {
	Host    string
	Port    int
	User    string
	Passwd  string
	Dbname  string
	Charset string
}

// 获取一个配置结构
func NewDbConfig(h string, p int, u string, pwd string, db string, c string) *DbConfig {
	return &DbConfig{
		Host:    h,
		Port:    p,
		User:    u,
		Passwd:  pwd,
		Dbname:  db,
		Charset: c,
	}
}

//获取一个新的query结构
func NewQuery(conn *sql.DB) *query {
	return &query{
		conn:   conn,
		where:  make([]where, 0, 5),
		fields: make([]string, 0, 10),
		order:  make(map[string]string),
		group:  make([]string, 0, 1),
		save:   make([]map[string]interface{}, 0, 5),
		offset: 0,
		limit:  0,
	}
}

//根据配置，连接db并返回一个query结构
func NewQueryWithConfig(conf *DbConfig) (*query, error) {
	db, err := dbConnection(conf)
	if err != nil {
		return nil, err
	}
	return NewQuery(db), nil

}

//mysql连接
func dbConnection(c *DbConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", c.User, c.Passwd, c.Host, c.Port, c.Dbname, c.Charset))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

//设置数据表
func (q *query) Table(t string) *query {
	q.table = t
	return q
}

//注入sql的where条件
func (q *query) Where(field string, value interface{}, compare string) *query {
	var w where = where{
		field:   field,
		value:   value,
		compare: compare,
	}
	q.where = append(q.where, w)
	return q
}

//设置要查询的字段
func (q *query) Select(fields []string) *query {
	for _, v := range fields {
		if q.inFields(v) {
			continue
		}
		q.fields = append(q.fields, v)
	}
	return q
}

//设置order
func (q *query) OrderBy(m map[string]string) *query {
	for k, v := range m {
		q.order[k] = v
	}
	return q
}

//设置group分组
func (q *query) GroupBy(g []string) *query {
	for _, v := range g {
		q.group = append(q.group, v)
	}
	return q
}

//设置limit值
func (q *query) Limit(i ...int) *query {
	if len(i) == 1 {
		q.limit = i[0]
	}
	if len(i) == 2 {
		q.offset = i[0]
		q.limit = i[1]
	}
	return q
}

//设置offset
func (q *query) Offset(i int) *query {
	q.offset = i
	return q
}

//判断是否含有元素
func (q *query) inFields(field string) bool {
	for _, v := range q.fields {
		if v == field {
			return true
		}
	}
	return false
}

//执行查询
func (q *query) Query() *queryResult {
	//先拼接sql
	q.compactQuery()
	//确认是否是事务处理
	var stmt *sql.Stmt
	var err error
	if q.txStatus {
		stmt, err = q.tx.Prepare(q.sql)
	} else {
		stmt, err = q.conn.Prepare(q.sql)
	}
	//然后执行查询
	var res *sql.Rows
	if stmt != nil && err == nil {
		res, err = stmt.Query()
		defer res.Close()
		stmt.Close()
		if err != nil {
			return nil
		}
		stmt.Close()
	} else {
		return nil
	}
	return q.get(res)
}

func (q *query) get(rows *sql.Rows) *queryResult {
	cols, err := rows.Columns()
	if err != nil {
		println(cols)
		return nil
	}
	rawResult := make([][]byte, len(cols))
	result := make(map[string]string)
	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}
	results := make([]map[string]string, 0, 10)
	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			return nil
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[cols[i]] = ""
			} else {
				result[cols[i]] = string(raw)
			}
		}

		results = append(results, result)
	}
	return &queryResult{
		Columns: cols,
		Value:   results,
		RowsNum: len(results),
	}
}

//查询单条记录
func (q *query) QueryOne() map[string]string {
	q.limit = 1
	//先拼接sql
	q.compactQuery()
	//确认是否是事务处理
	var stmt *sql.Stmt
	var err error
	if q.txStatus {
		stmt, err = q.tx.Prepare(q.sql)
	} else {
		stmt, err = q.conn.Prepare(q.sql)
	}
	//然后执行查询
	var res *sql.Rows
	if stmt != nil {
		res, err = stmt.Query()
		defer res.Close()
		stmt.Close()
		if err != nil {
			return nil
		}
	}
	return q.getOne(res)
}

func (q *query) getOne(rows *sql.Rows) map[string]string {
	cols, _ := rows.Columns()
	rawResult := make([][]byte, len(cols))
	result := make(map[string]string)
	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}
	if rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			return nil
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[cols[i]] = ""
			} else {
				result[cols[i]] = string(raw)
			}
		}
	}
	return result
}

//拼接sql
func (q *query) compactQuery() {
	q.compactSelect()
	q.compactTable()
	q.compactWhere()
	q.compactGroup()
	q.compactOrder()
	q.compactLimit()
}

//拼接select
func (q *query) compactSelect() {
	q.sql = "select "
	if len(q.fields) > 0 {
		for _, field := range q.fields {
			q.sql += field + ","
		}
		q.sql = strings.TrimRight(q.sql, ",")
	} else {
		q.sql += "*"
	}
}

//拼接table
func (q *query) compactTable() {

	q.sql += " from `" + q.table + "` "
}

//拼接where
func (q *query) compactWhere() {
	//拼接where
	q.sql += " where "
	if len(q.where) > 0 {
		for _, where := range q.where {
			q.sql += "`" + where.field + "` " + where.compare + " "
			switch v := where.value.(type) {
			case int:
				q.sql += strconv.Itoa(v)
			case string:
				q.sql += "'" + v + "'"
			case []int:
				var intString string
				intString = "("
				for _, i := range v {
					intString += strconv.Itoa(i) + ","
				}
				intString = strings.Trim(intString, ",")
				intString += ")"
				q.sql += intString
			case []string:
				var sString string
				sString = "("
				for _, s := range v {
					sString += "'" + s + "',"
				}
				sString = strings.Trim(sString, ",")
				sString += ")"
				q.sql += sString
			default:
				q.sql += "1"
			}
		}
	} else {
		q.sql += "1"
	}
}

//拼接groupby
func (q *query) compactGroup() {
	if len(q.group) > 0 {
		q.sql += " group by "
		var gString string
		for _, v := range q.group {
			gString += "`" + v + "`,"
		}
		gString = strings.Trim(gString, ",")
		q.sql += gString
	}
}

func (q *query) compactOrder() {
	if len(q.order) > 0 {
		q.sql += " order by "
		for k, v := range q.order {
			q.sql += "`" + k + "` " + v + ","
		}
		q.sql = strings.Trim(q.sql, ",")
	}
}

//拼接limit
func (q *query) compactLimit() {
	if q.limit > 0 {
		q.sql += " limit " + strconv.Itoa(q.limit)
	}
	if q.offset > 0 {
		q.sql += " offset " + strconv.Itoa(q.offset)
	}
}

func (q *query) Close() error {
	return q.conn.Close()
}
func (q *query) Ping() error {
	return q.conn.Ping()
}

func (q *query) BeginTransaction() error {
	if q.tx == nil {
		tx, err := q.conn.Begin()
		if err != nil {
			return err
		}
		q.tx = tx
	}
	q.txStatus = true
	return nil
}
func (q *query) Commit() error {
	if q.tx != nil {
		err := q.tx.Commit()
		if err != nil {
			return err
		}
		q.tx = nil
	}
	q.txStatus = false
	return nil
}
func (q *query) Rollback() error {
	if q.tx != nil {
		err := q.tx.Rollback()
		if err != nil {
			return err
		}
		q.tx = nil
	}
	q.txStatus = false
	return nil
}
