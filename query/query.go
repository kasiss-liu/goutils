package query

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type query struct {
	conn       *sql.DB
	tx         *sql.Tx
	txStatus   bool
	table      string
	where      []where
	fields     []string
	order      []map[string]string
	group      []string
	offset     int
	limit      int
	sql        string
	stmtValue  []interface{}
	queryLog   []string
	isLogQuery bool
	errors     []string
	lastSql    string
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
func NewDbConfig(host string, port int, user string, pwd string, db string, charset string) *DbConfig {
	return &DbConfig{
		Host:    host,
		Port:    port,
		User:    user,
		Passwd:  pwd,
		Dbname:  db,
		Charset: charset,
	}
}

//获取一个新的query结构
func NewQuery(conn *sql.DB) *query {
	return &query{
		conn:       conn,
		where:      make([]where, 0, 5),
		fields:     make([]string, 0, 10),
		order:      make([]map[string]string, 0, 1),
		group:      make([]string, 0, 1),
		offset:     0,
		limit:      0,
		stmtValue:  make([]interface{}, 0, 10),
		queryLog:   make([]string, 0, 10),
		isLogQuery: false,
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
func (q *query) OrderBy(field string, sort string) *query {
	order := map[string]string{field: sort}
	q.order = append(q.order, order)
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

	stmt, err := q.getStmt()
	//然后执行查询
	var rows *sql.Rows
	if err == nil {
		rows, err = stmt.Query(q.stmtValue...)
		defer rows.Close()
		defer stmt.Close()
		if err != nil {
			return nil
		}
	} else {
		q.saveError(err.Error())
		return nil
	}
	//清理临时数据
	q.resetAll()
	return q.get(rows)
}

//查询单条记录
func (q *query) QueryOne() map[string]string {
	q.limit = 1
	//先拼接sql
	q.compactQuery()

	//获取stmt
	stmt, err := q.getStmt()

	//然后执行查询
	var rows *sql.Rows
	if err != nil {
		q.saveError(err.Error())
		return nil

	}
	rows, err = stmt.Query(q.stmtValue...)
	defer rows.Close()
	defer stmt.Close()
	if err != nil {
		q.saveError(err.Error())
		return nil
	}
	//解析结果 并返回第一条数据
	results := q.get(rows)
	//清理临时数据
	q.resetAll()
	if len(results.Value) > 0 {
		return results.Value[0]
	} else {
		return nil
	}
}

//解析查询结果
func (q *query) get(rows *sql.Rows) *queryResult {
	cols, err := rows.Columns()
	if err != nil {
		println(cols)
		return nil
	}
	rawResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}
	results := make([]map[string]string, 0, 10)

	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			q.saveError(err.Error())
			return nil
		}
		result := make(map[string]string)

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

//根据不同的事务状态 返回不同的处理器
func (q *query) getStmt() (*sql.Stmt, error) {
	var stmt *sql.Stmt
	var err error
	//判断是否处于事务中
	if q.txStatus == true {
		stmt, err = q.tx.Prepare(q.sql)
	} else {
		stmt, err = q.conn.Prepare(q.sql)
	}
	//是否开启query日志 如果开启则保存sql
	if q.isLogQuery {
		q.queryLog = append(q.queryLog, q.sql)
	}
	//保存最后一条执行sql
	q.lastSql = q.sql

	return stmt, err
}

//拼接sql
func (q *query) compactQuery() {

	q.resetStmt()
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
	//将where.value 放入stmtValue中 最后拼接
	q.sql += " where 1 "
	if len(q.where) > 0 {
		for _, where := range q.where {
			q.sql += " and `" + where.field + "` " + where.compare + " "
			//断言分析where.value
			switch v := where.value.(type) {
			case int:
				q.sql += "?"
				q.stmtValue = append(q.stmtValue, strconv.Itoa(v))
			case string:
				q.sql += "?"
				q.stmtValue = append(q.stmtValue, v)
			case []int:
				q.sql += "("
				for _, s := range v {
					q.sql += " ? ,"
					q.stmtValue = append(q.stmtValue, s)
				}
				q.sql = strings.Trim(q.sql, ",")
				q.sql += ") "
			case []string:
				q.sql += "("
				for _, s := range v {
					q.sql += " ? ,"
					q.stmtValue = append(q.stmtValue, s)
				}
				q.sql = strings.Trim(q.sql, ",")
				q.sql += ") "
			}
		}
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

//拼接orderby
func (q *query) compactOrder() {
	if len(q.order) > 0 {
		q.sql += " order by "
		for _, stringMap := range q.order {
			for k, v := range stringMap {
				q.sql += "`" + k + "` " + v + ","
			}
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

//创建新数据
func (q *query) Create(data map[string]interface{}) int {
	//初始化sql
	q.resetStmt()
	keys := q.compactCreateFields(data)
	q.compactCreateData([]map[string]interface{}{data}, keys)
	//初始化stmt
	var stmt *sql.Stmt
	var err error
	var res sql.Result
	stmt, err = q.getStmt()

	//stmt 执行操作
	if err != nil {
		q.saveError(err.Error())
		return 0
	}

	defer stmt.Close()
	res, err = stmt.Exec(q.stmtValue...)
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//分析结果 返回数据
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//清理临时数据
	q.resetAll()
	return int(lastInsertId)
}

//批量创建数据
func (q *query) CreateBatch(data []map[string]interface{}) int {
	//初始化sql
	q.resetStmt()
	keys := q.compactCreateFields(data[0])
	q.compactCreateData(data, keys)
	//初始化stmt
	var stmt *sql.Stmt
	var err error
	stmt, err = q.getStmt()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	defer stmt.Close()
	res, err := stmt.Exec(q.stmtValue...)
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//清理临时数据
	q.resetAll()
	return int(lastInsertId)
}

//整理新建数据字段
func (q *query) compactCreateFields(data map[string]interface{}) []string {
	keys := make([]string, 0, 10)
	//准备sql的fields
	q.sql = "insert into `" + q.table + "` ("
	for key, _ := range data {
		q.sql += "`" + key + "`,"
		keys = append(keys, key)
	}
	q.sql = strings.TrimRight(q.sql, ",")
	q.sql += ") values "
	return keys
}

//兼容批量输入
func (q *query) compactCreateData(datas []map[string]interface{}, keys []string) {
	for _, data := range datas {
		q.sql += "("
		for _, key := range keys {
			q.sql += "?,"
			q.stmtValue = append(q.stmtValue, data[key])
		}
		q.sql = strings.TrimRight(q.sql, ",")
		q.sql += "),"
	}
	q.sql = strings.TrimRight(q.sql, ",")
}

//保存/更新数据
func (q *query) Update(data map[string]interface{}) int {
	q.resetStmt()
	//组装update数据
	q.compactUpdate(data)
	//组装where
	q.compactWhere()

	var stmt *sql.Stmt
	var err error
	stmt, err = q.getStmt()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	defer stmt.Close()
	//执行sql
	res, err := stmt.Exec(q.stmtValue...)
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//清理临时数据
	q.resetAll()
	return int(affectedRows)

}

//组装update数据
func (q *query) compactUpdate(data map[string]interface{}) {
	q.sql = "update `" + q.table + "` set "
	for k, v := range data {
		q.sql += "`" + k + "` = ? ,"
		q.stmtValue = append(q.stmtValue, v)
	}
	q.sql = strings.TrimRight(q.sql, ",")
}

//删除数据
func (q *query) Delete() int {
	q.resetStmt()
	//禁止无条件删除
	if len(q.where) == 0 {
		q.saveError("can not delete without where condition")
		return 0
	}
	//拼装sql
	q.compactDelete()
	//拼装where
	q.compactWhere()
	//进行预处理
	var stmt *sql.Stmt
	var err error
	stmt, err = q.getStmt()
	//预处理检查 记录错误
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	defer stmt.Close()
	//执行语句
	res, err := stmt.Exec(q.stmtValue...)
	//结果分析
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//获取受影响行数
	affectedRows, err := res.RowsAffected()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//清理临时数据
	q.resetAll()
	return int(affectedRows)

}

//预处理delete语句
func (q *query) compactDelete() {
	q.sql = "delete from `" + q.table + "` "
}

//关闭数据库链接
func (q *query) Close() error {
	return q.conn.Close()
}

//ping
func (q *query) Ping() error {
	return q.conn.Ping()
}

//开启事务
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

//事务提交
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

//事务回滚
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

//清空query数据
func (q *query) resetStmt() {
	q.stmtValue = make([]interface{}, 0, 10)
}

//清空执行sql
func (q *query) resetSql() {
	q.sql = ""
}

//清空errors
func (q *query) resetErrors() {
	q.errors = make([]string, 0, 10)
}

//清空sql日志
func (q *query) resetQueryLog() {
	q.queryLog = make([]string, 0, 10)
}

//清空where条件
func (q *query) resetWhere() {
	q.where = make([]where, 0, 10)
}

//清空group条件
func (q *query) resetGroup() {
	q.group = make([]string, 0, 1)
}

//清空orderby条件
func (q *query) resetOrder() {
	q.order = make([]map[string]string, 0, 1)
}

//清空fields条件
func (q *query) resetFields() {
	q.fields = make([]string, 0, 10)
}

//清空limit、offset
func (q *query) resetLimit() {
	q.offset = 0
	q.limit = 0
}

//清空table
func (q *query) resetTable() {
	q.table = ""
}

//清空日志、事务以外的所有数据
func (q *query) resetAll() {
	q.resetStmt()
	q.resetWhere()
	q.resetGroup()
	q.resetOrder()
	q.resetFields()
	q.resetTable()
	q.resetLimit()

}

//清空最后一条执行sql
func (q *query) resetLastQuery() {
	q.lastSql = ""
}

//记录错误日志
func (q *query) saveError(err string) {
	q.errors = append(q.errors, err)
}

//获取最后一个错误信息
func (q *query) GetLastError() string {
	return q.errors[len(q.errors)-1]
}

//获取所有错误信息
func (q *query) GetAllErrors() []string {
	return q.errors
}

//获取最后一条执行的sql
func (q *query) GetLastSQL() string {
	return q.lastSql
}

//获取sql日志
func (q *query) GetSQLs() []string {
	return q.queryLog
}

//启用sql日志 默认关闭
func (q *query) StartLogQuery() {
	q.isLogQuery = true
}

//关闭sql日志
func (q *query) StopLogQuery() {
	q.isLogQuery = false
}

//重置mysql链接以外的所有数据
func (q *query) Reset() {
	if q.txStatus == true {
		q.Rollback()
	}
	q = NewQuery(q.conn)
}
