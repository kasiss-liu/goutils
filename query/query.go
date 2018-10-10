package query

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	//引入mysql层
	_ "github.com/go-sql-driver/mysql"
	//引入sqlite3库
	_ "github.com/mattn/go-sqlite3"
)

//Query query结构体
type Query struct {
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
	lastSQL    string
	joins      []join
}

type where struct {
	field   string
	compare string
	value   interface{}
}

type join struct {
	joinType  string
	joinTable string
	onLeft    string
	onRight   string
	compare   string
}

type queryResult struct {
	Columns []string
	RowsNum int
	Value   []map[string]string
}

//DbConfig mysql 配置结构
type DbConfig struct {
	Host    string
	Port    int
	User    string
	Passwd  string
	Dbname  string
	Charset string
}

//NewDbConfig 获取一个配置结构
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

//NewQuery 获取一个新的query结构
func NewQuery(conn *sql.DB) *Query {
	return &Query{
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

//NewQueryWithConfig 根据配置，连接db并返回一个query结构
func NewQueryWithConfig(conf *DbConfig) (*Query, error) {
	db, err := dbConnection(conf)
	if err != nil {
		return nil, err
	}
	return NewQuery(db), nil

}

//NewQueryWithSqlite 返回一个sqliteDB链接的query结构
func NewQueryWithSqlite(dbPath string) (*Query, error) {
	sqlite, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return NewQuery(sqlite), nil
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

//Table 设置数据表
func (q *Query) Table(t string) *Query {
	q.table = t
	return q
}

//Where 注入sql的where条件
func (q *Query) Where(field string, value interface{}, compare string) *Query {
	var w = where{
		field:   field,
		value:   value,
		compare: compare,
	}
	q.where = append(q.where, w)
	return q
}

//Select 设置要查询的字段
func (q *Query) Select(fields ...string) *Query {
	for _, v := range fields {
		if q.inFields(v) {
			continue
		}
		q.fields = append(q.fields, v)
	}
	return q
}

//OrderBy 设置order
func (q *Query) OrderBy(field string, sort string) *Query {
	order := map[string]string{field: sort}
	q.order = append(q.order, order)
	return q
}

//GroupBy 设置group分组
func (q *Query) GroupBy(g ...string) *Query {
	for _, v := range g {
		q.group = append(q.group, v)
	}
	return q
}

//Limit 设置limit值
func (q *Query) Limit(i ...int) *Query {
	if len(i) == 1 {
		q.limit = i[0]
	}
	if len(i) >= 2 {
		q.offset = i[0]
		q.limit = i[1]
	}
	return q
}

//Offset 设置offset
func (q *Query) Offset(i int) *Query {
	q.offset = i
	return q
}

//判断是否含有元素
func (q *Query) inFields(field string) bool {
	for _, v := range q.fields {
		if v == field {
			return true
		}
	}
	return false
}

//Query 执行查询
func (q *Query) Query() *queryResult {
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

//QueryOne 查询单条记录
func (q *Query) QueryOne() map[string]string {
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
	}
	return nil
}

//QueryRaw 自定义sql查询
func (q *Query) QueryRaw(query string, v ...interface{}) *queryResult {
	q.sql = query
	q.stmtValue = v
	var stmt *sql.Stmt
	var err error
	stmt, err = q.getStmt()
	if err != nil {
		q.saveError(err.Error())
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(q.stmtValue...)
	if err != nil {
		q.saveError(err.Error())
		return nil
	}
	defer rows.Close()
	//清空临时数据
	q.resetAll()
	return q.get(rows)
}

//解析查询结果
func (q *Query) get(rows *sql.Rows) *queryResult {
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
func (q *Query) getStmt() (*sql.Stmt, error) {
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
	q.lastSQL = q.sql

	return stmt, err
}

//拼接sql
func (q *Query) compactQuery() {

	q.resetStmt()
	q.compactSelect()
	q.compactTable()
	q.compactJoin()
	q.compactWhere()
	q.compactGroup()
	q.compactOrder()
	q.compactLimit()
}

//拼接select
func (q *Query) compactSelect() {
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
func (q *Query) compactTable() {

	q.sql += " from " + q.table + " "
}

//拼接where
func (q *Query) compactWhere() {
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
func (q *Query) compactGroup() {
	if len(q.group) > 0 {
		q.sql += " group by "
		var gString string
		for _, v := range q.group {
			gString += v + ","
		}
		gString = strings.Trim(gString, ",")
		q.sql += gString
	}
}

//拼接orderby
func (q *Query) compactOrder() {
	if len(q.order) > 0 {
		q.sql += " order by "
		for _, stringMap := range q.order {
			for k, v := range stringMap {
				q.sql += k + " " + v + ","
			}
		}
		q.sql = strings.Trim(q.sql, ",")
	}
}

//拼接limit
func (q *Query) compactLimit() {
	if q.limit > 0 {
		q.sql += " limit " + strconv.Itoa(q.limit)
	}
	if q.offset > 0 {
		q.sql += " offset " + strconv.Itoa(q.offset)
	}
}

//Create 创建新数据
func (q *Query) Create(data map[string]interface{}) int {
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
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//清理临时数据
	q.resetAll()
	return int(lastInsertID)
}

//CreateBatch 批量创建数据
func (q *Query) CreateBatch(data []map[string]interface{}) int {
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
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		q.saveError(err.Error())
		return 0
	}
	//清理临时数据
	q.resetAll()
	return int(lastInsertID)
}

//ExecRaw 执行自定义sql
func (q *Query) ExecRaw(query string, v ...interface{}) sql.Result {
	q.sql = query
	q.stmtValue = v

	var stmt *sql.Stmt
	var err error
	stmt, err = q.getStmt()
	if err != nil {
		q.saveError(err.Error())
		return nil
	}
	defer stmt.Close()
	res, err := stmt.Exec(q.stmtValue...)
	if err != nil {
		q.saveError(err.Error())
		return nil
	}
	return res
}

//整理新建数据字段
func (q *Query) compactCreateFields(data map[string]interface{}) []string {
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
func (q *Query) compactCreateData(datas []map[string]interface{}, keys []string) {
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

//Update 保存/更新数据
func (q *Query) Update(data map[string]interface{}) int {
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
func (q *Query) compactUpdate(data map[string]interface{}) {
	q.sql = "update `" + q.table + "` set "
	for k, v := range data {
		q.sql += "`" + k + "` = ? ,"
		q.stmtValue = append(q.stmtValue, v)
	}
	q.sql = strings.TrimRight(q.sql, ",")
}

//Delete 删除数据
func (q *Query) Delete() int {
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
func (q *Query) compactDelete() {
	q.sql = "delete from `" + q.table + "` "
}

//Close 关闭数据库链接
func (q *Query) Close() error {
	return q.conn.Close()
}

//Ping ping
func (q *Query) Ping() error {
	return q.conn.Ping()
}

//BeginTransaction 开启事务
func (q *Query) BeginTransaction() error {
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

//Commit 事务提交
func (q *Query) Commit() error {
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

//Rollback 事务回滚
func (q *Query) Rollback() error {
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
func (q *Query) resetStmt() {
	q.stmtValue = make([]interface{}, 0, 10)
}

//清空执行sql
func (q *Query) resetSQL() {
	q.sql = ""
}

//清空errors
func (q *Query) resetErrors() {
	q.errors = make([]string, 0, 10)
}

//清空sql日志
func (q *Query) resetQueryLog() {
	q.queryLog = make([]string, 0, 10)
}

//清空where条件
func (q *Query) resetWhere() {
	q.where = make([]where, 0, 10)
}

//清空group条件
func (q *Query) resetGroup() {
	q.group = make([]string, 0, 1)
}

//清空orderby条件
func (q *Query) resetOrder() {
	q.order = make([]map[string]string, 0, 1)
}

//清空fields条件
func (q *Query) resetFields() {
	q.fields = make([]string, 0, 10)
}

//清空limit、offset
func (q *Query) resetLimit() {
	q.offset = 0
	q.limit = 0
}

//清空table
func (q *Query) resetTable() {
	q.table = ""
}

//清空日志、事务以外的所有数据
func (q *Query) resetAll() {
	q.resetStmt()
	q.resetWhere()
	q.resetGroup()
	q.resetOrder()
	q.resetFields()
	q.resetTable()
	q.resetLimit()
	q.resetSQL()
}

//清空最后一条执行sql
func (q *Query) resetLastQuery() {
	q.lastSQL = ""
}

//记录错误日志
func (q *Query) saveError(err string) {
	q.errors = append(q.errors, err)
}

//GetLastError 获取最后一个错误信息
func (q *Query) GetLastError() string {
	length := len(q.errors)
	if length > 0 {
		return q.errors[length-1]
	}
	return ""
}

//GetAllErrors 获取所有错误信息
func (q *Query) GetAllErrors() []string {
	return q.errors
}

//GetLastSQL 获取最后一条执行的sql
func (q *Query) GetLastSQL() string {
	return q.lastSQL
}

//GetSQLs 获取sql日志
func (q *Query) GetSQLs() []string {
	return q.queryLog
}

//StartLogQuery 启用sql日志 默认关闭
func (q *Query) StartLogQuery() {
	q.isLogQuery = true
}

//StopLogQuery 关闭sql日志
func (q *Query) StopLogQuery() {
	q.isLogQuery = false
}

//ResetErrorLog 清空错误日志
func (q *Query) ResetErrorLog() {
	q.resetErrors()
}

//ResetQueryLog 清空sql日志
func (q *Query) ResetQueryLog() {
	q.resetLastQuery()
	q.resetQueryLog()
}

//Reset 重置mysql链接以外的所有数据
func (q *Query) Reset() {
	if q.txStatus == true {
		q.Rollback()
	}
	q = NewQuery(q.conn)
}

//LeftJoin 增加一个左联查询 并添加连接条件
func (q *Query) LeftJoin(table, onLeft, onRight, compare string) *Query {
	q.joins = append(q.joins, join{
		joinType:  "left join",
		joinTable: table,
		onLeft:    onLeft,
		onRight:   onRight,
		compare:   compare,
	})
	return q
}

//RightJoin 增加一个右联查询 并添加连接条件
func (q *Query) RightJoin(table, onLeft, onRight, compare string) *Query {
	q.joins = append(q.joins, join{
		joinType:  "right join",
		joinTable: table,
		onLeft:    onLeft,
		onRight:   onRight,
		compare:   compare,
	})
	return q
}

//InnerJoin 增加一个内联查询 并添加连接条件
func (q *Query) InnerJoin(table, onLeft, onRight, compare string) *Query {
	q.joins = append(q.joins, join{
		joinType:  "inner join",
		joinTable: table,
		onLeft:    onLeft,
		onRight:   onRight,
		compare:   compare,
	})
	return q
}

//拼接内联语句
func (q *Query) compactJoin() {
	if len(q.joins) < 1 {
		return
	}

	for _, v := range q.joins {
		q.sql += v.joinType + " on "
		q.sql += v.onLeft + " " + v.compare + " " + v.onRight + " "

	}
}

//Tosql 返回一个拼接后的sql
func (q *Query) ToSql() string {
	q.compactQuery()
	return q.sql
}
